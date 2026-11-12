package workbench

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"al.essio.dev/pkg/shellescape"
	"github.com/sylvinus/agent-vm/internal/cage"
	"github.com/sylvinus/agent-vm/internal/containerenv"
	"github.com/sylvinus/agent-vm/internal/diagnostics"
	"github.com/sylvinus/agent-vm/internal/exitpolicy"
	"github.com/sylvinus/agent-vm/internal/paths"
	"github.com/sylvinus/agent-vm/internal/usererror"
	"go.yaml.in/yaml/v3"
)

//go:embed safety_probe.py
var guestSafetyProbe string

type safetyState struct {
	MacOSExitCandidate *macOSExitCandidate `json:"macos_exit_candidate,omitempty"`
	pauseVersion       uint64
	monitorError       string
	Gate               *cage.GateStatus   `json:"gate,omitempty"`
	Report             *cage.SafetyReport `json:"report,omitempty"`
	Error              string             `json:"error,omitempty"`
	ExitIP             string             `json:"exit_ip,omitempty"`
	VerifiedAt         time.Time          `json:"verified_at,omitempty"`
}

type networkObservation struct {
	Status int    `json:"status"`
	IP     string `json:"ip"`
}
type macIPv6Probe struct {
	Name          string `json:"name,omitempty"`
	Result        string `json:"result"`
	GlobalSources int    `json:"globalSources"`
	BoundIndex    int    `json:"boundIndex"`
	RouteVerified bool   `json:"routeVerified"`
}
type guestObservation struct {
	IPv6UnscopedAbsent bool           `json:"ipv6_unscoped_absent,omitempty"`
	IPv6Unbound        macIPv6Probe   `json:"ipv6_unbound,omitempty"`
	IPv6Scoped         []macIPv6Probe `json:"ipv6_scoped,omitempty"`
	IPv6LinkSent       bool           `json:"ipv6_link_sent,omitempty"`
	IPv6LinkInterface  string         `json:"ipv6_link_interface,omitempty"`

	ProcessLocale       string             `json:"process_locale,omitempty"`
	SystemLocale        string             `json:"system_locale,omitempty"`
	SystemLanguage      string             `json:"system_language,omitempty"`
	TZ                  string             `json:"tz"`
	Lang                string             `json:"lang"`
	LCAll               string             `json:"lc_all"`
	SystemTimezone      string             `json:"system_timezone"`
	LocaleValid         bool               `json:"locale_valid"`
	LocaltimeTarget     string             `json:"localtime_target,omitempty"`
	UTCOffsetExpected   *int               `json:"utc_offset_expected,omitempty"`
	UTCOffsetActual     *int               `json:"utc_offset_actual,omitempty"`
	UTCOffsetSystem     *int               `json:"utc_offset_system,omitempty"`
	PublicConnected     bool               `json:"public_connected"`
	ManagementConnected bool               `json:"management_connected"`
	IPv6Connected       bool               `json:"ipv6_connected"`
	IPv6DefaultRoutes   int                `json:"ipv6_default_routes"`
	UDPAttempted        bool               `json:"udp_attempted"`
	PublicDNSAnswered   bool               `json:"public_dns_answered"`
	PrivateFile         string             `json:"private_file"`
	SharedFile          bool               `json:"shared_file"`
	Egress              networkObservation `json:"egress"`
	DNS                 networkObservation `json:"dns"`

	// Reported by the Linux probe only: direct DNS transports, host-side
	// endpoints and resolver configuration, which the host network layer
	// cannot observe from denied connections alone.
	LinkLocalConnected     *bool               `json:"link_local_connected,omitempty"`
	ProxyEndpointConnected *bool               `json:"proxy_endpoint_connected,omitempty"`
	DNSUDPAttempted        bool                `json:"dns_udp_attempted,omitempty"`
	DNSTCPConnected        *bool               `json:"dns_tcp_connected,omitempty"`
	DoTConnected           *bool               `json:"dot_connected,omitempty"`
	Resolvers              resolverObservation `json:"resolvers,omitempty"`
	IPv6GlobalAddresses    *int                `json:"ipv6_global_addresses,omitempty"`
}
type resolverObservation struct {
	Count    *int  `json:"count"`
	Confined *bool `json:"confined"`
}
type browserObservation struct {
	Timezone         string             `json:"timezone"`
	ExpectedTimezone string             `json:"expected_timezone"`
	Language         string             `json:"language"`
	Platform         string             `json:"platform"`
	Network          networkObservation `json:"network"`
}

func (s *Server) gateClient() (*cage.GateClient, error) {
	return cage.OpenGateClient(s.instanceDir())
}

func (s *Server) publishGate(state cage.GateStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.safety.Gate
	if old != nil && old.RunID == state.RunID && old.Epoch > state.Epoch {
		return
	}
	s.safety.Gate = &state
	if state.Mode == cage.GateOpen && state.Report != nil {
		s.safety.Report = state.Report
		// A repeated status poll of the old grant is not a new authorization.
		if old == nil || old.RunID != state.RunID || old.Epoch != state.Epoch {
			s.safety.Error = ""
			s.safety.monitorError = ""
			s.exitProbeError = ""
		}
	}
}

// A gate grant resolves a monitoring failure. Independent diagnostics, such
// as a location lookup failure, are cleared by the operation that owns them.
func (s *Server) currentNetworkErrorLocked() string {
	reason := s.networkErrorLocked()
	if recovery := s.networkRecoveryMessageLocked(time.Now()); recovery != "" {
		if reason == "" {
			return recovery
		}
		return reason + "；" + recovery
	}
	return reason
}

func (s *Server) networkErrorLocked() string {
	if s.networkIntentError != "" {
		return s.networkIntentError
	}
	if s.safety.monitorError != "" {
		return s.safety.monitorError
	}
	if s.exitProbeError != "" {
		return s.exitProbeError
	}
	return s.locationError
}

// A probe can confirm the IP and still fail to resolve its location. Only
// a successful location lookup resolves that second diagnostic.
func (s *Server) recordLocationProbeLocked(loc Location, err error) {
	if err != nil && loc.IP == "" {
		s.exitProbeError = err.Error()
		return
	}
	s.exitProbeError = ""
	if err != nil {
		s.locationError = err.Error()
		return
	}
	s.location = &loc
	if loc.CountryCode != "" {
		s.locationError = ""
	}
}

// Location probes outside the periodic monitor must revoke unsafe observations
// too. Bind the result to the applied proxy generation so a late response from
// a previous node cannot revoke the new node's authorization.
func (s *Server) finishLocationProbe(cfg cage.Config, loc Location, probeErr error) error {
	s.safetyOp.Lock()
	defer s.safetyOp.Unlock()
	current, err := loadConfig(s.configPath)
	if err != nil {
		return err
	}
	if current.ProxyURL != cfg.ProxyURL || current.ProxyRevision != cfg.ProxyRevision || current.DNSURL != cfg.DNSURL {
		return nil
	}
	observed := loc.exitLocation()
	state, reason := exitpolicy.CheckExitCountry(loc.IP, &observed)
	s.mu.Lock()
	// A fresh lookup supersedes any cached result for this address, including
	// an older non-mainland result followed by a rejection or lookup failure.
	delete(s.locations, loc.IP)
	if probeErr != nil && s.location != nil && s.location.IP == loc.IP {
		s.location = nil
	}
	s.recordLocationProbeLocked(loc, probeErr)
	s.safety.ExitIP = loc.IP
	unsafe := probeErr != nil || state != cage.CheckPassed
	if unsafe {
		if s.checkCancel != nil {
			s.checkCancel()
		}
		if s.monitorCancel != nil {
			s.monitorCancel()
		}
		s.safety.Error, s.safety.monitorError = reason, reason
	}
	s.mu.Unlock()
	if !unsafe || !s.vmRunning() {
		return nil
	}
	client, err := s.gateClient()
	if err != nil {
		return err
	}
	closed, err := client.Close(context.Background(), reason)
	if err != nil {
		return err
	}
	s.publishGate(closed)
	if state == cage.CheckUnknown {
		s.scheduleNetworkRecovery(closed)
	}
	return nil
}

func (s *Server) publishReport(report cage.SafetyReport) {
	copy := report
	copy.Checks = append([]cage.SafetyCheck(nil), report.Checks...)
	s.mu.Lock()
	s.safety.Report = &copy
	s.mu.Unlock()
}

func (s *Server) safetyCheck(w http.ResponseWriter, r *http.Request) {
	if s.macosGuest != nil {
		s.launch(w, "safety", s.checkMacOSNetwork)
		return
	}
	s.launch(w, "safety", s.runSafetyCheck)
}

func (s *Server) pauseNetwork(w http.ResponseWriter, r *http.Request) {
	if err := s.pauseNetworkNow("已手动暂停外网；离线终端和文件仍可使用"); err != nil {
		requestError(w, err)
		return
	}
	jsonReply(w, map[string]bool{"paused": true})
}

func (s *Server) pauseNetworkNow(reason string) error {
	return s.revokeNetworkNow(reason, true)
}

func (s *Server) revokeNetworkNow(reason string, pause bool) error {
	var intentErr error
	s.mu.Lock()
	if pause {
		s.recovery = networkRecovery{}
		s.networkPaused = true
		s.safety.pauseVersion++
		intentErr = s.saveNetworkIntentLocked(true)
	}
	if s.checkCancel != nil {
		s.checkCancel()
	}
	if s.monitorCancel != nil {
		s.monitorCancel()
	}
	s.mu.Unlock()
	if pause {
		s.hostTailnetForwards.closeAll()
	}
	s.safetyOp.Lock()
	defer s.safetyOp.Unlock()
	client, err := s.gateClient()
	if err != nil {
		return errors.Join(intentErr, err)
	}
	state, err := client.Close(context.Background(), reason)
	if err == nil {
		s.publishGate(state)
	}
	return errors.Join(intentErr, err)
}

func (s *Server) safetyExport(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	report := s.safety.Report
	s.mu.Unlock()
	if report == nil {
		http.Error(w, "尚无检测报告", 404)
		return
	}
	// Exports cannot authorize networking and contain no proxy hash, address,
	// credentials, launch ticket or host filesystem paths.
	jsonReply(w, map[string]any{"policy": report.Policy, "created_at": report.CreatedAt, "checks": report.Checks, "scope": "采样与运行时策略检查；不代表平台账号资格或完整隔离审计"})
}

// The probe runs in Linux, while its files are created on the host. Use the
// same guest path mapping as the VM's mounts for both sides of the boundary.
func boundaryProbeInput(port int, nonce, privateFile, sharedFile, challenge string) map[string]any {
	return map[string]any{"stage": "boundary", "control_port": port, "nonce": nonce,
		"management_host": "192.168.5.2",
		"private_file":    guestPath(privateFile), "shared_file": guestPath(sharedFile), "challenge": challenge}
}

func (s *Server) guestObservation(ctx context.Context, input map[string]any) (guestObservation, error) {
	if s.windowsGuest != nil {
		return s.windowsGuestObservation(ctx, input)
	}
	if s.macosGuest != nil {
		return s.macOSGuestObservation(ctx, input)
	}
	var observation guestObservation
	data, err := json.Marshal(input)
	if err != nil {
		return observation, err
	}
	command := "python3 -c " + shellescape.Quote(guestSafetyProbe) + " " + shellescape.Quote(base64.StdEncoding.EncodeToString(data))
	output, err := s.guestControl(ctx, "exec bash -lc "+shellescape.Quote(command), "safety-probe")
	if err != nil {
		return observation, fmt.Errorf("客机检测未返回完整结果: %w", err)
	}
	if err = json.Unmarshal(output, &observation); err != nil {
		return observation, errors.New("客机检测返回格式无效")
	}
	return observation, nil
}

func (s *Server) browserObservation(ctx context.Context, cfg cage.Config, network bool) (browserObservation, error) {
	if s.macosGuest != nil {
		return s.macOSBrowserObservation(ctx, cfg, network)
	}
	var observation browserObservation
	s.mu.Lock()
	token := s.credentials.Token
	s.mu.Unlock()
	body, _ := json.Marshal(map[string]any{"network": network, "timezone": cfg.Timezone})
	data, err := s.guestAPI(ctx, token, "POST", "/audit", body)
	if err != nil || json.Unmarshal(data, &observation) != nil {
		return observation, errors.New("浏览器检测未返回完整结果，请更新或重启环境")
	}
	return observation, nil
}

// Retry only a transient network sample, before the failure drill. Settings and
// policy failures retain their evidence; no request grants external access.
func (s *Server) browserExitObservation(ctx context.Context, cfg cage.Config) (browserObservation, error) {
	observation, err := s.browserObservation(ctx, cfg, true)
	if err != nil || s.browserSettingsIssue(observation, cfg) != "" ||
		(observation.Network.Status != 0 && observation.Network.Status != http.StatusServiceUnavailable) {
		return observation, err
	}
	_ = diagnostics.New(s.profile).Record("warn", "safety-probe", fmt.Sprintf("stage=browser retry=1 status=%d", observation.Network.Status))
	s.setMessage("浏览器出口检测暂时未成功，正在重试一次；外网保持关闭…")
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return observation, ctx.Err()
	case <-timer.C:
		return s.browserObservation(ctx, cfg, true)
	}
}

func guestSettingsMatch(v guestObservation, cfg cage.Config) bool {
	return linuxGuestSettingsMismatch(v, cfg) == ""
}

// linuxGuestSettingsMismatch names every configured value the Linux guest did
// not read back. Besides the variables and /etc/timezone it requires the UTC
// offset programs actually observe, so a stale /etc/localtime or missing
// tzdata cannot hide behind a correct TZ variable.
func linuxGuestSettingsMismatch(v guestObservation, cfg cage.Config) string {
	var mismatches []string
	if !v.LocaleValid {
		mismatches = append(mismatches, "系统语言环境无法加载")
	}
	if v.TZ != cfg.Timezone {
		mismatches = append(mismatches, fmt.Sprintf("TZ 变量不一致（期望 %q，实际 %q）", cfg.Timezone, v.TZ))
	}
	if v.SystemTimezone != cfg.Timezone {
		mismatches = append(mismatches, fmt.Sprintf("系统时区文件不一致（期望 %q，实际 %q）", cfg.Timezone, v.SystemTimezone))
	}
	if v.Lang != cfg.Locale {
		mismatches = append(mismatches, fmt.Sprintf("LANG 不一致（期望 %q，实际 %q）", cfg.Locale, v.Lang))
	}
	if v.LCAll != cfg.Locale {
		mismatches = append(mismatches, fmt.Sprintf("LC_ALL 不一致（期望 %q，实际 %q）", cfg.Locale, v.LCAll))
	}
	if issue := effectiveTimezoneIssue(v, cfg); issue != "" {
		mismatches = append(mismatches, issue)
	}
	return strings.Join(mismatches, "；")
}

func effectiveTimezoneIssue(v guestObservation, cfg cage.Config) string {
	if v.UTCOffsetExpected == nil || v.UTCOffsetActual == nil {
		return "未读到客机的实际 UTC 偏移"
	}
	if *v.UTCOffsetExpected != *v.UTCOffsetActual {
		return fmt.Sprintf("实际 UTC 偏移与时区不一致（期望 %+d 秒，实际 %+d 秒）", *v.UTCOffsetExpected, *v.UTCOffsetActual)
	}
	// Read the system file separately: time.localtime observes TZ, which can
	// mask a stale or missing /etc/localtime even when it is a regular file.
	if v.UTCOffsetSystem == nil {
		return "未读到系统时区文件的实际 UTC 偏移"
	}
	if *v.UTCOffsetSystem != *v.UTCOffsetExpected {
		return fmt.Sprintf("系统时区文件的 UTC 偏移不一致（期望 %+d 秒，实际 %+d 秒）", *v.UTCOffsetExpected, *v.UTCOffsetSystem)
	}
	if v.LocaltimeTarget != "" && v.LocaltimeTarget != "/usr/share/zoneinfo/"+cfg.Timezone {
		return fmt.Sprintf("系统本地时间链接不一致（期望 %q，实际 %q）", "/usr/share/zoneinfo/"+cfg.Timezone, v.LocaltimeTarget)
	}
	return ""
}
func browserSettingsMatch(v browserObservation, cfg cage.Config) bool {
	return browserSettingsMismatch(v, cfg) == ""
}

// guestProxyEndpoint is the user's proxy as the guest would have to address
// it: a proxy bound on the host is reached through the VM gateway address.
// The guest must never reach it directly, only through the host network layer.
func guestProxyEndpoint(cfg cage.Config, managementHost string) (string, int) {
	u, err := url.Parse(cfg.ProxyURL)
	if err != nil {
		return "", 0
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil || port <= 0 {
		return "", 0
	}
	host := u.Hostname()
	if ip := net.ParseIP(host); strings.EqualFold(host, "localhost") || (ip != nil && (ip.IsLoopback() || ip.IsUnspecified())) {
		host = managementHost
	}
	return host, port
}

// Missing Linux observations stay unknown; confirmed means a connection succeeded.
func hostReachIssue(v guestObservation, linux, proxyExpected bool) (issue string, confirmed bool) {
	switch {
	case v.LinkLocalConnected != nil && *v.LinkLocalConnected:
		return "客机直接连通了链路本地元数据地址（169.254.169.254）", true
	case v.ProxyEndpointConnected != nil && *v.ProxyEndpointConnected:
		return "客机绕过环境外网络层直接连通了用户代理端口", true
	case linux && v.LinkLocalConnected == nil:
		return "未读到链路本地元数据地址的直连结果", false
	case linux && proxyExpected && v.ProxyEndpointConnected == nil:
		return "未读到用户代理端口的直连结果", false
	}
	return "", false
}

// dnsBypassIssue reports a way around the encrypted DNS path. confirmed is
// true when the guest actually reached a resolver outside the environment.
func dnsBypassIssue(v guestObservation, linux bool) (issue string, confirmed bool) {
	switch {
	case v.DNSTCPConnected != nil && *v.DNSTCPConnected:
		return "客机直接连通了公共 DNS 的 TCP/53", true
	case v.DoTConnected != nil && *v.DoTConnected:
		return "客机直接连通了公共 DNS 的 DoT/853", true
	case !linux:
		return "", false
	case !v.DNSUDPAttempted:
		return "客机的 UDP/53 直连探针没有发出", false
	case v.DNSTCPConnected == nil || v.DoTConnected == nil:
		return "未读到公共 DNS 的 TCP/53 或 DoT/853 直连结果", false
	case v.Resolvers.Count == nil || *v.Resolvers.Count < 0 || v.Resolvers.Confined == nil:
		return "未能读取客机的解析器配置", false
	case *v.Resolvers.Count == 0:
		return "客机没有配置任何解析器", false
	case !*v.Resolvers.Confined:
		return "客机解析器指向环境外的地址，查询可能绕过加密 DNS", true
	}
	return "", false
}

func (s *Server) ipv6PolicyVerified(v guestObservation, before, after cage.GateStatus, networkLayout bool) bool {
	if v.IPv6Connected || !after.PolicyVerified {
		return false
	}
	if s.macosGuest != nil {
		return networkLayout && (macOSIPv6LinkVerified(v, before, after) ||
			(v.IPv6DefaultRoutes == 0 && v.IPv6UnscopedAbsent && v.IPv6Unbound.Result == "no-route"))
	}
	if s.windowsGuest == nil && (v.IPv6GlobalAddresses == nil || *v.IPv6GlobalAddresses != 0) {
		return false
	}
	return v.IPv6DefaultRoutes == 0 || after.Counters.IPv6Denied > before.Counters.IPv6Denied
}

func browserSettingsMismatch(v browserObservation, cfg cage.Config) string {
	var mismatches []string
	if v.Timezone == "" || v.Timezone != v.ExpectedTimezone {
		mismatches = append(mismatches, fmt.Sprintf("时区不一致（期望 %q，实际 %q）", cfg.Timezone, v.Timezone))
	}
	// Chromium's Japanese UI locale may report the language-only tag even
	// with ja-JP preferences. Keep other regional language distinctions.
	japanese := strings.EqualFold(cfg.BrowserLanguage, "ja-JP") && strings.EqualFold(v.Language, "ja")
	if !strings.EqualFold(v.Language, cfg.BrowserLanguage) && !japanese {
		mismatches = append(mismatches, fmt.Sprintf("语言不一致（期望 %q，实际 %q）", cfg.BrowserLanguage, v.Language))
	}
	if !strings.HasPrefix(v.Platform, "Linux") {
		mismatches = append(mismatches, fmt.Sprintf("系统标识不一致（期望 Linux，实际 %q）", v.Platform))
	}
	return strings.Join(mismatches, "；")
}

func (s *Server) inspectVMLayout() (bool, bool) {
	if s.windowsGuest != nil {
		return s.windowsVMLayout()
	}
	if s.macosGuest != nil {
		return s.macOSVMLayout()
	}
	if containerenv.Configured(s.profile) {
		environment, err := containerenv.Load(s.profile)
		if err != nil {
			return false, false
		}
		ctx, cancel := s.controlContext()
		defer cancel()
		inspection, err := environment.Inspect(ctx)
		if err != nil {
			return false, false
		}
		return inspection.FilesystemConfined, inspection.NetworkConfined
	}
	workspace, err := validateWorkspace(s.environmentDir(), s.profile)
	if err != nil {
		return false, false
	}
	data, err := os.ReadFile(filepath.Join(s.profile, "state", "lima", cage.InstanceName(s.profile, workspace), "lima.yaml"))
	if err != nil {
		return false, false
	}
	var spec struct {
		Mounts []struct {
			Location   string `yaml:"location"`
			MountPoint string `yaml:"mountPoint"`
			Writable   bool   `yaml:"writable"`
			SSHFS      struct {
				Driver string `yaml:"sftpDriver"`
			} `yaml:"sshfs"`
		} `yaml:"mounts"`
		Networks  []any  `yaml:"networks"`
		MountType string `yaml:"mountType"`
	}
	if yaml.Unmarshal(data, &spec) != nil {
		return false, false
	}
	shares := len(spec.Mounts) == 1 && spec.MountType == "reverse-sshfs"
	if shares {
		m := spec.Mounts[0]
		shares = paths.Equal(paths.Host(m.Location), paths.Host(workspace)) && m.MountPoint == guestPath(workspace) && m.Writable && m.SSHFS.Driver == "builtin"
	}
	return shares, len(spec.Networks) == 0
}

func (s *Server) managementRequiresSession(ctx context.Context) bool {
	client := &http.Client{Transport: &http.Transport{Proxy: nil}, Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("redirect refused") }}
	defer client.CloseIdleConnections()
	req, _ := http.NewRequestWithContext(ctx, "GET", s.origin+"/api/settings", nil)
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	response.Body.Close()
	return response.StatusCode == http.StatusUnauthorized
}

// A location is cached by its measured address, never borrowed from the last
// proxy request. The provider's free quota makes per-heartbeat lookups wasteful.
func (s *Server) locationForIP(ctx context.Context, cfg cage.Config, ip string) (Location, error) {
	s.mu.Lock()
	loc, ok := s.locations[ip]
	if !ok && s.location != nil && s.location.IP == ip {
		loc, ok = *s.location, true
	}
	s.mu.Unlock()
	if ok && time.Since(loc.CheckedAt) < time.Hour && loc.CountryCode != "" {
		return loc, nil
	}
	loc, err := probeLocationForIP(ctx, cfg, ip)
	if err != nil {
		return Location{}, err
	}
	s.mu.Lock()
	if s.locations == nil {
		s.locations = map[string]Location{}
	}
	if len(s.locations) >= 64 {
		for key, value := range s.locations {
			if time.Since(value.CheckedAt) >= time.Hour {
				delete(s.locations, key)
			}
		}
		if len(s.locations) >= 64 {
			for key := range s.locations {
				delete(s.locations, key)
				break
			}
		}
	}
	s.locations[ip] = loc
	s.mu.Unlock()
	return loc, nil
}

func (s *Server) checkExitPolicy(ctx context.Context, cfg cage.Config, ip string) (cage.CheckState, string) {
	policy := cfg.ExitPolicy.Normalized()
	if state, message := exitpolicy.CheckExitCountry(ip, nil); state == cage.CheckFailed {
		return state, message
	}
	loc, err := s.locationForIP(ctx, cfg, ip)
	if err != nil {
		return cage.CheckUnknown, "出口定位暂时不可用，无法确认非中国大陆出口，外网保持关闭"
	}
	observed := loc.exitLocation()
	state, message := policy.Evaluate(ip, &observed)
	s.mu.Lock()
	s.recordLocationProbeLocked(loc, nil)
	s.mu.Unlock()
	return state, message
}

func (s *Server) runSafetyCheck() (returned error) {
	defer func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.networkPaused {
			returned = context.Canceled
			s.safety.Error = ""
		}
	}()
	if err := s.attach(); err != nil {
		return err
	}
	cfg, err := loadConfig(s.configPath)
	if err != nil {
		return err
	}
	if err = cfg.RequireManagedProxy(s.profile); err != nil {
		_ = s.closeProxyGate(err.Error())
		return err
	}
	client, err := s.gateClient()
	if err != nil {
		s.mu.Lock()
		s.safety.Error = err.Error()
		s.mu.Unlock()
		return err
	}
	ctx, cancel := context.WithTimeout(s.operationContext(), 170*time.Second)
	defer cancel()
	s.safetyOp.Lock()
	s.mu.Lock()
	paused := s.networkPaused
	if !paused {
		s.checkCancel = cancel
		s.safety.MacOSExitCandidate = nil
	}
	s.mu.Unlock()
	if paused {
		s.safetyOp.Unlock()
		return errors.New("外网已手动暂停；可重新检测以恢复联网")
	}
	defer func() { s.mu.Lock(); s.checkCancel = nil; s.mu.Unlock() }()
	state, err := client.BeginCheck(ctx)
	if err == nil {
		s.publishGate(state)
	}
	s.safetyOp.Unlock()
	if err != nil {
		_, _ = client.Close(context.Background(), "检测未启动，外网保持关闭")
		return err
	}
	granted := false
	defer func() {
		s.safetyOp.Lock()
		defer s.safetyOp.Unlock()
		s.mu.Lock()
		paused := s.networkPaused
		s.mu.Unlock()
		if !granted && !paused {
			reason := "联网检查尚未完成，还没有取得全部检测结果。外网保持关闭。"
			if returned != nil {
				reason = returned.Error()
			}
			closed, e := client.Close(context.Background(), reason)
			if e == nil {
				s.publishGate(closed)
			}
		}
		s.mu.Lock()
		if returned != nil {
			s.safety.Error = returned.Error()
		} else {
			s.safety.Error = ""
		}
		s.mu.Unlock()
	}()
	if cfg.ExitPolicy.NeedsBinding() {
		s.setMessage("正在首次绑定此出口节点的 IP 或地区…")
		loc, err := probeLocation(ctx, cfg)
		if err != nil {
			return err
		}
		cfg.ExitPolicy = cfg.ExitPolicy.BindMissing(loc.exitLocation())
		observed := loc.exitLocation()
		if result, message := cfg.ExitPolicy.Evaluate(loc.IP, &observed); result != cage.CheckPassed {
			return errors.New(message)
		}
		if err := s.saveExitPolicyConfig(cfg); err != nil {
			return err
		}
		// Completing an empty anchor changes the fingerprint. Revoke the first
		// check and collect all evidence again under the bound configuration.
		s.safetyOp.Lock()
		_, err = client.SyncProxyRevision(ctx)
		if err == nil {
			state, err = client.BeginCheck(ctx)
		}
		s.safetyOp.Unlock()
		if err != nil {
			return err
		}
		s.publishGate(state)
	}
	report := cage.SafetyReport{Policy: cage.SafetyPolicyVersion, ID: randomToken()[:32], RunID: state.RunID, CheckID: state.CheckID, Configuration: cfg.Fingerprint(), CreatedAt: time.Now()}
	defer func() {
		// Preserve partial reports too: cancellation or a transport failure must
		// not leave the last successful report as the only diagnostic evidence.
		s.publishReport(report)
		if err := writePrivate(filepath.Join(s.profile, "last-safety-report.json"), report); err != nil {
			_ = diagnostics.New(s.profile).Record("error", "safety", "保存检测报告失败: "+err.Error())
		}
		for _, check := range report.Checks {
			level := "info"
			if check.Critical && check.State != cage.CheckPassed {
				level = "warn"
			}
			_ = diagnostics.New(s.profile).Record(level, "safety", fmt.Sprintf("check_started=%s check=%s state=%s critical=%t evidence=%s remedy=%s operation_error=%v", report.CreatedAt.UTC().Format(time.RFC3339Nano), check.ID, check.State, check.Critical, check.Evidence, check.Remedy, returned))
		}
	}()
	for _, item := range [][2]string{{"egress_policy", "环境外出站约束"}, {"proxy_transport", "终端检测通道"}, {"dns_transport", "DNS 与解析约束"}, {"guest_settings", "系统时区与语言"}, {"browser_settings", "浏览器时区与语言"}, {"filesystem_boundary", "共享文件与宿主私有目录"}, {"management_boundary", "宿主管理接口"}, {"fail_closed", "代理故障关闭"}, {"udp_policy", "UDP 出站约束"}, {"ipv6_policy", "IPv6 出站约束"}, {"exit_consistency", "出口符合已保存规则"}} {
		report.Checks = append(report.Checks, cage.SafetyCheck{ID: item[0], Title: item[1], Critical: true, State: cage.CheckUnknown, Evidence: "尚未取得证据", Remedy: "完成检测后查看结果；未知项目不会自动通过", CheckedAt: report.CreatedAt})
	}
	set := func(id string, result cage.CheckState, evidence, remedy string) {
		for i := range report.Checks {
			if report.Checks[i].ID == id {
				report.Checks[i].State = result
				report.Checks[i].Evidence = evidence
				report.Checks[i].Remedy = remedy
				report.Checks[i].CheckedAt = time.Now()
			}
		}
		s.publishReport(report)
	}
	s.publishReport(report)
	if s.macosGuest != nil {
		// BeginCheck has already revoked external access. This path is shared
		// by initial setup, startup, controller reconnects and later rechecks.
		if err := s.syncMacOSRegion(ctx, cfg); err != nil {
			set("guest_settings", cage.CheckUnknown, err.Error(), "核对已保存的时区和语言后重新检测")
			return err
		}
	}
	shares, networkLayout := s.inspectVMLayout()
	private, err := os.CreateTemp(s.profile, ".cage-private-probe-*")
	if err != nil {
		return usererror.Wrap("创建私有目录检查文件", err, "请检查 Cage 配置目录的写入权限和磁盘剩余空间，然后重新检测。外网保持关闭。")
	}
	defer os.Remove(private.Name())
	defer private.Close()
	shareRoot := s.environmentDir()
	if s.macosGuest != nil {
		shareRoot = s.workspace()
		approved, e := s.projectWorkspace(shareRoot)
		if e != nil || isClientMountPath(approved) {
			return errors.New("请先在 macOS 工作台连接默认工作区或已授权的主机项目，再验证文件边界")
		}
		if e = s.ensureProjectShare(approved); e != nil {
			set("filesystem_boundary", cage.CheckUnknown, e.Error(), "先完成客机登录、文件系统扩展或工作区访问授权，再重新运行检查")
			var failure *projectShareError
			if errors.As(e, &failure) {
				return &projectShareError{Path: approved, Recovery: failure.Recovery, Message: failure.Message}
			}
			return e
		}
		shareRoot = approved
	}
	shared, err := os.CreateTemp(shareRoot, ".cage-share-probe-*")
	if err != nil {
		return usererror.Wrap("创建项目共享检查文件", err, "请确认环境目录仍存在且可写，再重新检测；不要删除个人配置目录。外网保持关闭。")
	}
	defer os.Remove(shared.Name())
	defer shared.Close()
	challenge := randomToken()
	if _, err = shared.WriteString(challenge); err != nil {
		return err
	}
	if err = shared.Sync(); err != nil {
		return err
	}
	control, _ := url.Parse(s.origin)
	port, _ := strconv.Atoi(control.Port())
	s.setMessage("正在验证离线约束、共享文件与系统配置…")
	boundaryInput := boundaryProbeInput(port, report.ID, private.Name(), shared.Name(), challenge)
	boundaryInput["timezone"] = cfg.Timezone
	if s.windowsGuest != nil {
		if err := s.windowsBoundaryInput(shareRoot, private.Name(), shared.Name(), challenge, boundaryInput); err != nil {
			return err
		}
	}
	if s.macosGuest != nil {
		boundaryInput["run_id"] = state.RunID
		boundaryInput["check_id"] = state.CheckID
	}
	if containerenv.Configured(s.profile) {
		boundaryInput["management_host"] = containerenv.GatewayIP
	}
	linuxProbe := s.windowsGuest == nil && s.macosGuest == nil
	if linuxProbe {
		if host, port := guestProxyEndpoint(cfg, boundaryInput["management_host"].(string)); port != 0 {
			boundaryInput["proxy_host"], boundaryInput["proxy_port"] = host, port
		}
	}
	boundary, boundaryErr := s.guestObservation(ctx, boundaryInput)
	after, err := client.Status(ctx)
	_ = diagnostics.New(s.profile).Record("info", "safety-probe", fmt.Sprintf("stage=boundary error=%v status_error=%v observation=%+v counters_before=%+v counters_after=%+v shares=%t network_layout=%t", boundaryErr, err, boundary, state.Counters, after.Counters, shares, networkLayout))
	if err != nil {
		return err
	}
	if boundaryErr == nil {
		if s.guestSettingsOK(boundary, cfg) {
			evidence := "客机进程与系统的时区和语言读回结果符合已应用配置"
			if linuxProbe {
				evidence += "；进程与系统时区文件的实际 UTC 偏移均符合已应用时区"
			}
			set("guest_settings", cage.CheckPassed, evidence, "")
		} else {
			evidence := "系统读回结果与已应用配置不一致"
			if linuxProbe {
				if detail := linuxGuestSettingsMismatch(boundary, cfg); detail != "" {
					evidence += "：" + detail
				}
			}
			set("guest_settings", cage.CheckFailed, evidence, "重新应用时区与语言后重启环境")
		}
		_, _ = shared.Seek(0, io.SeekStart)
		sharedData, _ := io.ReadAll(io.LimitReader(shared, 256))
		if shares && boundary.PrivateFile == "unreachable" && boundary.SharedFile && string(sharedData) == challenge+":guest" {
			set("filesystem_boundary", cage.CheckPassed, "共享读写探针往返成功；宿主私有目录探针不可读；环境基础共享目录与明确打开的项目使用受目录约束的传输，项目内允许 Git 读写", "")
		} else {
			set("filesystem_boundary", cage.CheckFailed, "共享读写、私有目录探针或挂载范围未符合预期", "检查共享目录；不要共享主目录或 Cage 配置目录")
		}
		if after.PolicyVerified && networkLayout && after.Configuration == cfg.Fingerprint() && !boundary.PublicConnected && after.Counters.TCPDenied > state.Counters.TCPDenied {
			set("egress_policy", cage.CheckPassed, "运行中策略检查符合预期；未配置额外网卡；检查期间的公网 TCP 请求被 环境外网络层拒绝", "")
		}
		hostReach, hostReachConfirmed := hostReachIssue(boundary, linuxProbe, boundaryInput["proxy_port"] != nil)
		if hostReach == "" && !boundary.ManagementConnected && after.Counters.PrivateDenied > state.Counters.PrivateDenied && s.managementRequiresSession(ctx) {
			evidence := "guest 到宿主管理端口的连接被拒绝；本机接口无会话访问被拒绝"
			if linuxProbe {
				evidence += "；链路本地元数据地址的直连探针未连通"
				if boundaryInput["proxy_port"] != nil {
					evidence += "；用户代理端口的直连探针未连通"
				}
			}
			set("management_boundary", cage.CheckPassed, evidence, "")
		} else {
			result, evidence := cage.CheckUnknown, "管理接口拒绝访问的证据不完整"
			if boundary.ManagementConnected {
				result = cage.CheckFailed
			}
			if hostReach != "" {
				evidence = hostReach
				if hostReachConfirmed {
					result = cage.CheckFailed
				}
			}
			set("management_boundary", result, evidence, "停止环境并重新启动；不要开放宿主控制端口")
		}
		if boundary.UDPAttempted && after.Counters.UDPDenied > state.Counters.UDPDenied && after.PolicyVerified {
			set("udp_policy", cage.CheckPassed, "客机 UDP 探针已发起，环境外网络层记录了拒绝；无响应本身不作为通过依据", "")
		}
		ipv6Verified := s.ipv6PolicyVerified(boundary, state, after, networkLayout)
		if ipv6Verified && s.macosGuest != nil && macOSIPv6LinkVerified(boundary, state, after) {
			set("ipv6_policy", cage.CheckPassed, fmt.Sprintf("常规 IPv6 查路由明确不可达；%d 条接口限定的默认路由逐个绑定检查，均未连通或发现外部源地址。唯一客机网卡 %s 对本轮随机 IPv6 UDP 探针在以太帧入口实施丢弃。", boundary.IPv6DefaultRoutes, boundary.IPv6LinkInterface), "")
		} else if ipv6Verified {
			evidence := "客机公网 IPv6 策略为拒绝；连接采样未建立，路由或拒绝记录与策略一致。当前版本验证 IPv4 出口"
			if linuxProbe {
				evidence += fmt.Sprintf("；客机全局 IPv6 地址 %d 个、默认路由 %d 条", *boundary.IPv6GlobalAddresses, boundary.IPv6DefaultRoutes)
			}
			set("ipv6_policy", cage.CheckPassed, evidence, "")
		} else if s.macosGuest != nil {
			set("ipv6_policy", cage.CheckUnknown, fmt.Sprintf("IPv6 默认路由记录为 %d。普通连接、逐接口观测与当前 VZ 网卡丢弃证据尚未全部匹配。", boundary.IPv6DefaultRoutes), "保持外网关闭。确认使用配套的客机诊断工具和虚拟机运行器，重启后重新检测")
		} else if linuxProbe {
			set("ipv6_policy", cage.CheckUnknown, "未能确认客机全局 IPv6 地址为零，或 IPv6 路由与拒绝证据不完整", "保持外网关闭。确认客机 IPv6 配置后重新检测")
		}
	} else {
		message := usererror.Display("读取环境内的隔离检查结果", boundaryErr, "请确认环境仍在运行，再重新检测；持续失败时导出诊断包。").Error()
		for _, id := range []string{"guest_settings", "filesystem_boundary", "egress_policy", "management_boundary", "udp_policy", "ipv6_policy"} {
			set(id, cage.CheckUnknown, message, "重新检测；若同一检查持续无结果，请在设置中导出诊断包")
		}
	}
	s.setMessage("正在通过固定检测通道检查代理与 DNS…")
	network, networkErr := s.guestObservation(ctx, map[string]any{"stage": "network"})
	normalState, err := client.Status(ctx)
	_ = diagnostics.New(s.profile).Record("info", "safety-probe", fmt.Sprintf("stage=network error=%v status_error=%v observation=%+v counters=%+v", networkErr, err, network, normalState.Counters))
	if err != nil {
		return err
	}
	if networkErr == nil && network.Egress.Status == 200 && network.Egress.IP != "" && normalState.Counters.ProbeTCP > after.Counters.ProbeTCP {
		set("proxy_transport", cage.CheckPassed, "终端请求经过 环境网络层的固定检测通道，并由指定代理完成 HTTPS 出口查询；未放开任意网站", "")
	} else {
		evidence := fmt.Sprintf("环境未能完成代理出口查询（检测通道 HTTP %d），尚不能确认代理连接可用", network.Egress.Status)
		if network.Egress.Status == 0 {
			evidence = "环境没有收到代理出口查询的结果，尚不能确认代理连接可用"
		}
		if networkErr != nil {
			evidence = usererror.Display("读取代理检测结果", networkErr, "请确认环境仍在运行。").Error()
		}
		set("proxy_transport", cage.CheckUnknown, evidence, "在网络与代理设置中确认节点已应用并能连接；再重新检测")
	}
	dnsBypass, dnsBypassConfirmed := dnsBypassIssue(boundary, linuxProbe)
	if networkErr == nil && boundaryErr == nil && !boundary.PublicDNSAnswered && after.Counters.DNSDenied > state.Counters.DNSDenied && network.DNS.Status == 200 && normalState.Counters.ProbeDNS > after.Counters.ProbeDNS && dnsBypass == "" {
		evidence := "检查期间公共解析被拒绝；固定 DNS 查询经指定 DoH 与代理完成"
		if linuxProbe {
			evidence += "；客机解析器仅指向环境内地址，到公共 DNS 的 UDP/53、TCP/53 与 DoT/853 直连均未绕过网关"
		}
		set("dns_transport", cage.CheckPassed, evidence, "")
	} else {
		result, evidence := cage.CheckUnknown, "公共解析拒绝记录或加密 DNS 查询证据不完整"
		if boundaryErr == nil && dnsBypass != "" {
			evidence = dnsBypass
			if dnsBypassConfirmed {
				result = cage.CheckFailed
			}
		}
		set("dns_transport", result, evidence, "检查 DNS 配置、代理连通性与环境组件")
	}
	s.setMessage("正在读取隔离浏览器的实际配置…")
	browser, browserErr := s.browserExitObservation(ctx, cfg)
	_ = diagnostics.New(s.profile).Record("info", "safety-probe", fmt.Sprintf("stage=browser error=%v observation=%+v", browserErr, browser))
	if browserErr == nil {
		if s.browserSettingsIssue(browser, cfg) == "" {
			set("browser_settings", cage.CheckPassed, "浏览器诊断页面读到客机平台；实际时区与首选语言符合已应用配置", "")
		} else {
			set("browser_settings", cage.CheckFailed, s.browserSettingsIssue(browser, cfg), "更新环境组件以补齐浏览器语言包并应用已保存设置；仍失败时查看各字段的期望值与实际值")
		}
	} else {
		set("browser_settings", cage.CheckUnknown, browserErr.Error(), "停止并重启环境以加载最新浏览器组件，再重新检测")
	}
	s.setMessage("正在演练代理失败；现有代理进程保持运行…")
	faultStart, err := client.SetFault(ctx, true)
	if err != nil {
		return err
	}
	fault, faultErr := s.guestObservation(ctx, map[string]any{"stage": "network"})
	faultBrowser, faultBrowserErr := s.browserObservation(ctx, cfg, true)
	faultEnd, err := client.Status(ctx)
	_ = diagnostics.New(s.profile).Record("info", "safety-probe", fmt.Sprintf("stage=fault terminal_error=%v browser_error=%v status_error=%v terminal=%+v browser=%+v counters_before=%+v counters_after=%+v", faultErr, faultBrowserErr, err, fault, faultBrowser, faultStart.Counters, faultEnd.Counters))
	if err != nil {
		return err
	}
	if faultErr == nil && faultBrowserErr == nil && fault.Egress.Status == 503 && fault.DNS.Status == 503 && faultBrowser.Network.Status == 503 && faultEnd.Counters.FaultTCP >= faultStart.Counters.FaultTCP+2 && faultEnd.Counters.FaultDNS > faultStart.Counters.FaultDNS {
		set("fail_closed", cage.CheckPassed, "在检测专用连接中注入代理不可用：终端、浏览器和 DNS 均失败，未改为直连；未停止用户的代理进程", "")
	} else {
		set("fail_closed", cage.CheckUnknown, "故障演练缺少终端、浏览器或 DNS 的失败记录", "修复检测通道或浏览器组件后重试；未知不等于已验证安全")
	}
	if _, err = client.SetFault(ctx, false); err != nil {
		return err
	}
	s.setMessage("正在核对静态 IP 或动态地区范围…")
	hostIP, hostErr := cage.ProbeExitIP(ctx, cfg)
	_ = diagnostics.New(s.profile).Record("info", "safety-probe", fmt.Sprintf("stage=host-exit error=%v ip=%s", hostErr, hostIP))
	result, message := cage.CheckPassed, "终端、浏览器与宿主检测出口均未命中大陆网段，定位为非中国大陆且符合已保存的 IP 变化规则；此结论限于检测采样"
	if missing := incompleteExitSamples(hostIP, hostErr, network.Egress, networkErr, browser.Network, browserErr); missing != "" {
		result, message = cage.CheckUnknown, missing
	} else {
		for _, ip := range []string{hostIP, network.Egress.IP, browser.Network.IP} {
			observed, why := s.checkExitPolicy(ctx, cfg, ip)
			if observed == cage.CheckFailed {
				result, message = observed, why
				break
			}
			if observed == cage.CheckUnknown {
				result, message = observed, why
			}
		}
	}
	set("exit_consistency", result, message, "选择非中国大陆出口，并核对静态 IP 或允许地区；静态 IP 和不限地区均须确认出口国家，已有基准不会随重检改变")
	report.Checks = append(report.Checks, cage.SafetyCheck{ID: "third_party_reputation", Title: "第三方 IP 信誉", State: cage.CheckUnknown, CheckedAt: time.Now(), Evidence: "未自动获取社区评分。可在 Cage 浏览器中打开补充检查；结果不替代底层隔离证据。"})
	s.mu.Lock()
	autoLanguage, location := s.preferences.AutoLanguage, s.location
	s.mu.Unlock()
	if autoLanguage && (location == nil || location.SuggestedLocale == "") {
		report.Checks = append(report.Checks, cage.SafetyCheck{ID: "language_suggestion", Title: "地区语言建议", State: cage.CheckUnknown, CheckedAt: time.Now(), Evidence: "当前地区没有可用的自动语言建议，已保留配置中的语言。", Remedy: "可在时区与语言设置中关闭自动，并手动选择"})
	}
	s.publishReport(report)
	s.mu.Lock()
	s.safety.ExitIP = hostIP
	if s.macosGuest != nil {
		s.safety.MacOSExitCandidate = nil
		if hostErr == nil && networkErr == nil && browserErr == nil {
			s.safety.MacOSExitCandidate = macOSStaticCandidate(cfg, state.RunID, hostIP, network.Egress, browser.Network)
		}
	}
	s.mu.Unlock()
	if err = report.Authorizes(state.RunID, cfg.Fingerprint(), time.Now()); err != nil {
		return err
	}
	s.safetyOp.Lock()
	defer s.safetyOp.Unlock()
	s.mu.Lock()
	paused = s.networkPaused
	s.mu.Unlock()
	if paused || ctx.Err() != nil {
		return errors.New("检测已取消，外网保持关闭")
	}
	state, err = client.Grant(ctx, report)
	if err != nil {
		return err
	}
	granted = true
	s.publishGate(state)
	s.mu.Lock()
	s.safety.VerifiedAt = time.Now()
	s.mu.Unlock()
	return nil
}

// Monitoring renews a short host-owned lease without replaying the disruptive
// failure drill. Lost supervision eventually closes existing guest connections.
func (s *Server) networkMonitorBlockedLocked() bool {
	job := s.jobs.snapshot()
	return s.closing || s.networkPaused || (job.Running && job.Kind != "environment-tools" && job.Kind != "claude-desktop-install")
}

func (s *Server) monitorSafety() {
	s.monitorSafetyWithProbe(cage.ProbeExitIP)
}

func (s *Server) monitorSafetyWithProbe(probe func(context.Context, cage.Config) (string, error)) {
	if !s.monitorSafetyOp.TryLock() {
		return
	}
	defer s.monitorSafetyOp.Unlock()
	timing := s.startMonitorTiming()
	defer s.recordMonitorTiming("safety", timing)
	s.mu.Lock()
	ready, busy, last := s.ready, s.networkMonitorBlockedLocked(), s.safety.VerifiedAt
	paused, closing := s.networkPaused, s.closing
	s.mu.Unlock()
	switch {
	case closing:
		timing.result = "closing"
		return
	case paused:
		timing.result = "network_paused"
		return
	case busy:
		timing.result = "job_running"
		return
	case !ready:
		timing.result = "environment_not_ready"
		return
	case time.Since(last) < 30*time.Second:
		timing.result = "not_due"
		return
	}
	client, err := s.gateClient()
	timing.stage("gate_client")
	if err != nil {
		timing.result = "gate_client_failed"
		s.mu.Lock()
		s.safety.Gate = nil
		s.safety.Error = err.Error()
		s.mu.Unlock()
		return
	}
	ctx, cancel := context.WithTimeout(s.lifetimeContext(), 45*time.Second)
	s.mu.Lock()
	s.monitorCancel = cancel
	s.mu.Unlock()
	defer func() { cancel(); s.mu.Lock(); s.monitorCancel = nil; s.mu.Unlock() }()
	state, err := client.Status(ctx)
	timing.stage("gate_status")
	if err != nil {
		timing.result = "gate_status_failed"
		s.mu.Lock()
		s.safety.Gate = nil
		s.safety.Error = err.Error()
		s.mu.Unlock()
		return
	}
	timing.until = state.Until
	s.publishGate(state)
	if state.Mode != cage.GateOpen || state.Report == nil {
		timing.result = "gate_not_open"
		// A packet helper can reconnect before the next GUI status tick. Its
		// gate stays closed until a fresh check, even though the helper never
		// exited and the new receipt is already available.
		recheckClosed := state.Reason == "检查或监测授权已到期" || s.windowsGuest != nil && state.Reason == windowsPacketDisconnectedReason
		if state.Mode == cage.GateClosed && recheckClosed && state.Report != nil {
			s.mu.Lock()
			pending := s.recovery.RunID != ""
			s.mu.Unlock()
			if !pending {
				s.scheduleNetworkRecovery(state)
			}
		}
		return
	}
	var issue string
	var recoverable bool
	cfg, err := loadConfig(s.configPath)
	if err != nil || cfg.Fingerprint() != state.Configuration {
		issue = "配置与运行环境不一致，请应用设置并重新检测"
	}
	if issue == "" {
		if err = state.Report.Authorizes(state.RunID, state.Configuration, time.Now()); err != nil {
			issue = err.Error()
		}
	}
	timing.stage("configuration_report")
	var ip string
	if issue == "" {
		ip, err = probe(ctx, cfg)
		timing.stage("exit_ip")
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil && s.deferTimedOutExitProbe(ctx, client, state) {
				timing.result = "exit_timeout_retry"
				return
			}
			issue = "宿主出口采样失败，外网已暂停：" + diagnostics.Redact(usererror.Display("读取宿主出口", err, "请重新检测。").Error())
			recoverable = true
		} else {
			result, message := s.checkExitPolicy(ctx, cfg, ip)
			timing.stage("exit_policy")
			if result != cage.CheckPassed {
				issue = message
				recoverable = result == cage.CheckUnknown
			}
		}
	}
	if issue == "" {
		guest, e := s.guestObservation(ctx, map[string]any{"stage": "environment", "timezone": cfg.Timezone})
		timing.stage("guest_settings")
		if e != nil || !s.guestSettingsOK(guest, cfg) {
			issue = "系统时区或语言的运行状态无法确认，请重新检测"
			recoverable = e != nil
		}
	}
	if issue == "" {
		browser, e := s.browserObservation(ctx, cfg, false)
		timing.stage("browser_settings")
		if e != nil {
			issue = "浏览器时区或语言的运行状态无法确认，请重新检测"
			recoverable = true
		} else if mismatch := s.browserSettingsIssue(browser, cfg); mismatch != "" {
			issue = "浏览器配置检查未通过：" + mismatch
		}
	}
	if issue != "" {
		timing.result = "observation_failed"
	}
	s.finishSafetyObservation(ctx, client, state, ip, issue, recoverable)
	timing.stage("finish_observation")
}

// Keep the already granted lease unchanged while retrying a timed-out exit
// query. The gate's own timer still closes all connections at its deadline.
func (s *Server) deferTimedOutExitProbe(ctx context.Context, client *cage.GateClient, observed cage.GateStatus) bool {
	s.safetyOp.Lock()
	defer s.safetyOp.Unlock()
	s.mu.Lock()
	blocked := s.networkMonitorBlockedLocked()
	s.mu.Unlock()
	if blocked || ctx.Err() != nil {
		return true
	}
	cfg, err := loadConfig(s.configPath)
	if err != nil || cfg.Fingerprint() != observed.Configuration {
		return false
	}
	current, err := client.Status(ctx)
	if err != nil {
		return true // The existing gate expires independently if it is unreachable.
	}
	if current.RunID != observed.RunID || current.Epoch != observed.Epoch || current.Mode != cage.GateOpen {
		return true // A newer transition owns this gate.
	}
	s.mu.Lock()
	s.safety.monitorError = "出口检测暂时超时，正在重试；现有联网许可未延长"
	s.mu.Unlock()
	return true
}

// Recheck the user operation and configuration at the actual lease write.
func (s *Server) finishSafetyMonitor(ctx context.Context, client *cage.GateClient, state cage.GateStatus, ip, issue string) {
	s.finishSafetyObservation(ctx, client, state, ip, issue, false)
}

func (s *Server) finishSafetyObservation(ctx context.Context, client *cage.GateClient, state cage.GateStatus, ip, issue string, recoverable bool) {
	timing := s.startMonitorTiming()
	timing.until = state.Until
	defer s.recordMonitorTiming("renewal", timing)
	s.safetyOp.Lock()
	defer s.safetyOp.Unlock()
	timing.stage("safety_lock")
	s.mu.Lock()
	busy := s.networkMonitorBlockedLocked()
	s.mu.Unlock()
	if busy {
		timing.result = "operation_cancelled"
		return
	} // A user operation cancelled this monitor; it owns the next transition.
	current, e := loadConfig(s.configPath)
	if e != nil || current.Fingerprint() != state.Configuration {
		issue = "运行配置已变化，请重新应用设置"
		recoverable = false
	}
	timing.stage("configuration")
	if issue != "" {
		timing.result = "observation_failed"
		closed, e := client.Close(context.Background(), issue)
		timing.stage("gate_close")
		if e == nil {
			s.publishGate(closed)
			if recoverable {
				s.scheduleNetworkRecovery(closed)
			}
		}
		s.mu.Lock()
		s.safety.Error = issue
		s.safety.ExitIP = ip
		s.safety.monitorError = issue
		s.mu.Unlock()
		return
	}
	renewed, err := client.Renew(ctx, state.Report.ID)
	timing.stage("gate_renew")
	if err != nil {
		timing.result = "renew_failed"
		closed, e := client.Close(context.Background(), "运行监测未能续期，请重新检测")
		timing.stage("gate_close")
		if e == nil {
			s.publishGate(closed)
			s.scheduleNetworkRecovery(closed)
		}
		s.mu.Lock()
		s.safety.Error = "运行监测未能续期，请重新检测"
		s.mu.Unlock()
		return
	}
	timing.result = "renewed"
	s.publishGate(renewed)
	s.mu.Lock()
	s.safety.VerifiedAt = time.Now()
	s.safety.Error = ""
	s.safety.monitorError = ""
	s.exitProbeError = ""
	s.safety.ExitIP = ip
	s.mu.Unlock()
}
