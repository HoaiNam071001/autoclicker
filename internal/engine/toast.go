package engine

import (
	"context"
	"errors"
	"strings"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/toast"
)

// Node Thông báo hệ thống (bước 43 S9): hiện thông báo của hệ điều hành với tiêu đề + nội dung (nhận biến). Không
// chờ người dùng bấm – nút bấm / cổng theo nút là S9b.
func (e *Engine) notifyStep(ctx context.Context, s model.Step) (map[string]any, error) {
	title, text := strings.TrimSpace(s.NotifyTitle), strings.TrimSpace(s.Text)
	if title == "" && text == "" {
		return nil, &invalidErr{"thông báo trống – nhập tiêu đề hoặc nội dung"}
	}
	if title == "" {
		title = "Ghosthand"
	}
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := e.toasts.Show(cctx, title, text)
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case errors.Is(err, toast.ErrUnsupported):
		return nil, &invalidErr{err.Error()}
	case err != nil:
		return nil, err
	}
	return map[string]any{"title": title, "text": text}, nil
}
