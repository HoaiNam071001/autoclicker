package engine

import (
	"context"
	"errors"

	"autoclicker/internal/launch"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Node hành động: click, kéo, gõ, phím, cuộn, mở, chờ ảnh / pixel (imagestep.go), web (webstep.go) – thao tác qua device.Input.

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

// pointAt – toạ độ màn hình của điểm (neo theo vùng / màn hình, hoặc tính từ biến đã thay – F5b).
func (e *Engine) pointAt(p model.Point) (model.Point, error) {
	q, err := script.PointOnScreen(p, script.Bounds(e.sc, e.screen))
	if err != nil {
		return q, &invalidErr{err.Error()}
	}
	return q, nil
}

// runStep chạy hành động của node; output là giá trị node trả về cho node sau ({{nX.output.key}}).
// wait: thời gian chờ sau node (delayAfter đã tính) – output ms của node Chờ.
func (e *Engine) runStep(ctx context.Context, s model.Step, wait int) (map[string]any, error) {
	switch s.Type {
	case model.StepWaitImage:
		h, err := e.findImage(ctx, s)
		if err != nil {
			return nil, err
		}
		return h.output(), nil
	case model.StepClickImage:
		return e.clickImage(ctx, s)
	case model.StepWebWait:
		return e.webWait(ctx, s)
	case model.StepWebClick:
		return e.webClick(ctx, s)
	case model.StepWebType:
		return e.webType(ctx, s)
	case model.StepWebRead:
		return e.webRead(ctx, s)
	case model.StepWebScroll:
		return e.webScroll(ctx, s)
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
	return e.actionOutput(s, wait), nil
}

// doStep – hành động của các node không tìm ảnh (chuột, phím, cuộn, mở).
func (e *Engine) doStep(ctx context.Context, s model.Step) error {
	in := e.in
	// pt: toạ độ màn hình của điểm i (theo neo / biến), mỗi lần gọi lệch ngẫu nhiên ± JitterPx.
	pt := func(i int) (int, int, error) {
		p, err := e.pointAt(s.Points[i])
		return e.jitter(p.X, s.JitterPx), e.jitter(p.Y, s.JitterPx), err
	}

	switch s.Type {
	case model.StepClick:
		n := max(s.Clicks, 1)
		for i := 0; i < n; i++ {
			if err := e.gate.wait(ctx); err != nil {
				return err
			}
			x, y, err := pt(0) // mỗi click lệch riêng
			if err != nil {
				return err
			}
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
		x, y, err := pt(0)
		if err != nil {
			return err
		}
		if err := e.move(x, y); err != nil {
			return err
		}
		return e.click(s.Button, true)
	case model.StepHold:
		x, y, err := pt(0)
		if err != nil {
			return err
		}
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
		return e.typeText(ctx, s.Text)
	case model.StepKey:
		return pressKeys(ctx, in, s)
	case model.StepScroll:
		if len(s.Points) > 0 {
			x, y, err := pt(0)
			if err != nil {
				return err
			}
			if err := e.move(x, y); err != nil {
				return err
			}
		}
		return e.wheel(ctx, s.Direction, s.Amount)
	case model.StepOpen:
		if s.WebBrowser {
			return e.openInBrowser(ctx, launch.BrowserURL(s.Target))
		}
		return e.open(s.Target, s.App)
	}
	return nil
}

// typeText gõ từng ký tự để dừng / tạm dừng được giữa chừng với chuỗi dài.
func (e *Engine) typeText(ctx context.Context, text string) error {
	for _, r := range text {
		if err := e.gate.wait(ctx); err != nil {
			return err
		}
		if err := e.in.TypeStr(string(r)); err != nil {
			return err
		}
	}
	return nil
}

// wheel lăn n nấc theo dir ở chỗ chuột đang đứng.
func (e *Engine) wheel(ctx context.Context, dir string, n int) error {
	return e.wheelEach(ctx, n, func() error { return e.in.Scroll(dir) })
}

// wheelEach lăn n nấc bằng notch, cách nhau scrollGapMs, dừng / tạm dừng được giữa các nấc.
func (e *Engine) wheelEach(ctx context.Context, n int, notch func() error) error {
	for i := 0; i < n; i++ {
		if err := e.gate.wait(ctx); err != nil {
			return err
		}
		if err := notch(); err != nil {
			return err
		}
		if err := Sleep(ctx, scrollGapMs); err != nil {
			return err
		}
	}
	return nil
}

// drag giữ nút tại A, di dần tới B trong DragMs, rồi nhả.
// A: points[0] hoặc chỗ chuột dừng ở bước trước (DragFromPrev);
// B: points[1] hoặc A + points[1] (DragRelative – điểm dùng biến thì biến là độ lệch) – cả hai lệch ngẫu nhiên ± JitterPx.
func (e *Engine) drag(ctx context.Context, s model.Step, pt func(int) (int, int, error)) error {
	in := e.in
	var ax, ay, bx, by int
	var err error
	if s.DragFromPrev {
		if !e.hasLast { // đồ thị không biết trước node nào chạy trước – kiểm tra lúc chạy
			return errors.New("kéo từ điểm trước đó nhưng chưa node nào đưa chuột tới đâu")
		}
		ax, ay = e.last.X, e.last.Y
	} else if ax, ay, err = pt(0); err != nil {
		return err
	}
	if s.DragRelative {
		d := s.Points[1]
		if d.UsesVar() {
			if d, err = script.PointOnScreen(d, model.Rect{}); err != nil { // độ lệch: không neo, khung không dùng
				return &invalidErr{"điểm B: " + err.Error()}
			}
		}
		bx, by = e.jitter(ax+d.X, s.JitterPx), e.jitter(ay+d.Y, s.JitterPx)
	} else if bx, by, err = pt(1); err != nil {
		return err
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
