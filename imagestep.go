package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"math"
	"strings"
	"time"
)

// Engine cho chờ ảnh / click theo ảnh / chờ pixel (bước 25).

// WaitProgress – đang chờ (event engine:progress, trường wait): thanh mini vẽ vòng tiến độ theo thời gian chờ.
type WaitProgress struct {
	Kind      string `json:"kind"` // image | pixel
	ElapsedMs int64  `json:"elapsedMs"`
	TimeoutMs int64  `json:"timeoutMs"`
	Best      int    `json:"best"` // % khớp tốt nhất tới giờ (pixel: 100 = đúng màu)
}

// searchArea: vùng tìm ảnh theo toạ độ robotgo – vùng của kịch bản, hoặc cả màn chứa vùng đó (không có vùng: màn chính).
func (e *engine) searchArea(s Step) Rect {
	b := bounds(e.sc, e.screen)
	if s.SearchArea != "screen" {
		return b
	}
	cx, cy := b.X+b.W/2, b.Y+b.H/2
	for _, d := range e.displays {
		if cx >= d.X && cx < d.X+d.W && cy >= d.Y && cy < d.Y+d.H {
			return d
		}
	}
	if e.screen.W > 0 {
		return e.screen
	}
	return b
}

// template trả ảnh mẫu (xám) của bước, giải mã một lần mỗi lần chạy.
func (e *engine) template(s Step) (gray, error) {
	if g, ok := e.tpls[s.ID]; ok {
		return g, nil
	}
	img, err := decodePNG(s.Template)
	if err != nil {
		return gray{}, err
	}
	g := toGray(img)
	if e.tpls == nil {
		e.tpls = map[string]gray{}
	}
	e.tpls[s.ID] = g
	return g, nil
}

// waitLoop gọi check mỗi IntervalMs tới khi check trả true, hết TimeoutMs (không tính lúc tạm dừng) hoặc bị huỷ.
// check trả thêm độ khớp 0..100 để báo tiến độ.
func (e *engine) waitLoop(ctx context.Context, s Step, kind string, check func() (bool, int, error)) (best int, err error) {
	start, paused0, first := time.Now(), e.gate.pausedTotal(), true
	defer e.update(true, func(p *Progress) { p.Wait = nil })
	for {
		if err := e.gate.wait(ctx); err != nil {
			return best, err
		}
		ok, score, err := check()
		if err != nil {
			return best, err
		}
		best = max(best, score)
		if ok {
			return best, nil
		}
		elapsed := time.Since(start) - (e.gate.pausedTotal() - paused0)
		e.update(first, func(p *Progress) { // lần đầu báo ngay để thanh mini hiện "đang tìm"
			p.Wait = &WaitProgress{Kind: kind, ElapsedMs: elapsed.Milliseconds(), TimeoutMs: int64(s.TimeoutMs), Best: best}
		})
		first = false
		left := time.Duration(s.TimeoutMs)*time.Millisecond - elapsed
		if left <= 0 {
			return best, errTimeout
		}
		if err := sleep(ctx, int(min(time.Duration(s.IntervalMs)*time.Millisecond, left).Milliseconds())); err != nil {
			return best, err
		}
	}
}

var errTimeout = errors.New("hết giờ")

// hit – chỗ thấy ảnh / chữ (toạ độ robotgo): output của node chờ ảnh / chữ, node click theo ảnh lấy lại qua ImageFrom.
type hit struct {
	Box   Rect
	Score int    // % khớp (tìm theo chữ: 100)
	Text  string // tìm theo chữ: dòng chứa chữ
}

func (h hit) center() Point { return Point{X: h.Box.X + h.Box.W/2, Y: h.Box.Y + h.Box.H/2} }

func (h hit) output() map[string]any {
	c := h.center()
	out := map[string]any{"x": c.X, "y": c.Y, "w": h.Box.W, "h": h.Box.H, "score": h.Score}
	if h.Text != "" {
		out["text"] = h.Text
	}
	return out
}

// findImage chờ tới khi thấy ảnh mẫu (hoặc chữ, FindBy text); trả chỗ thấy.
func (e *engine) findImage(ctx context.Context, s Step) (hit, error) {
	area := e.searchArea(s)
	if err := captureArea(area, e.displays); err != nil {
		return hit{}, err
	}
	if s.FindBy == FindByText {
		return e.findTextStep(ctx, s, area)
	}
	tpl, err := e.template(s)
	if err != nil {
		return hit{}, err
	}
	var at hit
	var miss FindResult // chỗ giống nhất tới giờ – hết giờ thì báo lại
	best, err := e.waitLoop(ctx, s, "image", func() (bool, int, error) {
		img, err := e.scr.Capture(area)
		if err != nil {
			return false, 0, err
		}
		m, sc := matchScaled(img, area, tpl, s.TemplateScale)
		score := int(math.Round(m.Score * 100))
		if score >= miss.Score {
			miss = FindResult{Score: score, Area: area, X: area.X + int(float64(m.X)/sc), Y: area.Y + int(float64(m.Y)/sc),
				W: int(float64(m.W) / sc), H: int(float64(m.H) / sc)}
		}
		if score < s.Threshold {
			return false, score, nil
		}
		at = hit{Score: score, Box: Rect{X: area.X + int(float64(m.X)/sc), Y: area.Y + int(float64(m.Y)/sc), W: int(float64(m.W) / sc), H: int(float64(m.H) / sc)}}
		return true, score, nil
	})
	if errors.Is(err, errTimeout) {
		e.miss = &miss
		return hit{}, &notFoundErr{fmt.Sprintf("không thấy ảnh mẫu sau %s (khớp nhất %d%%, cần %d%%)", fmtSeconds(s.TimeoutMs), best, s.Threshold)}
	}
	if err != nil {
		return hit{}, err
	}
	return at, nil
}

// findTextStep chờ tới khi đọc được chữ s.FindText trong vùng; trả khung chữ (toạ độ robotgo).
func (e *engine) findTextStep(ctx context.Context, s Step, area Rect) (hit, error) {
	var at hit
	var read []string
	_, err := e.waitLoop(ctx, s, "text", func() (bool, int, error) {
		box, lines, ok, err := locateText(e.scr, e.ocr, area, s)
		if err != nil {
			return false, 0, err
		}
		if !ok {
			read = lines
			return false, 0, nil
		}
		at = hit{Box: box, Score: 100, Text: s.FindText}
		return true, 100, nil
	})
	if errors.Is(err, errTimeout) {
		e.miss = &FindResult{Area: area, Read: read}
		return hit{}, &notFoundErr{fmt.Sprintf("không thấy chữ “%s” sau %s%s", s.FindText, fmtSeconds(s.TimeoutMs), readHint(read))}
	}
	if err != nil {
		return hit{}, err
	}
	return at, nil
}

// locateText chụp vùng, đọc chữ, tìm theo s; hit theo toạ độ robotgo. Không thấy → trả các dòng đọc được.
func locateText(scr Screen, ocr TextReader, area Rect, s Step) (hit Rect, read []string, ok bool, err error) {
	img, err := scr.Capture(area)
	if err != nil {
		return Rect{}, nil, false, err
	}
	lines, err := ocr.ReadText(img)
	if err != nil {
		return Rect{}, nil, false, err
	}
	hits := findText(lines, TextQuery{Text: s.FindText, Exact: s.TextExact, NoAccent: s.TextNoAccent})
	n := max(s.Occurrence, 1)
	if len(hits) < n {
		return Rect{}, textLines(lines, 40), false, nil
	}
	k := imageScale(img, area)
	b := hits[n-1].Box
	return Rect{X: area.X + int(float64(b.X)/k), Y: area.Y + int(float64(b.Y)/k), W: int(float64(b.W) / k), H: int(float64(b.H) / k)}, nil, true, nil
}

// readHint – vài dòng đọc được, gợi ý vì sao không khớp (sai dấu, chữ khác…).
func readHint(read []string) string {
	if len(read) == 0 {
		return " (không đọc được chữ nào trong vùng)"
	}
	show := read[:min(len(read), 5)]
	return fmt.Sprintf(" – đọc được: “%s”%s", strings.Join(show, "”, “"), map[bool]string{true: "…", false: ""}[len(read) > 5])
}

// matchScaled đưa mẫu và ảnh chụp về cùng tỷ lệ px (mẫu chụp ở màn 1x, chạy trên Retina…) rồi tìm. Luôn thu nhỏ
// bên nét hơn (trung bình theo ô ≈ cách màn 1x vẽ) thay vì phóng to bên kia – phóng to làm mẫu nhoè, độ khớp tụt.
// Trả kết quả theo px ảnh chụp và tỷ lệ px / toạ độ của ảnh chụp.
func matchScaled(img image.Image, area Rect, tpl gray, tplScale float64) (Match, float64) {
	sc := imageScale(img, area)
	g := toGray(img)
	if tplScale <= 0 || math.Abs(sc-tplScale) <= 0.01 {
		return findTemplate(g, tpl), sc
	}
	if k := tplScale / sc; k > 1 { // mẫu nét hơn ảnh chụp → thu nhỏ mẫu
		return findTemplate(g, tpl.resize(int(math.Round(float64(tpl.w)/k)), int(math.Round(float64(tpl.h)/k)))), sc
	}
	k := sc / tplScale // ảnh chụp nét hơn → thu nhỏ ảnh chụp, đổi kết quả về px ảnh chụp
	m := findTemplate(g.resize(int(math.Round(float64(g.w)/k)), int(math.Round(float64(g.h)/k))), tpl)
	return Match{X: int(math.Round(float64(m.X) * k)), Y: int(math.Round(float64(m.Y) * k)),
		W: int(math.Round(float64(m.W) * k)), H: int(math.Round(float64(m.H) * k)), Score: m.Score}, sc
}

// clickImage click vào chỗ thấy (+ Dx / Dy): chỗ node ImageFrom thấy lần gần nhất trong vòng, hoặc tự tìm ảnh / chữ.
// Output: chỗ đã click, kèm chỗ thấy khi tự tìm.
func (e *engine) clickImage(ctx context.Context, s Step) (map[string]any, error) {
	var out map[string]any
	var p Point
	if s.ImageSource == ImageFromSearch {
		h, err := e.findImage(ctx, s)
		if err != nil {
			return nil, err
		}
		out, p = h.output(), h.center()
	} else {
		r := e.results[s.ImageFrom]
		ref := e.byID[s.ImageFrom].Ref
		if r == nil || r.Status != StatusSuccess {
			return nil, fmt.Errorf("node %s chưa tìm thấy gì trong vòng này (chưa chạy hoặc lần gần nhất thất bại)", ref)
		}
		x, okX := r.Output["x"].(int)
		y, okY := r.Output["y"].(int)
		if !okX || !okY {
			return nil, fmt.Errorf("node %s không có chỗ thấy (output x / y)", ref)
		}
		out, p = map[string]any{}, Point{X: x, Y: y}
	}
	x, y := e.jitter(p.X+s.Dx, s.JitterPx), e.jitter(p.Y+s.Dy, s.JitterPx)
	if err := e.move(x, y); err != nil {
		return nil, err
	}
	if err := e.click(s.Button, false); err != nil {
		return nil, err
	}
	out["clickX"], out["clickY"] = x, y
	return out, nil
}

// waitPixel chờ tới khi điểm Points[0] có màu Color (mỗi kênh lệch ≤ Tolerance); trả màu đọc được.
func (e *engine) waitPixel(ctx context.Context, s Step) (string, error) {
	want, err := parseColor(s.Color)
	if err != nil {
		return "", err
	}
	p := resolvePoint(s.Points[0], bounds(e.sc, e.screen))
	var last string
	_, err = e.waitLoop(ctx, s, "pixel", func() (bool, int, error) {
		c, err := e.scr.Pixel(p.X, p.Y)
		if err != nil {
			return false, 0, err
		}
		last = hexColor(c)
		d := colorDiff(c, want)
		return d <= s.Tolerance, 100 - d*100/255, nil
	})
	if errors.Is(err, errTimeout) {
		return "", &notFoundErr{fmt.Sprintf("điểm X %d Y %d vẫn màu %s sau %s (cần %s ± %d)", p.X, p.Y, last, fmtSeconds(s.TimeoutMs), s.Color, s.Tolerance)}
	}
	return last, err
}

func fmtSeconds(ms int) string {
	if ms%1000 == 0 {
		return fmt.Sprintf("%d s", ms/1000)
	}
	return fmt.Sprintf("%.1f s", float64(ms)/1000)
}
