package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveDir(t *testing.T) {
	tmp := t.TempDir()
	def := filepath.Join(tmp, "Documents", "Auto Clicker")
	if got, err := resolveDir("", def); err != nil || got != def {
		t.Fatalf("trống = %q %v, muốn mặc định %q", got, err, def)
	}
	if _, err := os.Stat(def); err != nil {
		t.Errorf("mặc định chưa được tạo: %v", err)
	}
	if got, err := resolveDir(tmp, def); err != nil || got != tmp {
		t.Errorf("thư mục đã chọn = %q %v", got, err)
	}
	gone := filepath.Join(tmp, "da-xoa")
	if _, err := resolveDir(gone, def); err == nil || !strings.Contains(err.Error(), "Về mặc định") {
		t.Errorf("thư mục không còn phải báo lỗi gợi ý về mặc định, được %v", err)
	}
	if _, err := os.Stat(gone); err == nil {
		t.Error("không được tự tạo lại thư mục đã chọn")
	}
}

func TestAIScriptsDirFallsBackToScripts(t *testing.T) {
	tmp := t.TempDir()
	d := StorageDefaults{ScriptsDir: filepath.Join(tmp, "def")}
	mine := filepath.Join(tmp, "mine")
	os.Mkdir(mine, 0o755)
	if got, _ := (StorageSettings{}).aiScripts(d); got != d.ScriptsDir {
		t.Errorf("không chọn gì: AI = %q, muốn %q", got, d.ScriptsDir)
	}
	if got, _ := (StorageSettings{ScriptsDir: mine}).aiScripts(d); got != mine {
		t.Errorf("AI trống phải theo thư mục kịch bản: %q", got)
	}
	ai := filepath.Join(tmp, "ai")
	os.Mkdir(ai, 0o755)
	if got, _ := (StorageSettings{ScriptsDir: mine, AIScriptsDir: ai}).aiScripts(d); got != ai {
		t.Errorf("AI riêng = %q, muốn %q", got, ai)
	}
}

func TestDefaultDirsNamedAutoClicker(t *testing.T) {
	d, err := DefaultDirs()
	if err != nil {
		t.Skip(err)
	}
	if filepath.Base(d.ScriptsDir) != "Auto Clicker" || filepath.Base(d.ShotsDir) != "Auto Clicker" {
		t.Errorf("mặc định %+v phải là …/Auto Clicker", d)
	}
}
