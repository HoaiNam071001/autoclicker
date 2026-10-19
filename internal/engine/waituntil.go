package engine

import (
	"context"
	"fmt"
	"time"

	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Node Chờ đến giờ (bước 35 N6): chờ tới giờ At (giờ máy, hoặc giờ chuẩn NTP). Còn xa thì nghỉ từng đoạn ≤ 1 s
// (Dừng / tạm dừng ăn ngay, báo tiến độ, đồng hồ máy chỉnh giữa chừng vẫn đúng), sát giờ (< 20 ms) chờ vòng ngắn.

func (e *Engine) waitUntil(ctx context.Context, s model.Step) (map[string]any, error) {
	var offset time.Duration
	if s.UseNTP {
		if off, err := ntpOffset(ctx); err == nil {
			offset = off
		} // không lấy được giờ chuẩn: dùng giờ máy (output ntpOffsetMs 0)
	}
	now := func() time.Time { return time.Now().Add(offset) }
	target, dated, err := script.ParseAt(s.At, now())
	if err != nil {
		return nil, &invalidErr{fmt.Sprintf("giờ “%s”: %v", s.At, err)}
	}
	if !target.After(now()) {
		switch {
		case s.IfPassed == model.PassedFail:
			return nil, &timeoutErr{fmt.Sprintf("đã qua giờ %s (trễ %s)", target.Format("15:04:05.000"), fmtDur(now().Sub(target)))}
		case s.IfPassed == model.PassedTomorrow && !dated:
			target = target.AddDate(0, 0, 1)
		}
	}
	start := now()
	total := max(target.Sub(start), 0)
	defer e.update(true, func(p *Progress) { p.Wait = nil })
	for first := true; ; first = false {
		if err := e.gate.wait(ctx); err != nil {
			return nil, err
		}
		rem := target.Sub(now())
		e.update(first, func(p *Progress) {
			p.Wait = &WaitProgress{Kind: "time", ElapsedMs: now().Sub(start).Milliseconds(), TimeoutMs: total.Milliseconds()}
		})
		if rem <= 0 {
			break
		}
		if rem > 20*time.Millisecond {
			if err := Sleep(ctx, int(min(rem-15*time.Millisecond, time.Second).Milliseconds())); err != nil {
				return nil, err
			}
			continue
		}
		for now().Before(target) { // sát giờ: vòng ngắn, lệch < 1 ms
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			time.Sleep(200 * time.Microsecond)
		}
	}
	fired := now()
	return map[string]any{
		"firedAt": fired.Format("2006-01-02T15:04:05.000Z07:00"), "target": target.Format("2006-01-02T15:04:05.000Z07:00"),
		"lateMs": int(fired.Sub(target).Milliseconds()), "ntpOffsetMs": int(offset.Milliseconds()),
	}, nil
}

// fmtDur – "2 phút 5 giây" / "350 ms" cho lời báo.
func fmtDur(d time.Duration) string {
	switch {
	case d < time.Second:
		return fmt.Sprintf("%d ms", d.Milliseconds())
	case d < time.Minute:
		return fmt.Sprintf("%.1f giây", d.Seconds())
	case d < time.Hour:
		return fmt.Sprintf("%d phút %d giây", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%d giờ %d phút", int(d.Hours()), int(d.Minutes())%60)
}
