package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Điều kiện (bước 20 F3) của node Nếu và case của Switch theo màn hình: xét một lần trên màn hình lúc này, trả đúng / sai
// kèm chi tiết từng điều kiện lá (chỗ thấy, % khớp, màu đọc được, 2 vế) – engine ghi vào log, nút Thử node hiện ra.

// Cond.Kind
const (
	CondImage = "image"
	CondText  = "text"
	CondPixel = "pixel"
	CondValue = "value" // so 2 vế Left Op Right (đã thay {{biến}})
	CondAll   = "all"
	CondAny   = "any"
)

var condOps = []string{"eq", "ne", "contains", "gt", "lt", "empty"}

var condOpText = map[string]string{"eq": "=", "ne": "≠", "contains": "chứa", "gt": ">", "lt": "<"}

func isLeafScreen(kind string) bool {
	return kind == CondImage || kind == CondText || kind == CondPixel
}

// CondPart – kết quả một điều kiện lá.
type CondPart struct {
	Path   []int  `json:"path"`           // vị trí trong conds lồng nhau ([] = điều kiện gốc); switchScreen: [case, …]
	Case   string `json:"case,omitempty"` // switchScreen: key của case
	Kind   string `json:"kind"`
	Ok     bool   `json:"ok"`            // kết quả sau khi đã đảo (not)
	Known  bool   `json:"known"`         // false: chỉ biết lúc chạy (vế có biến node {{n2…}} – nút Thử node)
	Detail string `json:"detail"`        // vd "thấy “Gửi” tại X 120 Y 40", "khớp 97% (cần 90%)"
	Hit    *Rect  `json:"hit,omitempty"` // chỗ thấy (toạ độ robotgo)
	Score  int    `json:"score,omitempty"`
	Text   string `json:"text,omitempty"`
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
	scr   Screen
	ocr   TextReader
	shots map[Rect]image.Image
	texts map[image.Image][]TextLine
}

func newSnapshot(scr Screen, ocr TextReader) *snapshot {
	return &snapshot{scr: scr, ocr: ocr, shots: map[Rect]image.Image{}, texts: map[image.Image][]TextLine{}}
}

func (s *snapshot) Capture(r Rect) (image.Image, error) {
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

func (s *snapshot) ReadText(img image.Image) ([]TextLine, error) {
	if l, ok := s.texts[img]; ok {
		return l, nil
	}
	l, err := s.ocr.ReadText(img)
	if err == nil {
		s.texts[img] = l
	}
	return l, err
}

// probe – node Chờ ảnh / chữ tương ứng điều kiện (dùng lại vùng tìm, ảnh mẫu, tìm chữ). key: id riêng để nhớ ảnh mẫu.
func (c Cond) probe(key string) Step {
	s := Step{
		ID: key, Type: StepWaitImage, Template: c.Template, TemplateScale: c.TemplateScale, Threshold: c.Threshold,
		SearchArea: c.SearchArea, FindText: c.FindText, TextExact: c.TextExact, TextNoAccent: c.TextNoAccent, Occurrence: c.Occurrence,
	}
	if c.Kind == CondText {
		s.FindBy = FindByText
	}
	return s
}

// evalCond xét c một lần trên snap; ghi kết quả từng điều kiện lá vào parts. known = false khi kết quả phụ thuộc biến
// node chưa biết (chỉ xảy ra ở Thử node – lúc chạy biến đã thay). Lỗi (thiếu quyền, ảnh hỏng…) là lỗi của node.
func (e *engine) evalCond(c Cond, key string, path []int, snap *snapshot, parts *[]CondPart) (ok, known bool, err error) {
	if c.Kind == CondAll || c.Kind == CondAny {
		all := c.Kind == CondAll
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
func (e *engine) checkLeaf(c Cond, key string, snap *snapshot) (CondPart, error) {
	part := CondPart{Kind: c.Kind, Known: true}
	var ok bool
	switch c.Kind {
	case CondValue:
		if nodeVarRe.MatchString(c.Left) || nodeVarRe.MatchString(c.Right) {
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
	case CondText:
		probe := c.probe(key)
		area := e.searchArea(probe)
		if err := captureArea(area, e.displays); err != nil {
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
			part.Detail = fmt.Sprintf("không thấy “%s”%s", c.FindText, readHint(read))
		}
	case CondImage:
		probe := c.probe(key)
		tpl, err := e.template(probe)
		if err != nil {
			return part, err
		}
		area := e.searchArea(probe)
		if err := captureArea(area, e.displays); err != nil {
			return part, err
		}
		img, err := snap.Capture(area)
		if err != nil {
			return part, err
		}
		m, k := matchScaled(img, area, tpl, probe.TemplateScale)
		score := int(math.Round(m.Score * 100))
		box := Rect{X: area.X + int(float64(m.X)/k), Y: area.Y + int(float64(m.Y)/k), W: int(float64(m.W) / k), H: int(float64(m.H) / k)}
		ok = score >= max(c.Threshold, 1)
		part.Hit, part.Score = &box, score
		part.Detail = fmt.Sprintf("%s tại X %d Y %d · khớp %d%% (cần %d%%)", map[bool]string{true: "thấy ảnh", false: "giống nhất"}[ok],
			box.X+box.W/2, box.Y+box.H/2, score, c.Threshold)
	case CondPixel:
		want, err := parseColor(c.Color)
		if err != nil {
			return part, err
		}
		if len(c.Points) == 0 {
			return part, errors.New("điều kiện màu chưa chọn điểm")
		}
		p := resolvePoint(c.Points[0], bounds(e.sc, e.screen))
		got, err := snap.Pixel(p.X, p.Y)
		if err != nil {
			return part, err
		}
		d := colorDiff(got, want)
		ok = d <= c.Tolerance
		part.Hit = &Rect{X: p.X, Y: p.Y}
		part.Detail = fmt.Sprintf("X %d Y %d màu %s (cần %s ± %d, lệch %d)", p.X, p.Y, hexColor(got), strings.ToUpper(c.Color), c.Tolerance, d)
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
		if p.Ok && p.Hit != nil && p.Kind != CondPixel {
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
		a, err1 := parseNumber(l)
		b, err2 := parseNumber(r)
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

func parseNumber(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(s), ",", "."), 64)
}

// matchCase – Switch theo giá trị: eq / contains như compareValues, regex = khớp mẫu (Go regexp, phân biệt hoa thường
// trừ khi mẫu có (?i)).
func matchCase(value string, c Case) (bool, error) {
	if c.Op == "regex" {
		re, err := regexp.Compile(c.Match)
		if err != nil {
			return false, fmt.Errorf("mẫu %q sai: %v", c.Match, err)
		}
		return re.MatchString(value), nil
	}
	return compareValues(value, c.Op, c.Match)
}

// checkCond kiểm cấu trúc điều kiện (mở file): loại, phép so, lồng nhóm.
func checkCond(c *Cond, depth int) error {
	switch c.Kind {
	case CondImage, CondText, CondPixel:
	case CondValue:
		if !slices.Contains(condOps, c.Op) {
			return fmt.Errorf("phép so không hợp lệ: %q", c.Op)
		}
	case CondAll, CondAny:
		if depth > 8 {
			return errors.New("điều kiện lồng quá sâu")
		}
		for i := range c.Conds {
			if err := checkCond(&c.Conds[i], depth+1); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("loại điều kiện không hợp lệ: %q", c.Kind)
	}
	for _, p := range c.Points {
		if !validAnchor(p.Ax) || !validAnchor(p.Ay) {
			return fmt.Errorf("neo điểm không hợp lệ: %q / %q", p.Ax, p.Ay)
		}
	}
	return nil
}

// validateCond – điều kiện đủ dữ liệu để chạy (trước khi chạy): chữ, ảnh mẫu, điểm + màu, nhóm không rỗng.
// inside: điểm nằm trong vùng.
func validateCond(c Cond, inside func(Point) bool) error {
	switch c.Kind {
	case CondText:
		if strings.TrimSpace(c.FindText) == "" {
			return errors.New("điều kiện chưa nhập chữ cần thấy")
		}
	case CondImage:
		if c.Template == "" {
			return errors.New("điều kiện chưa có ảnh mẫu")
		}
		if _, err := decodePNG(c.Template); err != nil {
			return fmt.Errorf("ảnh mẫu của điều kiện hỏng: %v", err)
		}
		if c.Threshold < 1 || c.Threshold > 100 {
			return errors.New("độ khớp của điều kiện phải trong 1–100%")
		}
	case CondPixel:
		if len(c.Points) == 0 {
			return errors.New("điều kiện màu chưa chọn điểm")
		}
		if !inside(c.Points[0]) {
			return errors.New("điểm của điều kiện màu nằm ngoài vùng")
		}
		if _, err := parseColor(c.Color); err != nil {
			return err
		}
	case CondValue:
		if !slices.Contains(condOps, c.Op) {
			return fmt.Errorf("phép so không hợp lệ %q", c.Op)
		}
	case CondAll, CondAny:
		if len(c.Conds) == 0 {
			return errors.New("nhóm điều kiện chưa có điều kiện nào")
		}
		for i, sub := range c.Conds {
			if err := validateCond(sub, inside); err != nil {
				return fmt.Errorf("điều kiện %d: %w", i+1, err)
			}
		}
	}
	return nil
}

// cloneCond chép sâu (thay biến trên bản sao không đụng kịch bản gốc).
func cloneCond(c *Cond) *Cond {
	if c == nil {
		return nil
	}
	x := *c
	x.Points = slices.Clone(c.Points)
	if c.Conds != nil {
		x.Conds = make([]Cond, len(c.Conds))
		for i := range c.Conds {
			x.Conds[i] = *cloneCond(&c.Conds[i])
		}
	}
	return &x
}
