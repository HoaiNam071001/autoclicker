package store

import "testing"

func TestDeleteConsent(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("AppData", dir)
	if DeleteConfirmed("/k.json", "abc") {
		t.Fatal("chưa xác nhận mà đã được phép")
	}
	if err := ConfirmDelete("/k.json", "abc"); err != nil {
		t.Fatal(err)
	}
	if !DeleteConfirmed("/k.json", "abc") || DeleteConfirmed("/k.json", "khac") || DeleteConfirmed("/khac.json", "abc") {
		t.Fatal("xác nhận phải theo đúng kịch bản + đúng dấu vân tay")
	}
	for i := 0; i < maxConsents; i++ {
		ConfirmDelete("/k.json", string(rune('a'+i))+"x")
	}
	if DeleteConfirmed("/k.json", "abc") {
		t.Error("giữ quá maxConsents dấu vân tay")
	}
}
