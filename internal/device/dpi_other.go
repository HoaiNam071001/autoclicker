//go:build !windows

package device

import "ghosthand/internal/model"

// DPIScale: macOS toạ độ theo điểm – độ nét Retina đã nằm ở tỷ lệ ảnh chụp / vùng.
func DPIScale(model.Rect) float64 { return 1 }
