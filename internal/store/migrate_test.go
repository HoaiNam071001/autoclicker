package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMoveOldDirsKeepsExisting(t *testing.T) {
	tmp := t.TempDir()
	old, neu := filepath.Join(tmp, "Auto Clicker"), filepath.Join(tmp, "Ghosthand")
	os.MkdirAll(old, 0o755)
	os.WriteFile(filepath.Join(old, "a.json"), []byte("{}"), 0o644)
	if got := moveOldDirs([]dirMove{{old, neu}, {filepath.Join(tmp, "khong-co"), filepath.Join(tmp, "x")}}); len(got) != 1 {
		t.Fatalf("chuyển %v, muốn 1 cặp", got)
	}
	if _, err := os.Stat(filepath.Join(neu, "a.json")); err != nil {
		t.Fatal("file phải sang thư mục mới")
	}
	// Mới đã có → không chuyển nữa (không ghi đè dữ liệu bản mới)
	os.MkdirAll(old, 0o755)
	if got := moveOldDirs([]dirMove{{old, neu}}); len(got) != 0 {
		t.Errorf("mới đã có vẫn chuyển: %v", got)
	}
}

func TestRewritePathsWholeDirOnly(t *testing.T) {
	tmp := t.TempDir()
	old, neu := filepath.Join(tmp, "Auto Clicker"), filepath.Join(tmp, "Ghosthand")
	other := filepath.Join(tmp, "Auto Clicker 2", "b.json")
	settings := `{"recent":[` + jq(filepath.Join(old, "a.json")) + `,` + jq(other) + `],"storage":{"shotsDir":` + jq(old) + `}}`
	os.WriteFile(filepath.Join(tmp, "settings.json"), []byte(settings), 0o644)
	if err := rewritePaths(tmp, []dirMove{{old, neu}}); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(tmp, "settings.json"))
	got := string(b)
	for _, want := range []string{jq(filepath.Join(neu, "a.json")), jq(other), `"shotsDir":` + jq(neu)} {
		if !strings.Contains(got, want) {
			t.Errorf("thiếu %s trong %s", want, got)
		}
	}
}

func TestMoveRunHistoryFollowsScriptPath(t *testing.T) {
	tmp := t.TempDir()
	cfg := filepath.Join(tmp, "cfg")
	m := dirMove{filepath.Join(tmp, "Auto Clicker"), filepath.Join(tmp, "Ghosthand")}
	os.MkdirAll(m.to, 0o755)
	os.WriteFile(filepath.Join(m.to, "dang-nhap.json"), []byte(`{"name":"Đăng nhập"}`), 0o644)
	oldRuns := RunsDir(cfg, "Đăng nhập", filepath.Join(m.from, "dang-nhap.json"))
	os.MkdirAll(oldRuns, 0o755)
	moveRunHistory(cfg, m)
	if _, err := os.Stat(RunsDir(cfg, "Đăng nhập", filepath.Join(m.to, "dang-nhap.json"))); err != nil {
		t.Error("lịch sử phải theo đường dẫn mới")
	}
}

func jq(s string) string { return `"` + jsonInner(s) + `"` }
