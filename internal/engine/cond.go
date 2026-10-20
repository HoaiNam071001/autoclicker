package engine

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"regexp"
	"slices"
	"strings"

	"autoclicker/internal/device"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
	"autoclicker/internal/vision"
)

// Điều kiện (bước 20 F3) của node Nếu và case của Switch theo màn hình: xét một lần trên màn hình lúc này, trả đúng / sai
// kèm chi tiết từng điều kiện lá (chỗ thấy, % khớp, màu đọc được, 2 vế) – engine ghi vào log, nút Thử node hiện ra.

var condOpText = map[string]string{"eq": "=", "ne": "≠", "contains": "chứa", "gt": ">", "lt": "<"}

// CondPart – kết quả một điều kiện lá.
type CondPart struct {
	Path   []int       `json:"path"`           // vị trí trong conds lồng nhau ([] = điều kiện gốc); switchScreen: [case, …]
	Case   string      `json:"case,omitempty"` // switchScreen: key của case
	Kind   string      `json:"kind"`
	Ok     bool        `json:"ok"`            // kết quả sau khi đã đảo (not)
	Known  bool        `json:"known"`         // false: chỉ biết lúc chạy (vế có biến node {{n2…}} – nút Thử node)
	Detail string      `json:"detail"`        // vd "thấy “Gửi” tại X 120 Y 40", "khớp 97% (cần 90%)"
	Hit    *model.Rect `json:"hit,omitempty"` // chỗ thấy (toạ độ robotgo)
	Score  int         `json:"score,omitempty"`
	Text   string      `json:"text,omitempty"`
}

// output – chỗ thấy của điều kiện lá dưới dạng output node (x, y tâm…) để Click theo ảnh dùng lại.
func (p CondPart) output() map[string]any {
	if p.Hit == nil {
		return map[string]any{}
	}
	return hit{Box: *p.Hit, Score: p.Score, Text: p.Text}.output()
}

// snapshot – một lần xét: chụp mỗi vùng một lần, đọc chữ mỗi ảnh một lần, dù nhiều điều kiện cùng vùng.
type snapshot struct {
	scr   device.Screen
	ocr   device.TextReader
	shots map[model.Rect]image.Image
	texts map[image.Image][]vision.TextLine
	elems map[string]elementRead // cây phần tử theo app (điều kiện "thấy phần tử")
}

func newSnapshot(scr device.Screen, ocr device.TextReader) *snapshot {
	return &snapshot{scr: scr, ocr: ocr, shots: map[model.Rect]image.Image{}, texts: map[image.Image][]vision.TextLine{}, elems: map[string]elementRead{}}
}

func (s *snapshot) Capture(r model.Rect) (image.Image, error) {
	if img, ok := s.shots[r]; ok {
		return img, nil
	}
	img, err := s.scr.Capture(r)
	if err == nil {
		s.shots[r] = img
	}
	return img, err
}

func (s *snapshot) Pixel(x, y int) (color.RGBA, error) { return s.scr.Pixel(x, y) }

func (s *snapshot) DPIScale(r model.Rect) float64 { return s.scr.DPIScale(r) }

func (s *snapshot) ReadText(img image.Image) ([]vision.TextLine, error) {
	if l, ok := s.texts[img]; ok {
		return l, nil
	}
	l, err := s.ocr.ReadText(img)
	if err == nil {
		s.texts[img] = l
	}
	return l, err
}

// evalCond xét c một lần trên snap; ghi kết quả từng điều kiện lá vào parts. known = false khi kết quả phụ thuộc biến
// node chưa biết (chỉ xảy ra ở Thử node – lúc chạy biến đã thay). Lỗi (thiếu quyền, ảnh hỏng…) là lỗi của node.
func (e *Engine) evalCond(c model.Cond, key string, path []int, snap *snapshot, parts *[]CondPart) (ok, known bool, err error) {
	if c.Kind == model.CondAll || c.Kind == model.CondAny {
		all := c.Kind == model.CondAll
		ok, known = all, true
		decided := false
		for i, sub := range c.Conds {
			v, k, err := e.evalCond(sub, fmt.Sprintf("%s.%d", key, i), append(slices.Clone(path), i), snap, parts)
			if err != nil {
				return false, false, err
			}
			if decided {
				continue // đã biết kết quả – vẫn xét hết để báo từng phần
			}
			if !k {
				known = false
				continue
			}
			if v != all { // all: gặp sai → sai; any: gặp đúng → đúng
				ok, known, decided = v, true, true
			}
		}
		return ok != c.Not, known, nil
	}
	part, err := e.checkLeaf(c, key, snap)
	if err != nil {
		return false, false, err
	}
	part.Path = slices.Clone(path)
	*parts = append(*parts, part)
	return part.Ok, part.Known, nil
}

// checkLeaf xét một điều kiện lá (ảnh / chữ / màu / so giá trị).
func (e *Engine) checkLeaf(c model.Cond, key string, snap *snapshot) (CondPart, error) {
	part := CondPart{Kind: c.Kind, Known: true}
	var ok bool
	switch c.Kind {
	case model.CondValue:
		if script.NodeVarRe.MatchString(c.Left) || script.NodeVarRe.MatchString(c.Right) {
			part.Known, part.Detail = false, fmt.Sprintf("“%s” – chỉ biết lúc chạy (biến node)", c.Left)
			return part, nil
		}
		v, err := compareValues(c.Left, c.Op, c.Right)
		if err != nil {
			return part, err
		}
		ok = v
		if c.Op == "empty" {
			part.Detail = fmt.Sprintf("“%s” trống", c.Left)
		} else {
			part.Detail = fmt.Sprintf("“%s” %s “%s”", c.Left, condOpText[c.Op], c.Right)
		}
	case model.CondText:
		probe := c.Probe(key)
		area := e.searchArea(probe)
		if err := device.CaptureArea(area, e.displays); err != nil {
			return part, err
		}
		box, read, found, err := locateText(snap, snap, area, probe)
		if err != nil {
			return part, err
		}
		ok = found
		if found {
			part.Hit, part.Score, part.Text = &box, 100, c.FindText
			part.Detail = fmt.Sprintf("thấy “%s” tại X %d Y %d", c.FindText, box.X+box.W/2, box.Y+box.H/2)
		} else {
			part.Detail = fmt.Sprintf("không thấy “%s”%s", c.FindText, readHint(vision.TextLines(read, 40)))
		}
	case model.CondElement:
		v, err := e.elementCond(c, snap, &part)
		if err != nil {
			return part, err
		}
		ok = v
	case model.CondImage:
		probe := c.Probe(key)
		tpl, err := e.template(probe)
		if err != nil {
			return part, err
		}
		area := e.searchArea(probe)
		if err := device.CaptureArea(area, e.displays); err != nil {
			return part, err
		}
		img, err := snap.Capture(area)
		if err != nil {
			return part, err
		}
		m, k := vision.MatchScaled(img, area, tpl, probe.TemplateScale, snap.DPIScale(area))
		score := int(math.Round(m.Score * 100))
		box := model.Rect{X: area.X + int(float64(m.X)/k), Y: area.Y + int(float64(m.Y)/k), W: int(float64(m.W) / k), H: int(float64(m.H) / k)}
		ok = score >= max(c.Threshold, 1)
		part.Hit, part.Score = &box, score
		part.Detail = fmt.Sprintf("%s tại X %d Y %d · khớp %d%% (cần %d%%)", map[bool]string{true: "thấy ảnh", false: "giống nhất"}[ok],
			box.X+box.W/2, box.Y+box.H/2, score, c.Threshold)
	case model.CondPixel:
		want, err := vision.ParseColor(c.Color)
		if err != nil {
			return part, err
		}
		if len(c.Points) == 0 {
			return part, errors.New("điều kiện màu chưa chọn điểm")
		}
		p := script.ResolvePoint(c.Points[0], script.Bounds(e.sc, e.screen))
		got, err := snap.Pixel(p.X, p.Y)
		if err != nil {
			return part, err
		}
		d := vision.ColorDiff(got, want)
		ok = d <= c.Tolerance
		part.Hit = &model.Rect{X: p.X, Y: p.Y}
		part.Detail = fmt.Sprintf("X %d Y %d màu %s (cần %s ± %d, lệch %d)", p.X, p.Y, vision.HexColor(got), strings.ToUpper(c.Color), c.Tolerance, d)
	default:
		return part, fmt.Errorf("loại điều kiện không hợp lệ %q", c.Kind)
	}
	part.Ok = ok != c.Not
	if c.Not {
		part.Detail += " → đảo"
	}
	return part, nil
}

// firstHit – output chỗ thấy của điều kiện lá đầu tiên đúng có chỗ thấy (Click theo ảnh lấy từ node Nếu / Switch).
func firstHit(parts []CondPart) map[string]any {
	for _, p := range parts {
		if p.Ok && p.Hit != nil && p.Kind != model.CondPixel {
			return p.output()
		}
	}
	return map[string]any{}
}

// compareValues so 2 vế (đã thay {{biến}}). eq / ne / contains không phân biệt hoa thường, bỏ khoảng trắng 2 đầu;
// gt / lt so số (nhận dấu phẩy thập phân).
func compareValues(left, op, right string) (bool, error) {
	l, r := strings.TrimSpace(left), strings.TrimSpace(right)
	switch op {
	case "eq":
		return strings.EqualFold(l, r), nil
	case "ne":
		return !strings.EqualFold(l, r), nil
	case "contains":
		return strings.Contains(strings.ToLower(l), strings.ToLower(r)), nil
	case "empty":
		return l == "", nil
	case "gt", "lt":
		a, err1 := script.ParseNumber(l)
		b, err2 := script.ParseNumber(r)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("so lớn / nhỏ cần 2 số, có “%s” và “%s”", l, r)
		}
		if op == "gt" {
			return a > b, nil
		}
		return a < b, nil
	}
	return false, fmt.Errorf("phép so không hợp lệ %q", op)
}

// matchCase – Switch theo giá trị: eq / contains như compareValues, regex = khớp mẫu (Go regexp, phân biệt hoa thường
// trừ khi mẫu có (?i)).
func matchCase(value string, c model.Case) (bool, error) {
	if c.Op == "regex" {
		re, err := regexp.Compile(c.Match)
		if err != nil {
			return false, fmt.Errorf("mẫu %q sai: %v", c.Match, err)
		}
		return re.MatchString(value), nil
	}
	return compareValues(value, c.Op, c.Match)
}
