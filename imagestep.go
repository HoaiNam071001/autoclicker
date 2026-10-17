package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"math"
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

// searchArea: vùng tìm ảnh theo toạ độ robotgo – vùng của kịch bản, hoặc cả màn hình chính.
func (e *engine) searchArea(s Step) Rect {
	if s.SearchArea == "screen" && e.screen.W > 0 {
		return e.screen
	}
	return bounds(e.sc, e.screen)
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

// findImage chờ tới khi thấy ảnh mẫu; trả tâm chỗ thấy (toạ độ robotgo) và nhớ lại cho bước click theo ảnh sau.
func (e *engine) findImage(ctx context.Context, s Step) (Point, error) {
	area := e.searchArea(s)
	if err := captureArea(area, e.screen); err != nil {
		return Point{}, err
	}
	tpl, err := e.template(s)
	if err != nil {
		return Point{}, err
	}
	var at Point
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
		at = Point{X: area.X + int((float64(m.X)+float64(m.W)/2)/sc), Y: area.Y + int((float64(m.Y)+float64(m.H)/2)/sc)}
		return true, score, nil
	})
	if errors.Is(err, errTimeout) {
		e.miss = &miss
		return Point{}, fmt.Errorf("không thấy ảnh mẫu sau %s (khớp nhất %d%%, cần %d%%)", fmtSeconds(s.TimeoutMs), best, s.Threshold)
	}
	if err != nil {
		return Point{}, err
	}
	e.found, e.hasFound = at, true
	return at, nil
}

// matchScaled co giãn mẫu về cùng tỷ lệ px với ảnh chụp (mẫu chụp trên Retina, chạy trên màn 1x…) rồi tìm.
// Trả kết quả theo px ảnh chụp và tỷ lệ px / toạ độ của ảnh chụp.
func matchScaled(img image.Image, area Rect, tpl gray, tplScale float64) (Match, float64) {
	sc := imageScale(img, area)
	if tplScale > 0 && math.Abs(sc-tplScale) > 0.01 {
		k := sc / tplScale
		tpl = tpl.resize(int(math.Round(float64(tpl.w)*k)), int(math.Round(float64(tpl.h)*k)))
	}
	return findTemplate(toGray(img), tpl), sc
}

// clickImage click vào chỗ thấy ảnh (+ Dx / Dy): kết quả bước chờ ảnh trước, hoặc tự tìm ảnh mẫu của bước.
func (e *engine) clickImage(ctx context.Context, s Step) error {
	p := e.found
	if s.ImageSource == ImageFromSearch {
		var err error
		if p, err = e.findImage(ctx, s); err != nil {
			return err
		}
	} else if !e.hasFound {
		return errors.New("phía trước chưa có bước chờ ảnh nào tìm thấy ảnh")
	}
	if err := e.move(e.jitter(p.X+s.Dx, s.JitterPx), e.jitter(p.Y+s.Dy, s.JitterPx)); err != nil {
		return err
	}
	return e.click(s.Button, false)
}

// waitPixel chờ tới khi điểm Points[0] có màu Color (mỗi kênh lệch ≤ Tolerance).
func (e *engine) waitPixel(ctx context.Context, s Step) error {
	want, err := parseColor(s.Color)
	if err != nil {
		return err
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
		return fmt.Errorf("điểm X %d Y %d vẫn màu %s sau %s (cần %s ± %d)", p.X, p.Y, last, fmtSeconds(s.TimeoutMs), s.Color, s.Tolerance)
	}
	return err
}

func fmtSeconds(ms int) string {
	if ms%1000 == 0 {
		return fmt.Sprintf("%d s", ms/1000)
	}
	return fmt.Sprintf("%.1f s", float64(ms)/1000)
}
