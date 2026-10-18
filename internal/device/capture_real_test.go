package device

import (
	"image"
	"os"
	"testing"
	"time"

	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

// Chụp màn hình thật rồi tìm lại một mảnh vừa chụp – chỉ chạy khi AUTOCLICKER_REAL=1 (cần quyền ghi màn hình).
func TestRealCaptureAndFind(t *testing.T) {
	if os.Getenv("AUTOCLICKER_REAL") != "1" {
		t.Skip("đặt AUTOCLICKER_REAL=1 để chạy trên màn hình thật")
	}
	t.Logf("quyền ghi màn hình: %v", ScreenCaptureAllowed())
	scr, main := RobotgoScreen{}, MainScreen()
	start := time.Now()
	full, err := scr.Capture(main)
	if err != nil {
		t.Fatal(err)
	}
	k := vision.ImageScale(full, main)
	t.Logf("màn chính %v, ảnh %v, tỷ lệ %.2f, chụp %v", main, full.Bounds().Size(), k, time.Since(start))
	r := model.Rect{X: main.W / 3, Y: main.H / 3, W: 160, H: 60}
	part, err := scr.Capture(r)
	if err != nil {
		t.Fatal(err)
	}
	start = time.Now()
	m, sc := vision.MatchScaled(full, main, vision.ToGray(part), vision.ImageScale(part, r))
	got := image.Pt(int(float64(m.X)/sc), int(float64(m.Y)/sc))
	t.Logf("tìm %v: thấy tại %v, khớp %.3f, %v", r, got, m.Score, time.Since(start))
	if m.Score < 0.95 {
		t.Errorf("khớp thấp %.3f", m.Score)
	}
}
