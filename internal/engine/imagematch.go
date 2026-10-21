package engine

import (
	"fmt"
	"image"
	"math"

	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// So ảnh mẫu của node (chờ / click theo ảnh, điều kiện thấy ảnh, Thử tìm): giải mã một lần, cách so + kiểm tra màu (bước 40).

// template trả ảnh mẫu của bước, giải mã một lần mỗi lần chạy.
func (e *Engine) template(s model.Step) (vision.Template, error) {
	if t, ok := e.tpls[s.ID]; ok {
		return t, nil
	}
	img, err := vision.DecodePNG(s.Template)
	if err != nil {
		return vision.Template{}, err
	}
	t := vision.NewTemplate(img)
	if e.tpls == nil {
		e.tpls = map[string]vision.Template{}
	}
	e.tpls[s.ID] = t
	return t, nil
}

// matchImage tìm ảnh mẫu tpl của s trong ảnh chụp img của vùng area (cách so, kiểm tra màu theo s); khung theo toạ độ
// robotgo. Found = khớp ≥ ngưỡng và đúng màu.
func matchImage(img image.Image, area model.Rect, tpl vision.Template, s model.Step, dpi float64) FindResult {
	opt := vision.MatchOptions{Mode: s.MatchMode, ColorTol: s.ColorTol, MinScore: float64(s.Threshold) / 100}
	m, k := vision.MatchScaled(img, area, tpl, opt, s.TemplateScale, dpi)
	r := FindResult{Score: int(math.Round(m.Score * 100)), Area: area, ColorDiff: m.ColorDiff,
		X: area.X + int(float64(m.X)/k), Y: area.Y + int(float64(m.Y)/k), W: int(float64(m.W) / k), H: int(float64(m.H) / k)}
	r.Found = r.Score >= max(s.Threshold, 1) && m.ColorOK(opt)
	return r
}

// imageMissMsg – lời báo không thấy ảnh: giống hình nhưng lệch màu thì nói rõ để người dùng / AI biết chỉnh gì.
func imageMissMsg(m FindResult, s model.Step) string {
	if s.ColorTol > 0 && m.Score >= s.Threshold {
		return fmt.Sprintf("thấy hình giống (khớp %d%%) nhưng màu lệch %d%% (cho phép %d%%)", m.Score, m.ColorDiff, s.ColorTol)
	}
	return fmt.Sprintf("khớp nhất %d%%, cần %d%%", m.Score, s.Threshold)
}
