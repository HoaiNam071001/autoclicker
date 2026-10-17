package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"sync"
	"time"
)

// dragStepMs – khoảng thời gian giữa các lần di chuột khi kéo thả.
const dragStepMs = 10

// scrollGapMs – nghỉ giữa các nấc cuộn (app kịp nhận từng nấc; dừng được giữa chừng).
const scrollGapMs = 30

// validateScript kiểm tra kịch bản (đã áp input) trước khi chạy để không dừng giữa chừng vì dữ liệu sai.
// Lỗi là *ScriptError (bước / trường) để UI và AI biết chỗ sửa.
// screen: khung màn hình khi Region = nil; W = 0 thì bỏ qua kiểm tra điểm nằm trong vùng (không biết màn hình).
func validateScript(sc Script, screen Rect) error {
	b := bounds(sc, screen)
	inside := func(p Point) bool {
		q := resolvePoint(p, b)
		return b.W == 0 || q.X >= b.X && q.X <= b.X+b.W && q.Y >= b.Y && q.Y <= b.Y+b.H
	}
	enabled := 0
	hasPrevPoint := false // đã có bước chuột (để lại vị trí chuột) phía trước chưa
	hasPrevImage := false // đã có bước chờ ảnh / click theo ảnh tự tìm phía trước chưa (click theo kết quả bước trước)
	for i, s := range sc.Steps {
		if !s.Enabled {
			continue
		}
		enabled++
		if s.Type == StepDrag && s.DragFromPrev && !hasPrevPoint {
			return stepErr(i, s, "dragFromPrev", "kéo từ điểm trước đó nhưng phía trước không có bước chuột nào")
		}
		need := 0
		switch s.Type {
		case StepClick, StepDouble, StepHold:
			need = 1
		case StepDrag:
			need = 2
		case StepScroll:
			need = min(len(s.Points), 1) // có điểm thì điểm phải trong vùng
			if !slices.Contains(scrollDirections, s.Direction) {
				return stepErr(i, s, "direction", "hướng cuộn không hợp lệ %q", s.Direction)
			}
			if s.Amount < 1 {
				return stepErr(i, s, "amount", "số nấc cuộn phải ≥ 1")
			}
		case StepOpen:
			if strings.TrimSpace(s.Target) == "" || s.Target == "https://" {
				return stepErr(i, s, "target", "chưa nhập URL / app cần mở")
			}
		case StepWaitImage, StepClickImage:
			if s.Type == StepClickImage && s.ImageSource != ImageFromSearch {
				if !hasPrevImage {
					return stepErr(i, s, "imageSource", "click theo kết quả bước trước nhưng phía trước không có bước chờ ảnh nào")
				}
			} else if err := checkImageStep(i, s); err != nil {
				return err
			}
			hasPrevImage = true
			if s.Type == StepClickImage {
				hasPrevPoint = true // click theo ảnh di chuột tới chỗ thấy ảnh
			}
		case StepWaitPixel:
			need = 1
			if _, err := parseColor(s.Color); err != nil {
				return stepErr(i, s, "color", "%v", err)
			}
			if s.Tolerance < 0 || s.Tolerance > 255 {
				return stepErr(i, s, "tolerance", "sai số màu phải trong 0–255")
			}
			if s.TimeoutMs <= 0 || s.IntervalMs <= 0 {
				return stepErr(i, s, "timeoutMs", "thời gian chờ tối đa / kiểm tra mỗi phải > 0")
			}
		case StepWait:
		case StepText:
			if s.Text == "" {
				return stepErr(i, s, "text", "chưa nhập nội dung gõ")
			}
		case StepKey:
			if len(s.Keys) == 0 {
				return stepErr(i, s, "keys", "chưa chọn phím")
			}
		default:
			return stepErr(i, s, "type", "loại bước không hợp lệ %q", s.Type)
		}
		if len(s.Points) < need {
			return stepErr(i, s, "points", "cần %d điểm, mới có %d", need, len(s.Points))
		}
		for j := 0; j < need; j++ {
			skip := s.Type == StepDrag && (j == 0 && s.DragFromPrev || j == 1 && s.DragRelative)
			if !skip && !inside(s.Points[j]) {
				q := resolvePoint(s.Points[j], b)
				return stepErr(i, s, "points", "điểm X %d Y %d nằm ngoài vùng", q.X, q.Y)
			}
		}
		hasPrevPoint = hasPrevPoint || need > 0 && s.Type != StepWaitPixel // chờ pixel không di chuột
	}
	if enabled == 0 {
		return &ScriptError{Msg: "kịch bản không có bước nào được bật"}
	}
	return nil
}

// Progress – tiến độ gửi sang frontend (event "engine:progress").
type Progress struct {
	Loop      int           `json:"loop"`  // vòng hiện tại, từ 1
	Loops     int           `json:"loops"` // 0 = vô hạn
	Step      int           `json:"step"`  // thứ tự bước hiện tại trong các bước đang bật, từ 1
	Total     int           `json:"total"` // số bước đang bật
	Clicks    int           `json:"clicks"`
	ElapsedMs int64         `json:"elapsedMs"` // không tính thời gian tạm dừng
	Paused    bool          `json:"paused"`
	Wait      *WaitProgress `json:"wait,omitempty"` // đang chờ ảnh / pixel
}

// progressEvery – giãn cách tối thiểu giữa 2 lần báo tiến độ trong cùng một bước (tránh dội event).
const progressEvery = 100 * time.Millisecond

// engine chạy một kịch bản. Tạo mới cho mỗi lần chạy.
type engine struct {
	in         Input
	sc         Script
	gate       pauseGate
	screen     Rect                           // khung màn hình khi Region = nil (runner đặt = mainScreen); zero = gốc (0, 0)
	displays   []Rect                         // khung các màn (tìm ảnh "cả màn", kiểm tra vùng chụp); nil = không kiểm tra
	onProgress func(Progress)                 // có thể nil
	intN       func(n int) int                // số ngẫu nhiên [0, n) – test thay bằng hàm cố định
	open       func(target, app string) error // bước Mở – test thay bằng hàm ghi lại
	scr        Screen                         // chụp màn hình / đọc màu (imagestep.go) – test thay bằng màn hình giả
	ocr        TextReader                     // đọc chữ (tìm theo chữ, ocr.go) – test thay bằng bộ đọc giả

	failed   string      // id bước làm lần chạy dừng vì lỗi (runner gửi engine:failed)
	miss     *FindResult // bước chờ ảnh hết giờ: chỗ giống nhất (để xem lại vì sao không thấy)
	found    Point       // chỗ bước chờ ảnh / click theo ảnh gần nhất tìm thấy (ImageFromPrev)
	hasFound bool
	tpls     map[string]gray // ảnh mẫu đã giải mã theo id bước

	last    Point // chỗ engine vừa đưa chuột tới – cho bước kéo "từ điểm trước đó"
	hasLast bool

	mu       sync.Mutex
	prog     Progress
	start    time.Time
	lastEmit time.Time
}

func newEngine(in Input, sc Script, onProgress func(Progress)) *engine {
	return &engine{in: in, sc: sc, onProgress: onProgress, intN: rand.IntN, open: openTarget, scr: robotgoScreen{}, ocr: systemOCR{}}
}

// jitter trả v lệch ngẫu nhiên đều trong [v-j, v+j]; j ≤ 0 → v. Dùng cho toạ độ (có thể âm ở màn phụ).
func (e *engine) jitter(v, j int) int {
	if j <= 0 {
		return v
	}
	return v + e.intN(2*j+1) - j
}

// jitterMs như jitter nhưng không âm – dùng cho thời gian chờ.
func (e *engine) jitterMs(v, j int) int { return max(e.jitter(v, j), 0) }

// run chạy tới khi đủ số vòng, hết MaxMinutes hoặc ctx bị huỷ.
// Hết thời gian là kết thúc bình thường (nil); bị huỷ trả ctx.Err().
// MaxMinutes tính theo đồng hồ thật, kể cả lúc tạm dừng.
func (e *engine) run(ctx context.Context) error {
	sc := e.sc
	if err := validateScript(sc, e.screen); err != nil {
		return err
	}
	runCtx := ctx
	if sc.MaxMinutes > 0 {
		var cancel context.CancelFunc
		runCtx, cancel = context.WithTimeout(ctx, time.Duration(sc.MaxMinutes)*time.Minute)
		defer cancel()
	}

	total := 0
	for _, s := range sc.Steps {
		if s.Enabled {
			total++
		}
	}
	e.mu.Lock()
	e.start = time.Now()
	e.prog = Progress{Loops: sc.Loops, Total: total}
	e.mu.Unlock()

	for loop := 1; sc.Loops == 0 || loop <= sc.Loops; loop++ {
		n := 0
		for i, s := range sc.Steps {
			if !s.Enabled {
				continue
			}
			n++
			err := e.gate.wait(runCtx)
			if err == nil {
				e.update(true, func(p *Progress) { p.Loop, p.Step = loop, n })
				err = e.runStep(runCtx, s)
			}
			if err == nil {
				err = sleep(runCtx, e.jitterMs(s.DelayAfter, s.JitterMs))
			}
			if err != nil {
				if ctx.Err() == nil && errors.Is(err, context.DeadlineExceeded) {
					return nil // hết MaxMinutes
				}
				if ctx.Err() != nil {
					return ctx.Err()
				}
				e.failed = s.ID
				return fmt.Errorf("vòng %d, bước %d: %w", loop, i+1, err)
			}
		}
	}
	return nil
}

// update sửa tiến độ rồi báo sang frontend (force = báo ngay, không giãn cách).
func (e *engine) update(force bool, f func(*Progress)) {
	e.mu.Lock()
	f(&e.prog)
	now := time.Now()
	if !force && now.Sub(e.lastEmit) < progressEvery {
		e.mu.Unlock()
		return
	}
	e.lastEmit = now
	p := e.snapshotLocked()
	e.mu.Unlock()
	if e.onProgress != nil {
		e.onProgress(p)
	}
}

// Snapshot trả tiến độ hiện tại (kèm thời gian đã chạy và trạng thái tạm dừng).
func (e *engine) Snapshot() Progress {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.snapshotLocked()
}

func (e *engine) snapshotLocked() Progress {
	p := e.prog
	p.ElapsedMs = e.gate.activeSince(e.start).Milliseconds()
	p.Paused = e.gate.paused()
	return p
}

// Pause / Resume: trả false nếu trạng thái không đổi. Có hiệu lực trước bước / lần click kế tiếp.
func (e *engine) Pause() bool {
	if !e.gate.pause() {
		return false
	}
	e.update(true, func(*Progress) {})
	return true
}

func (e *engine) Resume() bool {
	if !e.gate.resume() {
		return false
	}
	e.update(true, func(*Progress) {})
	return true
}

func (e *engine) click(button string, double bool) error {
	if err := e.in.Click(button, double); err != nil {
		return err
	}
	e.update(false, func(p *Progress) { p.Clicks++ })
	return nil
}

// move đưa chuột tới (x, y) và nhớ lại vị trí đó.
func (e *engine) move(x, y int) error {
	e.last, e.hasLast = Point{X: x, Y: y}, true
	return e.in.Move(x, y)
}

func (e *engine) runStep(ctx context.Context, s Step) error {
	in := e.in
	// pt: toạ độ màn hình của điểm i (theo neo), mỗi lần gọi lệch ngẫu nhiên ± JitterPx.
	pt := func(i int) (int, int) {
		p := resolvePoint(s.Points[i], bounds(e.sc, e.screen))
		return e.jitter(p.X, s.JitterPx), e.jitter(p.Y, s.JitterPx)
	}

	switch s.Type {
	case StepClick:
		n := max(s.Clicks, 1)
		for i := 0; i < n; i++ {
			if err := e.gate.wait(ctx); err != nil {
				return err
			}
			x, y := pt(0) // mỗi click lệch riêng
			if err := e.move(x, y); err != nil {
				return err
			}
			if err := e.click(s.Button, false); err != nil {
				return err
			}
			if i < n-1 {
				if err := sleep(ctx, e.jitterMs(s.ClickInterval, s.JitterMs)); err != nil {
					return err
				}
			}
		}
	case StepDouble:
		x, y := pt(0)
		if err := e.move(x, y); err != nil {
			return err
		}
		return e.click(s.Button, true)
	case StepHold:
		x, y := pt(0)
		if err := e.move(x, y); err != nil {
			return err
		}
		if err := in.Toggle(s.Button, true); err != nil {
			return err
		}
		defer in.Toggle(s.Button, false) // luôn nhả chuột, kể cả khi bị huỷ giữa chừng
		return sleep(ctx, s.HoldMs)
	case StepDrag:
		return e.drag(ctx, s, pt)
	case StepWait:
		// chỉ chờ DelayAfter
	case StepText:
		// Gõ từng ký tự để dừng / tạm dừng được giữa chừng với chuỗi dài.
		for _, r := range s.Text {
			if err := e.gate.wait(ctx); err != nil {
				return err
			}
			if err := in.TypeStr(string(r)); err != nil {
				return err
			}
		}
	case StepKey:
		return pressKeys(ctx, in, s)
	case StepScroll:
		if len(s.Points) > 0 {
			x, y := pt(0)
			if err := e.move(x, y); err != nil {
				return err
			}
		}
		for i := 0; i < s.Amount; i++ {
			if err := e.gate.wait(ctx); err != nil {
				return err
			}
			if err := in.Scroll(s.Direction); err != nil {
				return err
			}
			if err := sleep(ctx, scrollGapMs); err != nil {
				return err
			}
		}
	case StepOpen:
		return e.open(s.Target, s.App)
	case StepWaitImage:
		_, err := e.findImage(ctx, s)
		return err
	case StepClickImage:
		return e.clickImage(ctx, s)
	case StepWaitPixel:
		return e.waitPixel(ctx, s)
	}
	return nil
}

// Thời gian giữ khi bấm phím: chord giữ cả tổ hợp keyHoldMs; sequence nhấn mỗi phím keyHoldMs,
// cách phím sau keyGapMs.
const keyHoldMs, keyGapMs = 50, 50

// pressKeys bấm tổ hợp theo KeyMode; mọi phím đã nhấn luôn được nhả (ngược thứ tự), kể cả khi lỗi / bị huỷ.
func pressKeys(ctx context.Context, in Input, s Step) error {
	var held []string
	defer func() {
		for i := len(held) - 1; i >= 0; i-- {
			in.KeyToggle(held[i], false, nil)
		}
	}()
	press := func(k string, mods []string) error {
		if err := in.KeyToggle(k, true, mods); err != nil {
			return err
		}
		held = append(held, k)
		return nil
	}

	for _, m := range s.Modifiers {
		if err := press(m, nil); err != nil {
			return err
		}
	}
	if s.KeyMode == KeySequence {
		for i, k := range s.Keys {
			if err := press(k, s.Modifiers); err != nil {
				return err
			}
			if err := sleep(ctx, keyHoldMs); err != nil {
				return err
			}
			held = held[:len(held)-1]
			if err := in.KeyToggle(k, false, nil); err != nil {
				return err
			}
			if i < len(s.Keys)-1 {
				if err := sleep(ctx, keyGapMs); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for _, k := range s.Keys {
		if err := press(k, s.Modifiers); err != nil {
			return err
		}
	}
	return sleep(ctx, keyHoldMs)
}

// drag giữ nút tại A, di dần tới B trong DragMs, rồi nhả.
// A: points[0] hoặc chỗ chuột dừng ở bước trước (DragFromPrev);
// B: points[1] hoặc A + points[1] (DragRelative) – cả hai đều lệch ngẫu nhiên ± JitterPx.
func (e *engine) drag(ctx context.Context, s Step, pt func(int) (int, int)) error {
	in := e.in
	ax, ay := pt(0)
	if s.DragFromPrev && e.hasLast {
		ax, ay = e.last.X, e.last.Y
	}
	bx, by := pt(1)
	if s.DragRelative {
		d := s.Points[1]
		bx, by = e.jitter(ax+d.X, s.JitterPx), e.jitter(ay+d.Y, s.JitterPx)
	}
	if err := e.move(ax, ay); err != nil {
		return err
	}
	if err := in.Toggle(s.Button, true); err != nil {
		return err
	}
	defer in.Toggle(s.Button, false)

	n := max(s.DragMs/dragStepMs, 1)
	for i := 1; i <= n; i++ {
		if err := sleep(ctx, s.DragMs/n); err != nil {
			return err
		}
		x := ax + (bx-ax)*i/n
		y := ay + (by-ay)*i/n
		if err := in.Drag(x, y, s.Button); err != nil {
			return err
		}
		e.last = Point{X: x, Y: y}
	}
	return nil
}

// sleep chờ ms mili giây hoặc tới khi ctx bị huỷ.
func sleep(ctx context.Context, ms int) error {
	if ms <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(time.Duration(ms) * time.Millisecond)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// pauseGate – cổng tạm dừng: wait() chặn khi đang tạm dừng, tới khi resume hoặc ctx bị huỷ.
// Ghi lại tổng thời gian tạm dừng để tính thời gian chạy thật.
type pauseGate struct {
	mu        sync.Mutex
	ch        chan struct{} // khác nil khi đang tạm dừng; đóng khi resume
	since     time.Time     // lúc bắt đầu lần tạm dừng hiện tại
	pausedFor time.Duration // tổng thời gian đã tạm dừng (không tính lần đang dừng)
}

func (g *pauseGate) pause() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.ch != nil {
		return false
	}
	g.ch = make(chan struct{})
	g.since = time.Now()
	return true
}

func (g *pauseGate) resume() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.ch == nil {
		return false
	}
	close(g.ch)
	g.ch = nil
	g.pausedFor += time.Since(g.since)
	return true
}

func (g *pauseGate) paused() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.ch != nil
}

func (g *pauseGate) wait(ctx context.Context) error {
	g.mu.Lock()
	ch := g.ch
	g.mu.Unlock()
	if ch == nil {
		return ctx.Err()
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}

// pausedTotal – tổng thời gian đã tạm dừng tính tới giờ (kể cả lần đang dừng).
func (g *pauseGate) pausedTotal() time.Duration {
	g.mu.Lock()
	defer g.mu.Unlock()
	d := g.pausedFor
	if g.ch != nil {
		d += time.Since(g.since)
	}
	return d
}

// activeSince – thời gian từ start, trừ các khoảng tạm dừng.
func (g *pauseGate) activeSince(start time.Time) time.Duration {
	g.mu.Lock()
	defer g.mu.Unlock()
	d := time.Since(start) - g.pausedFor
	if g.ch != nil {
		d -= time.Since(g.since)
	}
	return d
}

// checkImageStep: bước chờ ảnh / click theo ảnh tự tìm có ảnh mẫu đọc được (hoặc chữ cần tìm) và thông số hợp lệ.
func checkImageStep(i int, s Step) error {
	if s.FindBy == FindByText {
		if strings.TrimSpace(s.FindText) == "" {
			return stepErr(i, s, "findText", "chưa nhập chữ cần tìm")
		}
		if s.Occurrence < 0 {
			return stepErr(i, s, "occurrence", "chỗ thứ mấy phải ≥ 1")
		}
		if s.TimeoutMs <= 0 || s.IntervalMs <= 0 {
			return stepErr(i, s, "timeoutMs", "thời gian chờ tối đa / kiểm tra mỗi phải > 0")
		}
		return nil
	}
	if s.Template == "" {
		return stepErr(i, s, "template", "chưa có ảnh mẫu – bấm “Khoanh ảnh mẫu”")
	}
	img, err := decodePNG(s.Template)
	if err != nil {
		return stepErr(i, s, "template", "%v", err)
	}
	if b := img.Bounds(); b.Dx() < 4 || b.Dy() < 4 {
		return stepErr(i, s, "template", "ảnh mẫu quá nhỏ (%d × %d px)", b.Dx(), b.Dy())
	}
	if s.Threshold < 1 || s.Threshold > 100 {
		return stepErr(i, s, "threshold", "độ khớp tối thiểu phải trong 1–100%%")
	}
	if s.TimeoutMs <= 0 || s.IntervalMs <= 0 {
		return stepErr(i, s, "timeoutMs", "thời gian chờ tối đa / kiểm tra mỗi phải > 0")
	}
	return nil
}
