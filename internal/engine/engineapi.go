package engine

import (
	"errors"
	"image"
	"math"
	"strings"
	"time"

	"autoclicker/internal/device"
	"autoclicker/internal/model"
	"autoclicker/internal/vision"
	"autoclicker/internal/web"
)

// API của engine cho app (RunScript, Thử node, Thử tìm) và sau này MCP – không cần biết trường nội bộ.

// Options – cấu hình một engine. Giá trị 0 = không dùng (Thử node / Thử tìm chỉ cần Screen + Displays).
type Options struct {
	Screen     model.Rect   // khung màn hình chính: kịch bản không có vùng thì neo theo khung này
	Displays   []model.Rect // khung các màn (tìm "cả màn", kiểm tra vùng chụp); nil = không kiểm tra
	OnProgress func(Progress)
	OnNode     func(model.StepResult) // mỗi node chạy xong
	Debug      bool                   // dừng ở node có Breakpoint (chạy từ app)
	StepMode   bool                   // dừng trước mỗi node (chạy từng node)
	KeepShots  bool                   // chụp màn hình lúc node thất bại (lịch sử chạy)
	Browser    web.Browser            // node web + node Mở bằng Chrome của app; nil = không có (node web lỗi code browser)
}

// New – engine cho một lần chạy kịch bản sc. in = nil khi không thao tác chuột / phím (Thử node, Thử tìm).
func New(in device.Input, sc model.Script, o Options) *Engine {
	e := newEngine(in, sc, o.OnProgress)
	e.screen, e.displays, e.onNode, e.debug, e.browser = o.Screen, o.Displays, o.OnNode, o.Debug, o.Browser
	e.step.Store(o.StepMode)
	if o.KeepShots {
		e.shots = map[string]image.Image{}
	}
	return e
}

// Failed – node làm lần chạy dừng vì lỗi ("" = không) và, nếu là chờ ảnh hết giờ, chỗ giống nhất.
func (e *Engine) Failed() (string, *FindResult) { return e.failed, e.miss }

// Shots – ảnh màn hình lúc node thất bại theo tên file (Options.KeepShots).
func (e *Engine) Shots() map[string]image.Image { return e.shots }

// FindResult – kết quả thử tìm ảnh mẫu / chữ trên màn hình hiện tại.
type FindResult struct {
	Found bool       `json:"found"` // độ khớp ≥ ngưỡng của bước
	Score int        `json:"score"` // % khớp tốt nhất
	X     int        `json:"x"`     // khung chỗ giống nhất (toạ độ robotgo)
	Y     int        `json:"y"`
	W     int        `json:"w"`
	H     int        `json:"h"`
	Area  model.Rect `json:"area"`           // vùng đã tìm
	Ms    int        `json:"ms"`             // thời gian tìm
	Read  []string   `json:"read,omitempty"` // tìm theo chữ mà không thấy: các dòng đọc được trong vùng
}

// FindOnce tìm ảnh mẫu / chữ của node s trên màn hình lúc này (một lần, không chờ) – nút "Thử tìm".
// {{biến}} trong s phải thay sẵn (script.WithInputDefaults).
func (e *Engine) FindOnce(s model.Step) (FindResult, error) {
	area := e.searchArea(s)
	if err := device.CaptureArea(area, e.displays); err != nil {
		return FindResult{}, err
	}
	start := time.Now()
	if s.FindBy == model.FindByText {
		if strings.TrimSpace(s.FindText) == "" {
			return FindResult{}, errors.New("chưa nhập chữ cần tìm")
		}
		hit, read, ok, err := locateText(e.scr, e.ocr, area, s)
		if err != nil {
			return FindResult{}, err
		}
		res := FindResult{Found: ok, Area: area, Read: read, Ms: int(time.Since(start).Milliseconds()), X: hit.X, Y: hit.Y, W: hit.W, H: hit.H}
		if ok {
			res.Score = 100
		}
		return res, nil
	}
	if s.Template == "" {
		return FindResult{}, errors.New("bước chưa có ảnh mẫu")
	}
	tpl, err := vision.DecodePNG(s.Template)
	if err != nil {
		return FindResult{}, err
	}
	shot, err := e.scr.Capture(area)
	if err != nil {
		return FindResult{}, err
	}
	m, k := vision.MatchScaled(shot, area, vision.ToGray(tpl), s.TemplateScale)
	res := FindResult{
		Score: int(math.Round(m.Score * 100)), Area: area, Ms: int(time.Since(start).Milliseconds()),
		X: area.X + int(float64(m.X)/k), Y: area.Y + int(float64(m.Y)/k), W: int(float64(m.W) / k), H: int(float64(m.H) / k),
	}
	res.Found = res.Score >= s.Threshold
	return res, nil
}
