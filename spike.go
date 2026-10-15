package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

// SpikeResult – kết quả spike click (bước 3), trả về cho frontend.
type SpikeResult struct {
	X int `json:"x"` // vị trí chuột thực tế sau khi di
	Y int `json:"y"`
}

// SpikeClick đợi delayMs, di chuột tới (x, y) rồi click trái.
// Dùng để kiểm tra robotgo click được ra ngoài cửa sổ app. Sẽ bỏ khi có engine (bước 7).
func (a *App) SpikeClick(x, y, delayMs int) (SpikeResult, error) {
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	if err := robotgo.Move(x, y); err != nil {
		return SpikeResult{}, fmt.Errorf("không di chuyển được chuột: %w", err)
	}
	time.Sleep(50 * time.Millisecond)

	cx, cy := robotgo.Location()
	res := SpikeResult{X: cx, Y: cy}
	// Toạ độ lệch → thường do chưa cấp quyền Accessibility (macOS).
	if abs(cx-x) > 2 || abs(cy-y) > 2 {
		return res, fmt.Errorf("chuột ở (%d, %d) thay vì (%d, %d) – kiểm tra quyền Accessibility", cx, cy, x, y)
	}

	if err := robotgo.Click("left"); err != nil {
		return res, fmt.Errorf("không click được: %w", err)
	}
	return res, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
