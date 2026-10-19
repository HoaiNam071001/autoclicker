package store

import "testing"

func TestDangerConsent(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("AppData", dir)
	if DangerConfirmed("/k.json", "abc") {
		t.Fatal("chưa xác nhận mà đã được phép")
	}
	if err := ConfirmDanger("/k.json", "abc"); err != nil {
		t.Fatal(err)
	}
	if !DangerConfirmed("/k.json", "abc") || DangerConfirmed("/k.json", "khac") || DangerConfirmed("/khac.json", "abc") {
		t.Fatal("xác nhận phải theo đúng kịch bản + đúng dấu vân tay")
	}
	for i := 0; i < maxConsents; i++ {
		ConfirmDanger("/k.json", string(rune('a'+i))+"x")
	}
	if DangerConfirmed("/k.json", "abc") {
		t.Error("giữ quá maxConsents dấu vân tay")
	}
}
