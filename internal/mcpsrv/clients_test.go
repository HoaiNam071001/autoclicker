package mcpsrv

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddToClaudeDesktop(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Claude", "claude_desktop_config.json")
	if err := AddToClaudeDesktop(path, `C:\App\ghosthand.exe`); err != nil {
		t.Fatal(err)
	}
	// file đã có server khác + cài đặt khác → giữ nguyên
	os.WriteFile(path, []byte(`{"theme":"dark","mcpServers":{"khac":{"command":"x"},"ghosthand":{"command":"cu"},"autoclicker":{"command":"cu"}}}`), 0o644)
	if err := AddToClaudeDesktop(path, `C:\App\ghosthand.exe`); err != nil {
		t.Fatal(err)
	}
	var cfg struct {
		Theme      string
		McpServers map[string]struct {
			Command string
			Args    []string
		}
	}
	b, _ := os.ReadFile(path)
	if err := json.Unmarshal(b, &cfg); err != nil {
		t.Fatal(err)
	}
	s := cfg.McpServers["ghosthand"]
	if cfg.Theme != "dark" || cfg.McpServers["khac"].Command != "x" || s.Command != `C:\App\ghosthand.exe` || len(s.Args) != 1 || s.Args[0] != "--mcp" || len(cfg.McpServers) != 2 {
		t.Errorf("cấu hình: %s", b)
	}
	os.WriteFile(path, []byte(`{hỏng`), 0o644)
	if err := AddToClaudeDesktop(path, "x"); err == nil {
		t.Error("file hỏng phải báo lỗi, không ghi đè")
	}
}

func TestClientConfigs(t *testing.T) {
	cs := ClientConfigs(`C:\Program Files\AC\ghosthand.exe`, 0, "")
	if len(cs) != 5 || !strings.Contains(cs[1].Text, `"C:\Program Files\AC\ghosthand.exe" --mcp`) {
		t.Errorf("stdio: %+v", cs[1])
	}
	if !json.Valid([]byte(cs[0].Text)) {
		t.Error("JSON Claude Desktop")
	}
	cs = ClientConfigs("/a/ghosthand", 47800, "tok")
	if last := cs[len(cs)-1]; !strings.Contains(last.Text, "http://127.0.0.1:47800/mcp") || !strings.Contains(last.Text, "Bearer tok") {
		t.Errorf("http: %+v", last)
	}
}
