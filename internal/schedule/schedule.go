// Package schedule – bộ hẹn giờ trong app (bước 22): mỗi nhịp xem hẹn giờ nào tới giờ (model.NextRun tính từ nhịp
// trước) rồi gọi Fire. Lúc app tắt thì lỡ giờ – không chạy bù (hộp thoại hẹn giờ ghi "app cần đang mở").
package schedule

import (
	"context"
	"time"

	"ghosthand/internal/model"
)

const Interval = 15 * time.Second

// Due – các hẹn giờ có lần chạy trong (last, now].
func Due(list []model.Schedule, last, now time.Time) []model.Schedule {
	var out []model.Schedule
	for _, s := range list {
		if t, ok := model.NextRun(s, last); ok && !t.After(now) {
			out = append(out, s)
		}
	}
	return out
}

// Loop chạy tới khi ctx huỷ: mỗi Interval đọc danh sách (load) và gọi fire cho hẹn giờ tới giờ.
func Loop(ctx context.Context, load func() []model.Schedule, fire func(model.Schedule)) {
	last := time.Now()
	tick := time.NewTicker(Interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-tick.C:
			for _, s := range Due(load(), last, now) {
				fire(s)
			}
			last = now
		}
	}
}
