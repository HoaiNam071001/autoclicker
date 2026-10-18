package engine

import (
	"context"

	"autoclicker/internal/humanize"
	"autoclicker/internal/model"
)

// Thao tác giống người (bước 34): mức theo node (Step.Human, trống = theo Script.Human) đặt ở đầu mỗi node (setHuman);
// e.human nil = tắt – chuột nhảy thẳng, click / gõ / cuộn như cũ.

func (e *Engine) setHuman(s model.Step) {
	e.human = nil
	if p, ok := humanize.ProfileOf(humanize.Effective(e.sc.Human, s.Human)); ok {
		e.human = &p
	}
}

// moveBounds – khung bao các màn hình: điểm giữa đường di chuột không ra ngoài / chạm góc dừng khẩn cấp.
func (e *Engine) moveBounds() model.Rect {
	if len(e.displays) == 0 {
		return e.screen
	}
	b := e.displays[0]
	for _, d := range e.displays[1:] {
		x0, y0 := min(b.X, d.X), min(b.Y, d.Y)
		x1, y1 := max(b.X+b.W, d.X+d.W), max(b.Y+b.H, d.Y+d.H)
		b = model.Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}
	}
	return b
}

// glide đưa chuột theo đường giống người tới (x, y); size = cỡ đích (Fitts), 0 = không rõ.
func (e *Engine) glide(ctx context.Context, x, y, size int) error {
	fx, fy := e.in.Location()
	path := humanize.Path(e.rng, *e.human, model.Point{X: fx, Y: fy}, model.Point{X: x, Y: y}, humanize.Target{Size: size}, e.moveBounds())
	return e.follow(ctx, path, e.in.MoveStep, e.in.Move)
}

// follow đi hết path: chờ từng bước (dừng / tạm dừng được giữa chừng), bước giữa bằng step, bước cuối bằng last.
func (e *Engine) follow(ctx context.Context, path []humanize.Move, step, last func(x, y int) error) error {
	for i, m := range path {
		if err := Sleep(ctx, m.WaitMs); err != nil {
			return err
		}
		if err := e.gate.wait(ctx); err != nil {
			return err
		}
		mv := step
		if i == len(path)-1 {
			mv = last
		}
		if err := mv(m.X, m.Y); err != nil {
			return err
		}
		e.last, e.hasLast = model.Point{X: m.X, Y: m.Y}, true
	}
	return nil
}

// humanClick – nhấn, giữ ClickHold rồi nhả (luôn nhả, kể cả khi bị huỷ lúc đang giữ).
func (e *Engine) humanClick(ctx context.Context, button string) error {
	if err := e.in.Toggle(button, true); err != nil {
		return err
	}
	err := Sleep(ctx, humanize.Between(e.rng, e.human.ClickHold))
	if up := e.in.Toggle(button, false); err == nil {
		err = up
	}
	return err
}

// keyTiming – thời gian giữ / nghỉ giữa các phím của node phím tắt.
func (e *Engine) keyTiming() (hold, gap func() int) {
	if e.human == nil {
		fixed := func(ms int) func() int { return func() int { return ms } }
		return fixed(keyHoldMs), fixed(keyGapMs)
	}
	p := e.human
	return func() int { return humanize.Between(e.rng, p.KeyHold) }, func() int { return humanize.Between(e.rng, p.KeyGap) }
}

// aim – điểm bấm của đích tâm c cỡ r: giống người và không lệch tay (dx = dy = 0) thì rơi ngẫu nhiên quanh tâm.
func (e *Engine) aim(c model.Point, r model.Rect, dx, dy int) model.Point {
	if e.human == nil || dx != 0 || dy != 0 {
		return c
	}
	x, y := humanize.Aim(e.rng, c.X, c.Y, r.W, r.H)
	return model.Point{X: x, Y: y}
}
