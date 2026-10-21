package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ghosthand/internal/model"
)

// Node Clipboard (bước 35 N2): đọc chữ / ghi chữ / ghi ảnh PNG từ file. Không đụng chuột / phím – copy / paste trọn gói
// thì ghép với node Phím tắt (Cmd/Ctrl+C, Cmd/Ctrl+V).
func (e *Engine) clipboard(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	switch s.ClipOp {
	case model.ClipWrite:
		if err := e.clip.WriteText(s.Text); err != nil {
			return nil, err
		}
		return map[string]any{"text": s.Text}, nil
	case model.ClipImage:
		p, err := userPath(s.FilePath)
		if err != nil {
			return nil, err
		}
		if !strings.EqualFold(filepath.Ext(p), ".png") {
			return nil, &invalidErr{fmt.Sprintf("chỉ chép được ảnh PNG, file %q không phải .png", filepath.Base(p))}
		}
		if _, err := os.Stat(p); err != nil {
			return nil, &invalidErr{fmt.Sprintf("không thấy file ảnh %s", p)}
		}
		if err := e.clip.WriteImage(p); err != nil {
			return nil, err
		}
		return map[string]any{"path": p}, nil
	}
	t, err := e.clip.ReadText()
	if err != nil {
		return nil, err
	}
	return map[string]any{"text": t}, nil
}
