package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/shortcuts"
)

// Node Shortcuts của hệ điều hành (bước 43 S2, macOS): chạy Shortcut trong ứng dụng Phím tắt bằng `shortcuts run`,
// input là chữ ở Text (nhận biến), output `text` là chữ Shortcut trả về. Tên và input đi qua đối số / file tạm,
// không qua shell. Hết giờ / Dừng → dừng lệnh.
func (e *Engine) osShortcut(ctx context.Context, s model.Step) (map[string]any, error) {
	name := strings.TrimSpace(s.ShortcutName)
	if name == "" {
		return nil, &invalidErr{"chưa có tên Shortcut"}
	}
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	timeout := s.TimeoutMs
	if timeout <= 0 {
		timeout = model.OSShortcutTimeoutMs
	}
	cctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()
	text, err := e.shortcuts.Run(cctx, name, s.Text)
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case errors.Is(err, context.DeadlineExceeded):
		return nil, &timeoutErr{msg: fmt.Sprintf("Shortcut %q chạy quá %s – đã dừng", name, fmtSeconds(timeout)), detail: map[string]any{"shortcut": name, "timeoutMs": timeout}}
	case errors.Is(err, shortcuts.ErrUnsupported):
		return nil, &invalidErr{"node Shortcuts chỉ chạy trên macOS"}
	case errors.Is(err, shortcuts.ErrNotFound):
		return nil, &notFoundErr{msg: fmt.Sprintf("không có Shortcut tên %q trong ứng dụng Phím tắt", name), detail: map[string]any{"kind": "shortcut", "name": name}}
	case err != nil:
		return nil, err
	}
	return map[string]any{"text": trimEOL(text)}, nil
}
