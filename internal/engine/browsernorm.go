package engine

import (
	"context"

	"ghosthand/internal/model"
)

// Node Chuẩn hoá trình duyệt (bước 43 B3): đưa cửa sổ trình duyệt thường về trạng thái cố định để toạ độ / ảnh mẫu
// không lệch – thoát toàn màn hình, đặt cỡ (preset như node Cửa sổ), zoom 100 %, đóng popup. Chọn tab như browserTab.
// Zoom và popup là phím gửi cho cửa sổ trước nhất nên node đưa tab lên trước; trình duyệt không cho đọc mức zoom nên
// output chỉ báo đã gửi phím.

func (e *Engine) browserNorm(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	t, err := e.waitTab(ctx, s)
	if err != nil {
		return nil, err
	}
	if err := e.tabs.Activate(ctx, t); err != nil {
		return nil, err
	}
	o := model.NormalizeOpts{Zoom: s.NormZoom, Fullscreen: s.NormFullscreen, Popups: s.NormPopups}
	if s.NormSize {
		cur, err := e.tabs.Frame(ctx, t)
		if err != nil {
			return nil, err
		}
		r, err := e.windowTarget(s, model.WindowInfo{Box: cur.Box})
		if err != nil {
			return nil, err
		}
		o.Bounds = &r
	}
	was, err := e.tabs.Normalize(ctx, t, o)
	if err != nil {
		return nil, err
	}
	fr, err := e.tabs.Frame(ctx, t)
	if err != nil {
		return nil, err
	}
	b := fr.Box
	out := tabOutput(t)
	out["wasFullscreen"], out["zoomReset"] = was, o.Zoom
	out["left"], out["top"], out["w"], out["h"] = b.X, b.Y, b.W, b.H
	out["x"], out["y"] = b.X+b.W/2, b.Y+b.H/2
	return out, nil
}
