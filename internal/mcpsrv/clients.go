package mcpsrv

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"ghosthand/internal/store"
)

// ServerName – tên server trong cấu hình của client AI.
const ServerName = "ghosthand"

// ClientConfig – đoạn cấu hình mẫu cho một client AI (Cài đặt → MCP).
type ClientConfig struct {
	Client string `json:"client"` // Claude Desktop, Claude Code…
	Where  string `json:"where"`  // dán vào đâu
	Text   string `json:"text"`
}

// ClientConfigs – cấu hình stdio (exe --mcp) cho từng client; httpPort > 0 (API cho chương trình khác đang bật) thì
// thêm cấu hình HTTP cho client nối thẳng URL.
func ClientConfigs(exe string, httpPort int, token string) []ClientConfig {
	stdio, _ := json.MarshalIndent(map[string]any{"mcpServers": map[string]any{ServerName: stdioEntry(exe)}}, "", "  ")
	q := strconv.Quote(exe)
	out := []ClientConfig{
		{"Claude Desktop", "claude_desktop_config.json (Settings → Developer → Edit Config) – hoặc bấm nút Thêm ở trên", string(stdio)},
		{"Claude Code", "chạy trong terminal", "claude mcp add " + ServerName + " -- \"" + exe + "\" --mcp"},
		{"Cursor / VS Code Copilot", "~/.cursor/mcp.json hoặc .vscode/mcp.json (VS Code đổi mcpServers → servers)", string(stdio)},
		{"Codex CLI", "~/.codex/config.toml", "[mcp_servers." + ServerName + "]\ncommand = " + q + "\nargs = [\"--mcp\"]"},
		{"Gemini CLI", "~/.gemini/settings.json", string(stdio)},
	}
	if httpPort > 0 {
		url := fmt.Sprintf("http://127.0.0.1:%d/mcp", httpPort)
		h, _ := json.MarshalIndent(map[string]any{"mcpServers": map[string]any{ServerName: map[string]any{
			"type": "http", "url": url, "headers": map[string]string{"Authorization": "Bearer " + token}}}}, "", "  ")
		out = append(out, ClientConfig{"HTTP (client nối URL)", "cổng của API cho chương trình khác", string(h)})
	}
	return out
}

func stdioEntry(exe string) map[string]any {
	return map[string]any{"command": exe, "args": []string{"--mcp"}}
}

// ClaudeDesktopConfig – đường dẫn claude_desktop_config.json (Windows %APPDATA%\Claude, macOS ~/Library/Application Support/Claude).
func ClaudeDesktopConfig() (string, error) {
	dir, err := os.UserConfigDir()
	return filepath.Join(dir, "Claude", "claude_desktop_config.json"), err
}

// AddToClaudeDesktop thêm / sửa mục ghosthand trong mcpServers của file cấu hình, giữ nguyên phần còn lại; bỏ mục tên cũ
// "autoclicker" (bước 45 – exe đã đổi tên, để lại thì Claude Desktop báo lỗi server; xoá cùng store/migrate.go).
func AddToClaudeDesktop(path, exe string) error {
	cfg := map[string]any{}
	b, err := os.ReadFile(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
	case err != nil:
		return err
	case len(b) > 0:
		if err := json.Unmarshal(b, &cfg); err != nil {
			return fmt.Errorf("%s không phải JSON hợp lệ – sửa tay trước: %w", path, err)
		}
	}
	servers, _ := cfg["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}
	delete(servers, "autoclicker")
	servers[ServerName] = stdioEntry(exe)
	cfg["mcpServers"] = servers
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return store.WriteFileAtomic(path, out)
}
