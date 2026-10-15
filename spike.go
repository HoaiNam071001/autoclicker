package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// spikeMaxClicks – giới hạn an toàn cho vòng lặp thử, phòng khi F8 không hoạt động.
const spikeMaxClicks = 50

type spikeRun struct{ cancel context.CancelFunc }

// StartSpikeLoop click trái liên tục tại (x, y), mỗi intervalMs một lần, cho tới khi
// StopSpikeLoop / F8 hoặc đủ spikeMaxClicks. Emit "spike:state" sau mỗi click. Spike bước 4.
func (a *App) StartSpikeLoop(x, y, intervalMs int) {
	a.mu.Lock()
	if a.spike != nil {
		a.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	run := &spikeRun{cancel: cancel}
	a.spike = run
	a.mu.Unlock()

	go func() {
		clicks := 0
		defer func() {
			cancel()
			a.mu.Lock()
			if a.spike == run {
				a.spike = nil
			}
			a.mu.Unlock()
			runtime.EventsEmit(a.ctx, "spike:state", false, clicks)
		}()
		ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
		defer ticker.Stop()
		for clicks < spikeMaxClicks {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
			robotgo.Move(x, y)
			robotgo.Click("left")
			clicks++
			runtime.EventsEmit(a.ctx, "spike:state", true, clicks)
		}
	}()
}

// StopSpikeLoop dừng vòng lặp click thử (không làm gì nếu đang không chạy).
func (a *App) StopSpikeLoop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.spike != nil {
		a.spike.cancel()
		a.spike = nil
	}
}
