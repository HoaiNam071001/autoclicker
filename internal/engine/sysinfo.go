package engine

import (
	"context"

	"ghosthand/internal/model"
)

// Node Thông tin máy (bước 43 S4): đọc một lần toàn bộ thông số, node sau chọn cái cần qua output.
// FilePath = ổ / thư mục đo chỗ trống.
func (e *Engine) sysInfo(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	return e.sysinfo.Read(ctx, s.FilePath)
}
