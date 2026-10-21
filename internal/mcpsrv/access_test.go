package mcpsrv

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"ghosthand/internal/store"
)

func post(t *testing.T, h http.Handler, url, client, body string) string {
	t.Helper()
	req := httptest.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if client != "" {
		req.Header.Set(ClientHeader, client)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w.Body.String()
}

func TestToolAccessPerClient(t *testing.T) {
	back := &fakeBackend{dir: t.TempDir(), revs: store.NewRevisions()}
	back.settings = store.MCPSettings{Enabled: true, OffTools: []string{"screenshot"},
		Clients: map[string][]string{"claude-code": {}, "cursor": {"screenshot", "run_script"}}}
	h := New(back).Handler()
	list := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
	shot := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"screenshot","arguments":{}}}`

	// mặc định chung (không rõ client) tắt screenshot; Claude Code ghi đè bật lại; Cursor tắt thêm run_script
	if got := post(t, h, "/mcp", "", list); strings.Contains(got, `"screenshot"`) || !strings.Contains(got, `"run_script"`) {
		t.Errorf("mặc định: %s", got)
	}
	if got := post(t, h, "/mcp", "claude-code", list); !strings.Contains(got, `"screenshot"`) {
		t.Errorf("claude-code phải thấy screenshot: %s", got)
	}
	if got := post(t, h, "/mcp", "cursor", list); strings.Contains(got, `"run_script"`) {
		t.Errorf("cursor không được thấy run_script: %s", got)
	}
	if got := post(t, h, "/mcp", "cursor", shot); !strings.Contains(got, "tool_disabled") {
		t.Errorf("gọi tool tắt phải bị từ chối: %s", got)
	}
	// client nối HTTP khai tên qua ?client=
	if got := post(t, h, "/mcp?client=cursor", "", list); strings.Contains(got, `"run_script"`) {
		t.Errorf("?client=: %s", got)
	}
}

func TestToolCatalogAndSeenClients(t *testing.T) {
	back := &fakeBackend{dir: t.TempDir(), revs: store.NewRevisions()}
	srv := New(back)
	var names []string
	for _, ti := range srv.Tools() {
		if ti.Kind == "" || ti.Desc == "" {
			t.Errorf("%s thiếu nhóm / mô tả", ti.Name)
		}
		names = append(names, ti.Name)
	}
	for _, n := range []string{"run_script", "screenshot", "update_script", "list_scripts", "open_in_app"} {
		if !slices.Contains(names, n) {
			t.Errorf("bảng tool thiếu %s", n)
		}
	}
	post(t, srv.Handler(), "/mcp", "", `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"cursor-vscode","version":"1"}}}`)
	if c := srv.Clients(); len(c) != 1 || c[0] != "cursor-vscode" {
		t.Errorf("Clients = %v", c)
	}
}
