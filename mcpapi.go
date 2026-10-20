package main

import (
	"os"

	"autoclicker/internal/device"
	"autoclicker/internal/mcpsrv"
	"autoclicker/internal/model"
	"autoclicker/internal/session"
	"autoclicker/internal/store"
)

// MCP cho AI (bước 31, internal/mcpsrv): route /mcp trên API local (cả kênh CLI – `autoclicker --mcp` chuyển tiếp
// vào đó – lẫn cổng cho chương trình khác). Công tắc trong Cài đặt, bật mặc định.

// startMCP – gọi trong startAPI, trước khi mở cổng.
func (a *App) startMCP() {
	a.mcp = mcpsrv.New(mcpBackend{a})
	a.apiSrv.Handle("/mcp", a.mcp.Handler())
}

// SetMCPSettings bật / tắt MCP.
// Chuyển về Tắt / Chỉ soạn thì huỷ ngay lượt AI đang chạy.
func (a *App) SetMCPSettings(m store.MCPSettings) error {
	if err := store.SetMCP(m); err != nil {
		return err
	}
	a.stopAIRun(m)
	return nil
}

// MCPLog – các lệnh AI đã gọi gần đây (mới nhất trước) cho Cài đặt.
func (a *App) MCPLog() []mcpsrv.Call {
	if a.mcp == nil {
		return []mcpsrv.Call{}
	}
	return a.mcp.Log()
}

// MCPClients – đoạn cấu hình mẫu cho từng client AI (exe của app + --mcp; HTTP khi API cho chương trình khác bật).
func (a *App) MCPClients() ([]mcpsrv.ClientConfig, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	port, api := 0, store.CurrentSettings().API
	if api.Enabled {
		port = api.Port
	}
	return mcpsrv.ClientConfigs(exe, port, api.Token), nil
}

// ClaudeDesktopConfigPath – file cấu hình của Claude Desktop (hiện trong câu hỏi trước khi ghi).
func (a *App) ClaudeDesktopConfigPath() (string, error) { return mcpsrv.ClaudeDesktopConfig() }

// AddToClaudeDesktop ghi mục autoclicker vào cấu hình Claude Desktop (frontend đã hỏi người dùng). Trả đường dẫn file.
func (a *App) AddToClaudeDesktop() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	path, err := mcpsrv.ClaudeDesktopConfig()
	if err != nil {
		return "", err
	}
	return path, mcpsrv.AddToClaudeDesktop(path, exe)
}

// AIScripts – trong danh sách mở gần đây, file nào do AI tạo (chip AI ở thư viện).
func (a *App) AIScripts() []string {
	out := []string{}
	for _, p := range store.RecentScripts() {
		if sc, err := store.ReadScriptFile(p); err == nil && sc.Meta != nil && sc.Meta.CreatedBy == model.CreatedByAI {
			out = append(out, p)
		}
	}
	return out
}

// mcpBackend – mcpsrv.Backend nối vào App.
type mcpBackend struct{ a *App }

func (b mcpBackend) Enabled() bool               { return store.CurrentSettings().MCP.Enabled }
func (b mcpBackend) Recent() []string            { return store.RecentScripts() }
func (b mcpBackend) ScriptsDir() (string, error) { return store.ScriptsDir() }

func (b mcpBackend) Editor() mcpsrv.Editor {
	b.a.mu.Lock()
	defer b.a.mu.Unlock()
	return mcpsrv.Editor{Path: b.a.editPath, Dirty: b.a.dirty}
}

func (b mcpBackend) ReadScript(path string) (model.Script, int, error) {
	return apiBackend(b).ReadScript(path)
}

func (b mcpBackend) WriteScript(path string, sc model.Script, rev int) (int, error) {
	rev, err := b.a.writeScriptExternal(path, sc, rev, session.SourceAI)
	if err == nil {
		store.EditRecent(func(l []string) []string { return store.PushRecent(l, path) })
	}
	return rev, err
}

// OpenInApp – frontend mở file (scriptFile.ts, nghe "mcp:open") và đưa cửa sổ lên trước.
func (b mcpBackend) OpenInApp(path string) error {
	b.a.showMain()
	b.a.emit("mcp:open", path)
	return nil
}

func (b mcpBackend) Displays() []mcpsrv.Display {
	main := device.MainScreen()
	out := []mcpsrv.Display{}
	for _, r := range device.Displays() {
		out = append(out, mcpsrv.Display{Rect: r, Scale: device.DPIScale(r), Main: r == main})
	}
	return out
}
