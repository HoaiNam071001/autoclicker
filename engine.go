package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// dragStepMs – khoảng thời gian giữa các lần di chuột khi kéo thả.
const dragStepMs = 10

// validateScript kiểm tra kịch bản trước khi chạy để không dừng giữa chừng vì dữ liệu sai.
func validateScript(sc Script) error {
	enabled := 0
	for i, s := range sc.Steps {
		if !s.Enabled {
			continue
		}
		enabled++
		need := 0
		switch s.Type {
		case StepClick, StepDouble, StepHold:
			need = 1
		case StepDrag:
			need = 2
		case StepWait:
		case StepText:
			if s.Text == "" {
				return fmt.Errorf("bước %d: chưa nhập nội dung gõ", i+1)
			}
		case StepKey:
			if len(s.Keys) == 0 {
				return fmt.Errorf("bước %d: chưa chọn phím", i+1)
			}
		default:
			return fmt.Errorf("bước %d: loại bước không hợp lệ %q", i+1, s.Type)
		}
		if len(s.Points) < need {
			return fmt.Errorf("bước %d: cần %d điểm, mới có %d", i+1, need, len(s.Points))
		}
	}
	if enabled == 0 {
		return errors.New("kịch bản không có bước nào được bật")
	}
	return nil
}

// Progress – tiến độ gửi sang frontend (event "engine:progress").
type Progress struct {
	Loop      int   `json:"loop"`  // vòng hiện tại, từ 1
	Loops     int   `json:"loops"` // 0 = vô hạn
	Step      int   `json:"step"`  // thứ tự bước hiện tại trong các bước đang bật, từ 1
	Total     int   `json:"total"` // số bước đang bật
	Clicks    int   `json:"clicks"`
	ElapsedMs int64 `json:"elapsedMs"` // không tính thời gian tạm dừng
	Paused    bool  `json:"paused"`
}

// progressEvery – giãn cách tối thiểu giữa 2 lần báo tiến độ trong cùng một bước (tránh dội event).
const progressEvery = 100 * time.Millisecond

// engine chạy một kịch bản. Tạo mới cho mỗi lần chạy.
type engine struct {
	in         Input
	sc         Script
	gate       pauseGate
	onProgress func(Progress) // có thể nil

	mu       sync.Mutex
	prog     Progress
	start    time.Time
	lastEmit time.Time
}

func newEngine(in Input, sc Script, onProgress func(Progress)) *engine {
	return &engine{in: in, sc: sc, onProgress: onProgress}
}

// run chạy tới khi đủ số vòng, hết MaxMinutes hoặc ctx bị huỷ.
// Hết thời gian là kết thúc bình thường (nil); bị huỷ trả ctx.Err().
// MaxMinutes tính theo đồng hồ thật, kể cả lúc tạm dừng.
func (e *engine) run(ctx context.Context) error {
	sc := e.sc
	if err := validateScript(sc); err != nil {
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
				err = sleep(runCtx, s.DelayAfter)
			}
			if err != nil {
				if ctx.Err() == nil && errors.Is(err, context.DeadlineExceeded) {
					return nil // hết MaxMinutes
				}
				if ctx.Err() != nil {
					return ctx.Err()
				}
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

func (e *engine) runStep(ctx context.Context, s Step) error {
	in := e.in
	pt := func(i int) (int, int) {
		p := s.Points[i]
		if r := e.sc.Region; r != nil {
			return r.X + p.X, r.Y + p.Y
		}
		return p.X, p.Y
	}

	switch s.Type {
	case StepClick:
		x, y := pt(0)
		n := max(s.Clicks, 1)
		for i := 0; i < n; i++ {
			if err := e.gate.wait(ctx); err != nil {
				return err
			}
			if err := in.Move(x, y); err != nil {
				return err
			}
			if err := e.click(s.Button, false); err != nil {
				return err
			}
			if i < n-1 {
				if err := sleep(ctx, s.ClickInterval); err != nil {
					return err
				}
			}
		}
	case StepDouble:
		x, y := pt(0)
		if err := in.Move(x, y); err != nil {
			return err
		}
		return e.click(s.Button, true)
	case StepHold:
		x, y := pt(0)
		if err := in.Move(x, y); err != nil {
			return err
		}
		if err := in.Toggle(s.Button, true); err != nil {
			return err
		}
		defer in.Toggle(s.Button, false) // luôn nhả chuột, kể cả khi bị huỷ giữa chừng
		return sleep(ctx, s.HoldMs)
	case StepDrag:
		return drag(ctx, in, s, pt)
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
func drag(ctx context.Context, in Input, s Step, pt func(int) (int, int)) error {
	ax, ay := pt(0)
	bx, by := pt(1)
	if err := in.Move(ax, ay); err != nil {
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
