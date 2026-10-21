package device

import (
	"os"
	"testing"
)

// Đặt chuột ở góc / giữa từng màn rồi đọc lại phải đúng điểm đó – bắt lệch DPI (Windows 125 %: chuột và khung màn
// phải cùng pixel thật). Chỉ chạy khi GHOSTHAND_REAL=1 (chuột thật di chuyển).
func TestRealCursorRoundTrip(t *testing.T) {
	if os.Getenv("GHOSTHAND_REAL") != "1" {
		t.Skip("đặt GHOSTHAND_REAL=1 để chạy với chuột thật")
	}
	ox, oy := Location()
	defer moveCursor(ox, oy)
	for i, d := range Displays() {
		t.Logf("màn %d: %v", i+1, d)
		for _, p := range [][2]int{{d.X + 1, d.Y + 1}, {d.X + d.W/2, d.Y + d.H/2}, {d.X + d.W - 2, d.Y + d.H - 2}} {
			if err := moveCursor(p[0], p[1]); err != nil {
				t.Fatal(err)
			}
			if x, y := Location(); x != p[0] || y != p[1] {
				t.Errorf("màn %d: đặt chuột (%d, %d) đọc lại (%d, %d)", i+1, p[0], p[1], x, y)
			}
		}
	}
}
