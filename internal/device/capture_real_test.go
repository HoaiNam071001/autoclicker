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
	t.Logf("màn chính %v, ảnh %v, tỷ lệ %.2f, DPI %.2f, chụp %v", main, full.Bounds().Size(), k, scr.DPIScale(main), time.Since(start))
	r := model.Rect{X: main.W / 3, Y: main.H / 3, W: 160, H: 60}
	part, err := scr.Capture(r)
	if err != nil {
		t.Fatal(err)
	}
	// Ảnh vùng phải trùng phần cắt cùng chỗ từ ảnh cả màn – sai gốc vùng (vd lệch DPI Windows) thì khác
	if diff := cropDiff(full, part, k, r, main); diff > 0.02 {
		t.Errorf("ảnh vùng %v khác phần cắt từ ảnh cả màn: %.1f%% điểm ảnh lệch", r, diff*100)
	}
	start = time.Now()
	m, sc := vision.MatchScaled(full, main, vision.NewTemplate(part), vision.MatchOptions{}, TemplateScale(part, r), scr.DPIScale(main))
	got := image.Pt(int(float64(m.X)/sc), int(float64(m.Y)/sc))
	t.Logf("tìm %v: thấy tại %v, khớp %.3f, %v", r, got, m.Score, time.Since(start))
	if m.Score < 0.95 {
		t.Errorf("khớp thấp %.3f", m.Score)
	}
}

// cropDiff – tỷ lệ điểm ảnh của part khác điểm cùng chỗ trong full (ảnh cả màn main, tỷ lệ ảnh / toạ độ k).
func cropDiff(full, part image.Image, k float64, r, main model.Rect) float64 {
	b, pb := full.Bounds(), part.Bounds()
	ox, oy := int(float64(r.X-main.X)*k), int(float64(r.Y-main.Y)*k)
	bad := 0
	for y := 0; y < pb.Dy(); y++ {
		for x := 0; x < pb.Dx(); x++ {
			r1, g1, b1, _ := part.At(pb.Min.X+x, pb.Min.Y+y).RGBA()
			r2, g2, b2, _ := full.At(b.Min.X+ox+x, b.Min.Y+oy+y).RGBA()
			if absDiff(r1, r2)+absDiff(g1, g2)+absDiff(b1, b2) > 3*8<<8 {
				bad++
			}
		}
	}
	return float64(bad) / float64(pb.Dx()*pb.Dy())
}

func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}
