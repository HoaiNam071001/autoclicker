package engine

import (
	"strings"

	"autoclicker/internal/model"
)

// actionOutput – output của node chuột / phím / chờ / mở sau khi chạy xong (khớp Outputs trong model/steptypes.go,
// flow/view.ts outputsOf). s đã thay biến + numVars nên output là giá trị thật đã dùng.
func (e *Engine) actionOutput(s model.Step, wait int) map[string]any {
	switch s.Type {
	case model.StepClick, model.StepDouble, model.StepDrag:
		return map[string]any{"x": e.last.X, "y": e.last.Y} // chỗ chuột dừng
	case model.StepHold:
		return map[string]any{"x": e.last.X, "y": e.last.Y, "ms": s.HoldMs}
	case model.StepWait:
		return map[string]any{"ms": wait}
	case model.StepText:
		return map[string]any{"text": s.Text}
	case model.StepKey:
		return map[string]any{"keys": strings.Join(append(append([]string{}, s.Modifiers...), s.Keys...), "+")}
	case model.StepScroll:
		out := map[string]any{"direction": s.Direction, "amount": s.Amount}
		if e.hasLast { // không có điểm và chưa node chuột nào chạy: không biết chuột ở đâu – không có x, y
			out["x"], out["y"] = e.last.X, e.last.Y
		}
		return out
	case model.StepOpen:
		app := s.App
		if s.WebBrowser {
			app = "Chrome của app"
		}
		return map[string]any{"target": s.Target, "app": app}
	}
	return nil
}
