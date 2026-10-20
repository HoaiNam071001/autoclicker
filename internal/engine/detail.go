package engine

import (
	"errors"

	"autoclicker/internal/model"
)

// Chi tiết lỗi có cấu trúc (StepError.Detail – bước 31b): AI / người đọc để tự sửa node thay vì đoán từ câu chữ.
//   chờ / click ảnh: {kind: image, best, threshold, bestBox, area}
//   chờ / click chữ: {kind: text, want, read[], area}
//   chờ màu:         {kind: pixel, point, actual, want, tolerance}
//   mọi node lỗi:    context {foreground {title, app, box}, mouse {x, y}, screen, displays, shotArea}

// errDetail – phần chi tiết riêng của lỗi (nil = không có).
func errDetail(err error) map[string]any {
	var nf *notFoundErr
	if errors.As(err, &nf) && nf.detail != nil {
		return nf.detail
	}
	return nil
}

func imageMiss(m FindResult, threshold int) map[string]any {
	return map[string]any{"kind": "image", "best": m.Score, "threshold": threshold,
		"bestBox": model.Rect{X: m.X, Y: m.Y, W: m.W, H: m.H}, "area": m.Area}
}

func textMiss(want string, read []string, area model.Rect) map[string]any {
	if read == nil {
		read = []string{}
	}
	return map[string]any{"kind": "text", "want": want, "read": read, "area": area}
}

func pixelMiss(p model.Point, actual, want string, tolerance int) map[string]any {
	return map[string]any{"kind": "pixel", "point": p, "actual": actual, "want": want, "tolerance": tolerance}
}

// addContext gắn ngữ cảnh lúc node lỗi: cửa sổ đang ở trước, vị trí chuột, màn hình. Đọc lỗi thì bỏ phần đó.
func (e *Engine) addContext(res *model.StepResult) {
	if res.Error == nil || res.Error.Code == model.ErrCancelled {
		return
	}
	// shotArea – vùng của ảnh chụp lúc lỗi (shots.go): toạ độ thật = shotArea.x + px ảnh × shotArea.w / rộng ảnh
	ctx := map[string]any{"screen": e.screen, "shotArea": e.searchArea(model.Step{SearchArea: "screen"})}
	if len(e.displays) > 0 {
		ctx["displays"] = e.displays
	}
	if e.in != nil {
		x, y := e.in.Location()
		ctx["mouse"] = model.Point{X: x, Y: y}
	}
	if ws, err := e.wins.List(); err == nil {
		for _, w := range ws { // List theo thứ tự chồng – cửa sổ không thu nhỏ đầu tiên là cửa sổ đang ở trước
			if !w.Minimized {
				ctx["foreground"] = map[string]any{"title": w.Title, "app": w.App, "box": w.Box}
				break
			}
		}
	}
	if res.Error.Detail == nil {
		res.Error.Detail = map[string]any{}
	}
	res.Error.Detail["context"] = ctx
}
