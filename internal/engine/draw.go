package engine

import (
	"context"
	"math"

	"ghosthand/internal/humanize"
	"ghosthand/internal/model"
)

// Node Vẽ đường (bước 35 N7): giữ nút, kéo một nét liền qua các đỉnh trong DragMs – thời gian chia theo độ dài từng
// đoạn. Giống người: mỗi đoạn theo humanize.Path (cong ít để hình không méo, tăng / giảm tốc). Luôn nhả nút.

// drawCurve – phần độ cong giữ lại khi vẽ (đường cong giống người làm méo cạnh thẳng của hình).
const drawCurve = 0.15

func (e *Engine) drawPath(ctx context.Context, s model.Step) (map[string]any, error) {
	pts, err := e.drawPoints(s)
	if err != nil {
		return nil, err
	}
	if len(pts) < 2 {
		return nil, &invalidErr{"đường vẽ cần ít nhất 2 điểm"}
	}
	in := e.in
	if err := e.move(ctx, pts[0].X, pts[0].Y, 0); err != nil {
		return nil, err
	}
	if err := in.Toggle(s.Button, true); err != nil {
		return nil, err
	}
	defer in.Toggle(s.Button, false)
	total := 0.0
	for i := 1; i < len(pts); i++ {
		total += math.Hypot(float64(pts[i].X-pts[i-1].X), float64(pts[i].Y-pts[i-1].Y))
	}
	dur := max(s.DragMs, 1)
	drag := func(x, y int) error { return in.Drag(x, y, s.Button) }
	for i := 1; i < len(pts); i++ {
		a, b := pts[i-1], pts[i]
		seg := 1
		if total > 0 {
			seg = max(int(float64(dur)*math.Hypot(float64(b.X-a.X), float64(b.Y-a.Y))/total), 1)
		}
		if e.human != nil {
			p := *e.human
			p.Curve *= drawCurve
			if err := e.follow(ctx, humanize.Path(e.rng, p, a, b, humanize.Target{DurMs: seg}, e.moveBounds()), drag, drag); err != nil {
				return nil, err
			}
			continue
		}
		n := max(seg/dragStepMs, 1)
		for j := 1; j <= n; j++ {
			if err := Sleep(ctx, seg/n); err != nil {
				return nil, err
			}
			if err := e.gate.wait(ctx); err != nil {
				return nil, err
			}
			x, y := a.X+(b.X-a.X)*j/n, a.Y+(b.Y-a.Y)*j/n
			if err := drag(x, y); err != nil {
				return nil, err
			}
			e.last, e.hasLast = model.Point{X: x, Y: y}, true
		}
	}
	end := pts[len(pts)-1]
	return map[string]any{"x": end.X, "y": end.Y, "points": len(pts)}, nil
}

// drawPoints – các đỉnh theo toạ độ màn hình: đường gấp khúc qua Points (Closed: thêm điểm đầu), hoặc hình quanh tâm.
// Mỗi đỉnh lệch ngẫu nhiên ± JitterPx (hình: lệch cả hình theo tâm).
func (e *Engine) drawPoints(s model.Step) ([]model.Point, error) {
	if s.Shape == model.ShapeLine {
		out := make([]model.Point, 0, len(s.Points)+1)
		for _, p := range s.Points {
			q, err := e.pointAt(s, p)
			if err != nil {
				return nil, err
			}
			out = append(out, model.Point{X: e.jitter(q.X, s.JitterPx), Y: e.jitter(q.Y, s.JitterPx)})
		}
		if s.Closed && len(out) > 1 {
			out = append(out, out[0])
		}
		return out, nil
	}
	c, err := e.pointAt(s, s.Points[0])
	if err != nil {
		return nil, err
	}
	cx, cy := float64(e.jitter(c.X, s.JitterPx)), float64(e.jitter(c.Y, s.JitterPx))
	var out []model.Point
	for _, p := range model.ShapePoints(s.Shape, cx, cy, float64(s.ShapeW), float64(s.ShapeH)) {
		q := model.Point{X: int(math.Round(p.X)), Y: int(math.Round(p.Y))}
		if len(out) == 0 || out[len(out)-1] != q {
			out = append(out, q)
		}
	}
	return out, nil
}
