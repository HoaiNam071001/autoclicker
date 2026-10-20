package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"autoclicker/internal/api"
	"autoclicker/internal/mcpsrv"
	"autoclicker/internal/model"
	"autoclicker/internal/store"
)

// emptyBackend – app không có kịch bản nào.
type emptyBackend struct{ dir string }

func (b emptyBackend) Enabled() bool               { return true }
func (b emptyBackend) Recent() []string            { return nil }
func (b emptyBackend) Editor() mcpsrv.Editor       { return mcpsrv.Editor{} }
func (b emptyBackend) ScriptsDir() (string, error) { return b.dir, nil }
func (b emptyBackend) ReadScript(string) (model.Script, int, error) {
	return model.Script{}, 0, os.ErrNotExist
}
func (b emptyBackend) WriteScript(string, model.Script, int) (int, error) {
	return 0, errors.New("không ghi")
}
func (b emptyBackend) OpenInApp(string) error     { return nil }
func (b emptyBackend) Displays() []mcpsrv.Display { return nil }

func TestMCPProxy(t *testing.T) {
	srv := api.New(nil, api.NewHub())
	srv.Handle("/mcp", mcpsrv.New(emptyBackend{t.TempDir()}).Handler())
	l, err := srv.Listen(0, "tok")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	c := &Client{Info: func() (store.APIInfo, error) { return store.APIInfo{Port: l.Port, Token: "tok"}, nil }, HTTP: http.DefaultClient,
		Launch: func() error { return errors.New("không được mở") }}

	in := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":"x","method":"tools/call","params":{"name":"list_scripts","arguments":{}}}`,
	}, "\n")
	var out, errOut bytes.Buffer
	if code := c.MCP(strings.NewReader(in), &out, &errOut); code != ExitOK {
		t.Fatal(code)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 { // notification không có trả lời
		t.Fatalf("stdout:\n%s\nstderr: %s", out.String(), errOut.String())
	}
	var list struct {
		ID     int `json:"id"`
		Result struct{ Tools []struct{ Name string } }
	}
	if err := json.Unmarshal([]byte(lines[1]), &list); err != nil || list.ID != 2 || len(list.Result.Tools) < 8 {
		t.Errorf("tools/list: %s", lines[1])
	}
	if !strings.Contains(lines[2], `"id":"x"`) || !strings.Contains(lines[2], "scripts") {
		t.Errorf("tools/call: %s", lines[2])
	}
}

func TestMCPNoApp(t *testing.T) {
	c := &Client{Info: func() (store.APIInfo, error) { return store.APIInfo{}, os.ErrNotExist }, HTTP: http.DefaultClient,
		Launch: func() error { return nil }, WaitApp: 30 * time.Millisecond, poll: 10 * time.Millisecond}
	var out, errOut bytes.Buffer
	c.MCP(strings.NewReader(`{"jsonrpc":"2.0","id":7,"method":"tools/list"}`+"\n"+`{"jsonrpc":"2.0","method":"notifications/x"}`), &out, &errOut)
	if s := strings.TrimSpace(out.String()); !strings.Contains(s, `"id":7`) || !strings.Contains(s, "không nối được app") || strings.Count(s, "\n") != 0 {
		t.Errorf("không có app → lỗi JSON-RPC cho request, bỏ qua notification: %s", s)
	}
	if !IsMCP([]string{"--mcp"}) || IsMCP([]string{"run"}) {
		t.Error("IsMCP")
	}
}
