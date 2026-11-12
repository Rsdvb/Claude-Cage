package workbench

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

func validTailnetAuthURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host == "login.tailscale.com" && u.User == nil && u.Fragment == "" && (len(u.Path) > 3 && (u.Path[:3] == "/a/" || u.Path[:3] == "/c/"))
}
func (s *Server) tailnet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	action := r.PathValue("action")
	if action == "host-peers" || action == "host-open-service" || action == "host-close-service" {
		s.hostTailnet(w, r, action)
		return
	}
	if r.Method == http.MethodGet {
		action = "status"
	}
	if action == "install" {
		var input struct{}
		if !readJSON(w, r, &input) {
			return
		}
		s.launch(w, "tailscale-install", s.installTailnet)
		return
	}
	switch action {
	case "status", "connect", "login", "disconnect", "logout", "add-key", "remove-key", "open-service":
	default:
		http.NotFound(w, r)
		return
	}
	var input struct {
		Key  string `json:"key"`
		Host string `json:"host,omitempty"`
		Port int    `json:"port,omitempty"`
	}
	if r.Method != http.MethodGet && !readJSON(w, r, &input) {
		return
	}
	if len(input.Key) > 8192 {
		http.Error(w, "SSH 公钥过长", 400)
		return
	}
	if action == "open-service" && (len(input.Host) == 0 || len(input.Host) > 253 || input.Port < 1 || input.Port > 65535) {
		http.Error(w, "请选择私网设备并填写 1–65535 之间的端口", 400)
		return
	}
	s.mu.Lock()
	ready, token := s.ready, s.credentials.Token
	busy := s.jobs.snapshot().Running
	s.mu.Unlock()
	if !ready || busy {
		http.Error(w, "请先启动 Cage 环境，并等待当前操作完成", 409)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	body, _ := json.Marshal(struct {
		Action string `json:"action"`
		Key    string `json:"key,omitempty"`
		Host   string `json:"host,omitempty"`
		Port   int    `json:"port,omitempty"`
	}{action, input.Key, input.Host, input.Port})
	data, err := s.guestAPI(ctx, token, http.MethodPost, "/tailnet", body)
	if err != nil {
		http.Error(w, "Tailscale 操作未完成。请更新环境组件后重试；原设备身份和 SSH 公钥会保留。", 502)
		return
	}
	var result map[string]json.RawMessage
	if json.Unmarshal(data, &result) != nil {
		http.Error(w, "Tailscale 组件返回了无效状态", 502)
		return
	}
	if b := result["tailnet_error"]; len(b) > 0 {
		var reason string
		_ = json.Unmarshal(b, &reason)
		http.Error(w, reason, 502)
		return
	}
	var authURL string
	_ = json.Unmarshal(result["auth_url"], &authURL)
	if authURL != "" && !validTailnetAuthURL(authURL) {
		http.Error(w, "Tailscale 登录地址无效，请重新发起登录", 502)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

var errTailnetNativeInstall = errors.New("请在设置中更新环境组件，Tailscale 随客机组件一同安装")
