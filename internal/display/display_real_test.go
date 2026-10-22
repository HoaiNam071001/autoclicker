package display

import "testing"

// Chỉ đọc (không đổi gì trên máy): kiểm tra các cầu cgo / lệnh trả dữ liệu hợp lý.
func TestSystemReadOnly(t *testing.T) {
	var b Backend = System{}
	list, err := b.AudioDevices()
	if err != nil {
		t.Skipf("không đọc được thiết bị âm thanh: %v", err)
	}
	var out, in int
	for _, a := range list {
		t.Logf("%q %q input=%v default=%v", a.ID, a.Name, a.Input, a.Default)
		if a.Default {
			if a.Input {
				in++
			} else {
				out++
			}
		}
	}
	if len(list) > 0 && out > 1 || in > 1 {
		t.Errorf("mỗi loại tối đa một thiết bị mặc định: out=%d in=%d", out, in)
	}
	dark, err := b.Dark()
	t.Logf("dark=%v err=%v", dark, err)
	pct, err := b.Brightness()
	t.Logf("brightness=%d err=%v", pct, err)
}
