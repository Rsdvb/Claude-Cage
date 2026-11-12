package workbench

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sylvinus/agent-vm/internal/containerenv"
	"github.com/sylvinus/agent-vm/internal/host"
)

var tailnetIPv4 = netip.MustParsePrefix("100.64.0.0/10")

type hostTailnetPeer struct {
	ID       string       `json:"ID"`
	HostName string       `json:"HostName"`
	DNSName  string       `json:"DNSName"`
	IPs      []netip.Addr `json:"TailscaleIPs"`
}

type hostTailnetService struct {
	Device string `json:"device"`
	Port   int    `json:"port"`
}

type hostTailnetForward struct {
	hostTailnetService
	Address string `json:"address"`
	cancel  context.CancelFunc
	done    chan struct{}
}

type hostTailnetForwards struct {
	mu      sync.Mutex
	entries map[hostTailnetService]*hostTailnetForward
}

func (f *hostTailnetForwards) prune() {
	for key, entry := range f.entries {
		select {
		case <-entry.done:
			delete(f.entries, key)
		default:
		}
	}
}

func (f *hostTailnetForwards) list() []hostTailnetForward {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.prune()
	result := make([]hostTailnetForward, 0, len(f.entries))
	for _, entry := range f.entries {
		result = append(result, *entry)
	}
	return result
}

func (f *hostTailnetForwards) closeAll() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for key, entry := range f.entries {
		entry.cancel()
		delete(f.entries, key)
	}
}

func hostTailscaleCLI() (string, error) {
	path, err := exec.LookPath("tailscale")
	if err == nil && filepath.IsAbs(path) {
		return path, nil
	}
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(os.Getenv("ProgramFiles"), "Tailscale", "tailscale.exe")
	case "darwin":
		path = "/Applications/Tailscale.app/Contents/MacOS/Tailscale"
	}
	if filepath.IsAbs(path) {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}
	return "", errors.New("请在宿主机安装并登录 Tailscale")
}

func hostTailnetPeers(ctx context.Context, cli string) ([]hostTailnetPeer, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	data, err := hostTailscaleCommand(ctx, cli, "status", "--json").Output()
	var status struct {
		BackendState string
		Peer         map[string]hostTailnetPeer
	}
	if err != nil || json.Unmarshal(data, &status) != nil || status.BackendState != "Running" {
		return nil, errors.New("宿主机 Tailscale 尚未连接，请登录后重试")
	}
	peers := make([]hostTailnetPeer, 0, len(status.Peer))
	for _, peer := range status.Peer {
		if peer.ID != "" && peer.ip().IsValid() {
			peers = append(peers, peer)
		}
	}
	slices.SortFunc(peers, func(a, b hostTailnetPeer) int {
		if byName := strings.Compare(a.HostName, b.HostName); byName != 0 {
			return byName
		}
		return strings.Compare(a.ID, b.ID)
	})
	return peers, nil
}

func (p hostTailnetPeer) ip() netip.Addr {
	for _, ip := range p.IPs {
		if tailnetIPv4.Contains(ip) {
			return ip
		}
	}
	return netip.Addr{}
}

func (s *Server) hostTailnet(w http.ResponseWriter, r *http.Request, action string) {
	var service hostTailnetService
	if !readJSON(w, r, &service) {
		return
	}
	if s.busy() {
		http.Error(w, "请等待当前 Cage 操作完成", http.StatusConflict)
		return
	}
	// Serialize creation with environment stop or replacement.
	s.op.Lock()
	defer s.op.Unlock()
	s.mu.Lock()
	ready := s.ready && !s.closing
	paused := s.networkPaused
	lifetime := s.lifetimeContextLocked()
	s.mu.Unlock()
	if !ready || s.busy() || s.clientOnly || s.windowsGuest != nil || containerenv.Configured(s.profile) {
		http.Error(w, "请启动本机 Linux 或 macOS 客机后再使用宿主机转发", http.StatusConflict)
		return
	}
	if action != "host-peers" && (service.Device == "" || len(service.Device) > 128 || service.Port < 1 || service.Port > 65535) {
		http.Error(w, "请选择设备并填写 1–65535 之间的 TCP 端口", http.StatusBadRequest)
		return
	}
	if action == "host-open-service" && paused {
		http.Error(w, "请先恢复 Cage 联网，再创建宿主机转发", http.StatusConflict)
		return
	}
	var peers []hostTailnetPeer
	var warning string
	if action == "host-close-service" {
		s.hostTailnetForwards.mu.Lock()
		if entry := s.hostTailnetForwards.entries[service]; entry != nil {
			entry.cancel()
			delete(s.hostTailnetForwards.entries, service)
		}
		s.hostTailnetForwards.mu.Unlock()
	} else {
		cli, err := hostTailscaleCLI()
		if err == nil {
			peers, err = hostTailnetPeers(r.Context(), cli)
		}
		if err != nil {
			if action == "host-open-service" {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			warning = err.Error()
		}
		if action == "host-open-service" {
			index := slices.IndexFunc(peers, func(peer hostTailnetPeer) bool { return peer.ID == service.Device })
			if index < 0 {
				http.Error(w, "该设备不在宿主机的 Tailscale 私网设备列表中", http.StatusBadRequest)
				return
			}
			epoch, unlocked := s.account.accessFor(s.hasTrustedView(r))
			if !unlocked {
				http.Error(w, "Cage 已锁定，请先登录", http.StatusLocked)
				return
			}
			if err := s.openHostTailnetForward(r.Context(), lifetime, epoch, cli, peers[index].ip(), service); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Peers    []hostTailnetPeer    `json:"peers,omitempty"`
		Forwards []hostTailnetForward `json:"forwards"`
		Warning  string               `json:"warning,omitempty"`
	}{peers, s.hostTailnetForwards.list(), warning})
}

func (s *Server) openHostTailnetForward(request, lifetime, account context.Context, cli string, ip netip.Addr, service hostTailnetService) error {
	forwards := &s.hostTailnetForwards
	forwards.mu.Lock()
	forwards.prune()
	exists, count := forwards.entries[service] != nil, len(forwards.entries)
	forwards.mu.Unlock()
	if exists {
		return nil
	}
	if count >= 16 {
		return errors.New("已创建 16 条宿主机连接，请关闭不再使用的连接")
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return errors.New("无法创建宿主机转发监听")
	}
	ctx, cancel := context.WithCancel(lifetime)
	stopAccount := context.AfterFunc(account, cancel)
	failed := true
	defer func() {
		if failed {
			stopAccount()
			cancel()
			listener.Close()
		}
	}()
	forward := "127.0.0.1:0:" + listener.Addr().String()
	cmd := s.sshCommandWithForwarding(ctx, "-v", "-N", "-T", "-o", "ExitOnForwardFailure=yes", "-R", forward, s.sshTarget())
	if cmd.Err != nil {
		return errors.New("无法创建客机 SSH 转发")
	}
	cmd.Args = slices.DeleteFunc(cmd.Args, func(arg string) bool { return arg == "-q" })
	for i, arg := range cmd.Args {
		if arg == "ClearAllForwardings=yes" {
			cmd.Args[i] = "ClearAllForwardings=no"
		}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil || cmd.Start() != nil {
		return errors.New("无法启动客机 SSH 转发")
	}

	allocated := make(chan int, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer stopAccount()
		defer cancel()
		defer listener.Close()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if port := allocatedTailnetPort(scanner.Text(), listener.Addr().String()); port != 0 {
				select {
				case allocated <- port:
				default:
				}
			}
		}
		_ = cmd.Wait()
	}()
	startup, finish := context.WithTimeout(request, 8*time.Second)
	defer finish()
	var guestPort int
	select {
	case guestPort = <-allocated:
	case <-done:
		return errors.New("客机拒绝了 SSH 端口转发")
	case <-startup.Done():
		return errors.New("等待客机转发端口超时")
	}
	if err := errors.Join(startup.Err(), ctx.Err()); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.networkPaused || !s.ready || s.closing {
		return errors.New("Cage 环境已暂停或停止")
	}
	forwards.mu.Lock()
	defer forwards.mu.Unlock()
	if forwards.entries == nil {
		forwards.entries = make(map[hostTailnetService]*hostTailnetForward)
	}
	forwards.entries[service] = &hostTailnetForward{
		hostTailnetService: service,
		Address:            net.JoinHostPort("127.0.0.1", strconv.Itoa(guestPort)),
		cancel:             cancel,
		done:               done,
	}
	failed = false
	go serveHostTailnetForward(ctx, listener, cli, ip.String(), strconv.Itoa(service.Port))
	return nil
}

func allocatedTailnetPort(line, target string) int {
	rest, ok := strings.CutPrefix(line, "Allocated port ")
	if !ok {
		return 0
	}
	port, destination, ok := strings.Cut(rest, " for remote forward to ")
	if !ok || destination != target {
		return 0
	}
	value, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0
	}
	return int(value)
}

func serveHostTailnetForward(ctx context.Context, listener net.Listener, cli, address, port string) {
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	defer listener.Close()
	slots := make(chan struct{}, 32)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		select {
		case slots <- struct{}{}:
			go func() {
				defer func() { <-slots }()
				defer conn.Close()
				stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
				defer stop()
				runHostTailnetStream(ctx, conn, cli, address, port)
			}()
		default:
			_ = conn.Close()
		}
	}
}

func runHostTailnetStream(ctx context.Context, conn net.Conn, cli, address, port string) {
	cmd := hostTailscaleCommand(ctx, cli, "nc", address, port)
	input, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	cmd.Stdout = conn
	if cmd.Start() != nil {
		input.Close()
		return
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer input.Close()
		_, _ = io.Copy(input, conn)
	}()
	_ = cmd.Wait()
	_ = conn.Close()
	<-done
}

func hostTailscaleCommand(ctx context.Context, cli string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, cli, args...)
	if runtime.GOOS == "darwin" {
		cmd.Env = append(cmd.Environ(), "TAILSCALE_BE_CLI=1")
	}
	host.HideCommandWindow(cmd)
	return cmd
}
