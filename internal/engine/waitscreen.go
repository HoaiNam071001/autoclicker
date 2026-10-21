package engine

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"ghosthand/internal/device"
	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// Node Chờ màn hình (bước 42 G1): chụp vùng (như node Chụp màn hình) mỗi IntervalMs, so ảnh xám thu nhỏ.
// Đứng yên: lệch giữa 2 lần chụp liền nhau ≤ DiffPct liên tục StableMs. Thay đổi: khác lần chụp đầu > DiffPct.

func (e *Engine) waitScreen(ctx context.Context, s model.Step) (map[string]any, error) {
	area, err := e.shotArea(s)
	if err != nil {
		return nil, err
	}
	if err := device.CaptureArea(area, e.displays); err != nil {
		return nil, err
	}
	start, paused0 := time.Now(), e.gate.pausedTotal()
	elapsed := func() time.Duration { return time.Since(start) - (e.gate.pausedTotal() - paused0) } // không tính lúc tạm dừng
	stable := time.Duration(s.StableMs) * time.Millisecond
	var base, prev vision.Gray
	var stableFrom time.Duration
	diff, shots := 0.0, 0
	_, err = e.waitLoop(ctx, s, "screen", func() (bool, int, error) {
		img, err := e.scr.Capture(area)
		if err != nil {
			return false, 0, err
		}
		g := vision.Thumb(img)
		if shots++; shots == 1 {
			base, prev, stableFrom = g, g, elapsed()
			return s.WaitMode == model.WaitStable && stable == 0, 0, nil
		}
		if s.WaitMode == model.WaitChange {
			diff = vision.DiffPercent(base, g)
			return diff > float64(s.DiffPct), 0, nil
		}
		diff = vision.DiffPercent(prev, g)
		prev = g
		if diff > float64(s.DiffPct) {
			stableFrom = elapsed()
		}
		still := elapsed() - stableFrom
		return still >= stable, int(min(100, still*100/max(stable, 1))), nil
	})
	diff = math.Round(diff*100) / 100
	if errors.Is(err, errTimeout) {
		return nil, &timeoutErr{msg: screenTimeoutMsg(s, diff), detail: map[string]any{"kind": "screen", "waitMode": s.WaitMode, "diff": diff, "diffPct": s.DiffPct, "area": area}}
	}
	if err != nil {
		return nil, err
	}
	return map[string]any{"ms": int(elapsed().Milliseconds()), "diff": diff}, nil
}

func screenTimeoutMsg(s model.Step, diff float64) string {
	if s.WaitMode == model.WaitChange {
		return fmt.Sprintf("màn hình không đổi sau %s (lệch %.2f %%, cần > %d %%)", fmtSeconds(s.TimeoutMs), diff, s.DiffPct)
	}
	return fmt.Sprintf("màn hình chưa đứng yên %s liền sau %s (lần so cuối lệch %.2f %%, cần ≤ %d %%)", fmtSeconds(s.StableMs), fmtSeconds(s.TimeoutMs), diff, s.DiffPct)
}
