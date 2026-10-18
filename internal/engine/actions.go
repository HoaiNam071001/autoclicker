package engine

import (
	"context"
	"errors"

	"autoclicker/internal/device"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Node hành động: click, kéo, gõ, phím, cuộn, mở, chờ ảnh / pixel (imagestep.go) – thao tác qua device.Input.

// dragStepMs – khoảng thời gian giữa các lần di chuột khi kéo thả.
const dragStepMs = 10

// scrollGapMs – nghỉ giữa các nấc cuộn (app kịp nhận từng nấc; dừng được giữa chừng).
const scrollGapMs = 30

func (e *Engine) click(button string, double bool) error {
	if err := e.in.Click(button, double); err != nil {
		return err
	}
	e.update(false, func(p *Progress) { p.Clicks++ })
	return nil
}

// move đưa chuột tới (x, y) và nhớ lại vị trí đó.
func (e *Engine) move(x, y int) error {
	e.last, e.hasLast = model.Point{X: x, Y: y}, true
	return e.in.Move(x, y)
}

// runStep chạy hành động của node; output là giá trị node trả về cho node sau ({{nX.output.key}}).
func (e *Engine) runStep(ctx context.Context, s model.Step) (map[string]any, error) {
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
func (e *Engine) doStep(ctx context.Context, s model.Step) error {
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
				if err := Sleep(ctx, e.jitterMs(s.ClickInterval, s.JitterMs)); err != nil {
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
		return Sleep(ctx, s.HoldMs)
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
			if err := Sleep(ctx, scrollGapMs); err != nil {
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
			if err := Sleep(ctx, keyHoldMs); err != nil {
				return err
			}
			held = held[:len(held)-1]
			if err := in.KeyToggle(k, false, nil); err != nil {
				return err
			}
			if i < len(s.Keys)-1 {
				if err := Sleep(ctx, keyGapMs); err != nil {
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
	return Sleep(ctx, keyHoldMs)
}

// drag giữ nút tại A, di dần tới B trong DragMs, rồi nhả.
// A: points[0] hoặc chỗ chuột dừng ở bước trước (DragFromPrev);
// B: points[1] hoặc A + points[1] (DragRelative) – cả hai đều lệch ngẫu nhiên ± JitterPx.
func (e *Engine) drag(ctx context.Context, s model.Step, pt func(int) (int, int)) error {
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
		if err := Sleep(ctx, s.DragMs/n); err != nil {
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
