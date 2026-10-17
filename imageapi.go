package main

import (
	"errors"
	"math"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Hàm cho giao diện (và MCP – bước 31) quanh ảnh mẫu / màu pixel (bước 25). Cửa sổ app (overlay phủ màn hình, hoặc
// màn chính) che màn hình thật nên ẩn cửa sổ một chút trong lúc chụp.

// hideDelay – chờ cửa sổ ẩn hẳn (hiệu ứng của hệ điều hành) trước khi chụp.
const hideDelay = 250 * time.Millisecond

// TemplateShot – ảnh mẫu vừa chụp.
type TemplateShot struct {
	Png   string  `json:"png"`   // PNG base64 – ghi vào Step.Template
	Scale float64 `json:"scale"` // ghi vào Step.TemplateScale
	W     int     `json:"w"`     // cỡ theo toạ độ robotgo
	H     int     `json:"h"`
}

// FindResult – kết quả thử tìm ảnh mẫu trên màn hình hiện tại.
type FindResult struct {
	Found bool `json:"found"` // độ khớp ≥ ngưỡng của bước
	Score int  `json:"score"` // % khớp tốt nhất
	X     int  `json:"x"`     // khung chỗ giống nhất (toạ độ robotgo)
	Y     int  `json:"y"`
	W     int  `json:"w"`
	H     int  `json:"h"`
	Area  Rect `json:"area"` // vùng đã tìm
	Ms    int  `json:"ms"`   // thời gian tìm
}

// withHidden ẩn cửa sổ app, chạy f, rồi hiện lại – chỉ khi không chụp được mà bỏ qua cửa sổ app (macOS < 14, Linux…).
func (a *App) withHidden(f func() error) error {
	if a.ctx != nil && !canCaptureWithoutSelf() {
		runtime.WindowHide(a.ctx)
		defer runtime.WindowShow(a.ctx)
		time.Sleep(hideDelay)
	}
	return f()
}

// CaptureTemplate chụp vùng r (toạ độ robotgo, trên màn hình chính) làm ảnh mẫu.
func (a *App) CaptureTemplate(r Rect) (TemplateShot, error) {
	if err := captureArea(r, mainScreen()); err != nil {
		return TemplateShot{}, err
	}
	if r.W < 4 || r.H < 4 {
		return TemplateShot{}, errors.New("vùng khoanh quá nhỏ (cần ít nhất 4 × 4)")
	}
	var shot TemplateShot
	err := a.withHidden(func() error {
		img, err := robotgoScreen{}.Capture(r)
		if err != nil {
			return err
		}
		png, err := encodePNG(img)
		shot = TemplateShot{Png: png, Scale: imageScale(img, r), W: r.W, H: r.H}
		return err
	})
	return shot, err
}

// FindImage thử tìm ảnh mẫu của bước stepID trên màn hình hiện tại (một lần, không chờ) – nút "Thử tìm".
func (a *App) FindImage(sc Script, stepID string) (FindResult, error) {
	var s *Step
	for i := range sc.Steps {
		if sc.Steps[i].ID == stepID {
			s = &sc.Steps[i]
		}
	}
	if s == nil || s.Template == "" {
		return FindResult{}, errors.New("bước chưa có ảnh mẫu")
	}
	img, err := decodePNG(s.Template)
	if err != nil {
		return FindResult{}, err
	}
	e := newEngine(nil, sc, nil)
	e.screen = mainScreen()
	area := e.searchArea(*s)
	if err := captureArea(area, e.screen); err != nil {
		return FindResult{}, err
	}
	var res FindResult
	err = a.withHidden(func() error {
		shot, err := robotgoScreen{}.Capture(area)
		if err != nil {
			return err
		}
		start := time.Now()
		m, k := matchScaled(shot, area, toGray(img), s.TemplateScale)
		res = FindResult{
			Score: int(math.Round(m.Score * 100)), Area: area, Ms: int(time.Since(start).Milliseconds()),
			X: area.X + int(float64(m.X)/k), Y: area.Y + int(float64(m.Y)/k), W: int(float64(m.W) / k), H: int(float64(m.H) / k),
		}
		res.Found = res.Score >= s.Threshold
		return nil
	})
	return res, err
}

// PixelColorAt trả màu (#RRGGBB) đang hiện tại (x, y) – nút "Lấy màu" của bước chờ pixel.
func (a *App) PixelColorAt(x, y int) (string, error) {
	var hex string
	err := a.withHidden(func() error {
		c, err := robotgoScreen{}.Pixel(x, y)
		hex = hexColor(c)
		return err
	})
	return hex, err
}

// ScreenCaptureAllowed: đã có quyền chụp màn hình chưa (macOS Screen Recording; hệ khác luôn true).
func (a *App) ScreenCaptureAllowed() bool { return screenCaptureAllowed() }

// RequestScreenCapture mở hộp thoại xin quyền chụp màn hình của macOS (cấp xong thường phải mở lại app).
func (a *App) RequestScreenCapture() bool { return requestScreenCapture() }
