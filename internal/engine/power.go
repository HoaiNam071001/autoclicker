package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/power"
)

// Node Nguồn (bước 43 S3): khoá màn, tắt màn hình, ngủ, đăng xuất, tắt, khởi động lại. Ngủ / đăng xuất / tắt / khởi động
// lại làm mất phiên làm việc nên cần người dùng đồng ý trước (script.DangerPower) và luôn đếm ngược – bấm Dừng (thanh
// mini, F8, góc màn hình) trong lúc đếm là huỷ, máy không bị động tới.

func (e *Engine) powerStep(ctx context.Context, s model.Step) (map[string]any, error) {
	need := model.PowerNeedsConsent(s.PowerOp)
	if need && !e.allowDanger {
		return nil, dangerErr("thao tác nguồn máy")
	}
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	secs := s.CountdownSec
	if secs == 0 && need {
		secs = model.PowerCountdownSec
	}
	if secs > 0 {
		if err := e.countdown(ctx, secs); err != nil {
			return nil, err
		}
	}
	err := e.powers.Do(ctx, s.PowerOp)
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case errors.Is(err, power.ErrUnsupported):
		return nil, &invalidErr{err.Error()}
	case err != nil:
		return nil, err
	}
	return map[string]any{"action": s.PowerOp, "countdownSec": secs}, nil
}

// countdown đếm ngược secs giây, báo tiến độ cho thanh mini; Dừng giữa chừng trả lỗi huỷ.
func (e *Engine) countdown(ctx context.Context, secs int) error {
	defer e.update(true, func(p *Progress) { p.Wait = nil })
	start := time.Now()
	total := time.Duration(secs) * time.Second
	for {
		if err := e.gate.wait(ctx); err != nil {
			return err
		}
		elapsed := time.Since(start)
		e.update(true, func(p *Progress) {
			p.Wait = &WaitProgress{Kind: "power", ElapsedMs: elapsed.Milliseconds(), TimeoutMs: total.Milliseconds()}
		})
		if elapsed >= total {
			return nil
		}
		if err := Sleep(ctx, int(min(total-elapsed, 250*time.Millisecond).Milliseconds())); err != nil {
			return fmt.Errorf("đã huỷ trước khi thực hiện: %w", err)
		}
	}
}
