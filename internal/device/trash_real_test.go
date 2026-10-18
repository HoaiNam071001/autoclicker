//go:build darwin

package device

import (
	"os"
	"path/filepath"
	"testing"
)

// Chuyển file thật vào Thùng rác của máy – chỉ khi AUTOCLICKER_REAL=1.
func TestRealMoveToTrash(t *testing.T) {
	if os.Getenv("AUTOCLICKER_REAL") != "1" {
		t.Skip("đặt AUTOCLICKER_REAL=1 để thử Thùng rác thật")
	}
	home, _ := os.UserHomeDir()
	dir, err := os.MkdirTemp(filepath.Join(home, "Documents"), "autoclicker-trash-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	p := filepath.Join(dir, "autoclicker-test-xoa.txt")
	os.WriteFile(p, []byte("x"), 0o644)
	if err := (SystemTrash{}).MoveToTrash([]string{p}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("file vẫn còn chỗ cũ")
	}
	if _, err := os.Stat(filepath.Join(home, ".Trash", "autoclicker-test-xoa.txt")); err != nil {
		t.Errorf("không thấy trong ~/.Trash: %v", err)
	}
}
