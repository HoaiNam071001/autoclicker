package main

import (
	"context"
	"errors"
	"image"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"autoclicker/internal/device"
	"autoclicker/internal/launch"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
	"autoclicker/internal/vision"
)

// dragStepMs – khoảng thời gian giữa các lần di chuột khi kéo thả.
const dragStepMs = 10

// scrollGapMs – nghỉ giữa các nấc cuộn (app kịp nhận từng nấc; dừng được giữa chừng).
const scrollGapMs = 30

// Progress – tiến độ gửi sang frontend (event "engine:progress").
type Progress struct {
	Loop      int           `json:"loop"`  // vòng hiện tại, từ 1
	Loops     int           `json:"loops"` // 0 = vô hạn
	Step      int           `json:"step"`  // số thứ tự node đang chạy (theo model.FlowOrder, trong các node đang bật), từ 1
	Total     int           `json:"total"` // số node đang bật
	NodeID    string        `json:"nodeId,omitempty"`
	Ref       string        `json:"ref,omitempty"`
	Clicks    int           `json:"clicks"`
	ElapsedMs int64         `json:"elapsedMs"` // không tính thời gian tạm dừng
	Paused    bool          `json:"paused"`
	Wait      *WaitProgress `json:"wait,omitempty"`    // đang chờ ảnh / pixel
	Retry     int           `json:"retry,omitempty"`   // đang thử lại lần thứ mấy (0 = lần đầu)
	Retries   int           `json:"retries,omitempty"` // số lần thử lại tối đa của bước hiện tại
	Break     string        `json:"break,omitempty"`   // đang dừng trước node NodeID vì debug: breakpoint | step
}

// progressEvery – giãn cách tối thiểu giữa 2 lần báo tiến độ trong cùng một bước (tránh dội event).
const progressEvery = 100 * time.Millisecond

// engine chạy một kịch bản. Tạo mới cho mỗi lần chạy.
type engine struct {
	in         device.Input
	sc         model.Script
	gate       pauseGate
	screen     model.Rect                     // khung màn hình khi Region = nil (runner đặt = mainScreen); zero = gốc (0, 0)
	displays   []model.Rect                   // khung các màn (tìm ảnh "cả màn", kiểm tra vùng chụp); nil = không kiểm tra
	onProgress func(Progress)                 // có thể nil
	intN       func(n int) int                // số ngẫu nhiên [0, n) – test thay bằng hàm cố định
	open       func(target, app string) error // bước Mở – test thay bằng hàm ghi lại
	scr        device.Screen                  // chụp màn hình / đọc màu (imagestep.go) – test thay bằng màn hình giả
	ocr        device.TextReader              // đọc chữ (tìm theo chữ, ocr.go) – test thay bằng bộ đọc giả

	onNode  func(model.StepResult)       // mỗi node chạy xong (event engine:node) – có thể nil
	shots   map[string]image.Image       // ảnh màn hình lúc node thất bại theo tên file (nil = không chụp – runlog.go)
	debug   bool                         // dừng ở node có Breakpoint (chạy từ app); không thì bỏ qua breakpoint
	step    atomic.Bool                  // chạy từng node: dừng trước mỗi node (StepScript bật, ResumeScript tắt)
	failed  string                       // id node làm lần chạy dừng vì lỗi (runner gửi engine:failed)
	miss    *FindResult                  // node chờ ảnh hết giờ: chỗ giống nhất (để xem lại vì sao không thấy)
	tpls    map[string]vision.Gray       // ảnh mẫu đã giải mã theo id node
	byID    map[string]model.Step        // node theo id (run đặt)
	order   map[string]int               // số thứ tự hiển thị của node (model.FlowOrder) – thanh mini
	results map[string]*model.StepResult // kết quả lần chạy gần nhất của từng node trong vòng này
	rv      runVars                      // bảng tra biến {{nX…}} / {{last…}} / {{vars…}}
	counts  map[string]int               // node Đếm: số lần đã đi qua trong vòng này

	last    model.Point // chỗ engine vừa đưa chuột tới – cho bước kéo "từ điểm trước đó"
	hasLast bool

	mu       sync.Mutex
	prog     Progress
	start    time.Time
	lastEmit time.Time
}

func newEngine(in device.Input, sc model.Script, onProgress func(Progress)) *engine {
	return &engine{in: in, sc: sc, onProgress: onProgress, intN: rand.IntN, open: launch.OpenTarget, scr: device.RobotgoScreen{}, ocr: device.SystemOCR{}}
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
	e.step.Store(false)
	if !e.gate.resume() {
		return false
	}
	e.update(true, func(p *Progress) { p.Break = "" })
	return true
}

// Step chạy node kế rồi dừng lại (debug, F10); đang chạy thì dừng trước node sau.
func (e *engine) Step() bool {
	e.step.Store(true)
	e.gate.resume()
	e.update(true, func(p *Progress) { p.Break = "" })
	return true
}

// breakBefore – debug: tạm dừng trước node s nếu đang chạy từng node / s có breakpoint (wait sau đó chặn tới khi tiếp).
func (e *engine) breakBefore(s model.Step, loop int) {
	why := ""
	switch {
	case !s.Enabled: // node tắt chỉ đi qua, không đáng dừng
		return
	case e.step.Load():
		why = "step"
	case e.debug && s.Breakpoint:
		why = "breakpoint"
	default:
		return
	}
	e.gate.pause()
	e.update(true, func(p *Progress) { p.Loop, p.Step, p.NodeID, p.Ref, p.Break = loop, e.order[s.ID], s.ID, s.Ref, why })
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
	e.last, e.hasLast = model.Point{X: x, Y: y}, true
	return e.in.Move(x, y)
}

// runStep chạy hành động của node; output là giá trị node trả về cho node sau ({{nX.output.key}}).
func (e *engine) runStep(ctx context.Context, s model.Step) (map[string]any, error) {
	switch s.Type {
	case model.StepWaitImage:
		h, err := e.findImage(ctx, s)
		if err != nil {
			return nil, err
		}
		return h.output(), nil
	case model.StepClickImage:
		return e.clickImage(ctx, s)
	case model.StepWaitPixel:
		c, err := e.waitPixel(ctx, s)
		if err != nil {
			return nil, err
		}
		return map[string]any{"color": c}, nil
	}
	if err := e.doStep(ctx, s); err != nil {
		return nil, err
	}
	switch s.Type {
	case model.StepClick, model.StepDouble, model.StepHold, model.StepDrag:
		return map[string]any{"x": e.last.X, "y": e.last.Y}, nil // chỗ chuột dừng
	case model.StepText:
		return map[string]any{"text": s.Text}, nil
	case model.StepOpen:
		return map[string]any{"target": s.Target}, nil
	}
	return nil, nil
}

// doStep – hành động của các node không tìm ảnh (chuột, phím, cuộn, mở).
func (e *engine) doStep(ctx context.Context, s model.Step) error {
	in := e.in
	// pt: toạ độ màn hình của điểm i (theo neo), mỗi lần gọi lệch ngẫu nhiên ± JitterPx.
	pt := func(i int) (int, int) {
		p := script.ResolvePoint(s.Points[i], script.Bounds(e.sc, e.screen))
		return e.jitter(p.X, s.JitterPx), e.jitter(p.Y, s.JitterPx)
	}

	switch s.Type {
	case model.StepClick:
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
	case model.StepDouble:
		x, y := pt(0)
		if err := e.move(x, y); err != nil {
			return err
		}
		return e.click(s.Button, true)
	case model.StepHold:
		x, y := pt(0)
		if err := e.move(x, y); err != nil {
			return err
		}
		if err := in.Toggle(s.Button, true); err != nil {
			return err
		}
		defer in.Toggle(s.Button, false) // luôn nhả chuột, kể cả khi bị huỷ giữa chừng
		return sleep(ctx, s.HoldMs)
	case model.StepDrag:
		return e.drag(ctx, s, pt)
	case model.StepWait:
		// chỉ chờ DelayAfter
	case model.StepText:
		// Gõ từng ký tự để dừng / tạm dừng được giữa chừng với chuỗi dài.
		for _, r := range s.Text {
			if err := e.gate.wait(ctx); err != nil {
				return err
			}
			if err := in.TypeStr(string(r)); err != nil {
				return err
			}
		}
	case model.StepKey:
		return pressKeys(ctx, in, s)
	case model.StepScroll:
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
	case model.StepOpen:
		return e.open(s.Target, s.App)
	}
	return nil
}

// Thời gian giữ khi bấm phím: chord giữ cả tổ hợp keyHoldMs; sequence nhấn mỗi phím keyHoldMs,
// cách phím sau keyGapMs.
const keyHoldMs, keyGapMs = 50, 50

// pressKeys bấm tổ hợp theo KeyMode; mọi phím đã nhấn luôn được nhả (ngược thứ tự), kể cả khi lỗi / bị huỷ.
func pressKeys(ctx context.Context, in device.Input, s model.Step) error {
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
	if s.KeyMode == model.KeySequence {
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
func (e *engine) drag(ctx context.Context, s model.Step, pt func(int) (int, int)) error {
	in := e.in
	ax, ay := pt(0)
	if s.DragFromPrev {
		if !e.hasLast { // đồ thị không biết trước node nào chạy trước – kiểm tra lúc chạy
			return errors.New("kéo từ điểm trước đó nhưng chưa node nào đưa chuột tới đâu")
		}
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
		e.last = model.Point{X: x, Y: y}
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
