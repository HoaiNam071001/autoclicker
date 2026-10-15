package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// dragStepMs – khoảng thời gian giữa các lần di chuột khi kéo thả.
const dragStepMs = 10

// validateScript kiểm tra kịch bản trước khi chạy để không dừng giữa chừng vì dữ liệu sai.
func validateScript(sc Script) error {
	enabled := 0
	for i, s := range sc.Steps {
		if !s.Enabled {
			continue
		}
		enabled++
		need := 0
		switch s.Type {
		case StepClick, StepDouble, StepHold:
			need = 1
		case StepDrag:
			need = 2
		case StepWait:
		case StepText:
			if s.Text == "" {
				return fmt.Errorf("bước %d: chưa nhập nội dung gõ", i+1)
			}
		case StepKey:
			if s.Key == "" {
				return fmt.Errorf("bước %d: chưa chọn phím", i+1)
			}
		default:
			return fmt.Errorf("bước %d: loại bước không hợp lệ %q", i+1, s.Type)
		}
		if len(s.Points) < need {
			return fmt.Errorf("bước %d: cần %d điểm, mới có %d", i+1, need, len(s.Points))
		}
	}
	if enabled == 0 {
		return errors.New("kịch bản không có bước nào được bật")
	}
	return nil
}

// runScript chạy kịch bản tới khi đủ số vòng, hết MaxMinutes hoặc ctx bị huỷ.
// Hết thời gian là kết thúc bình thường (nil); bị huỷ trả ctx.Err().
func runScript(ctx context.Context, in Input, sc Script) error {
	if err := validateScript(sc); err != nil {
		return err
	}
	runCtx := ctx
	if sc.MaxMinutes > 0 {
		var cancel context.CancelFunc
		runCtx, cancel = context.WithTimeout(ctx, time.Duration(sc.MaxMinutes)*time.Minute)
		defer cancel()
	}

	for loop := 1; sc.Loops == 0 || loop <= sc.Loops; loop++ {
		for i, s := range sc.Steps {
			if !s.Enabled {
				continue
			}
			err := runStep(runCtx, in, sc.Region, s)
			if err == nil {
				err = sleep(runCtx, s.DelayAfter)
			}
			if err != nil {
				if ctx.Err() == nil && errors.Is(err, context.DeadlineExceeded) {
					return nil // hết MaxMinutes
				}
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return fmt.Errorf("vòng %d, bước %d: %w", loop, i+1, err)
			}
		}
	}
	return nil
}

func runStep(ctx context.Context, in Input, region *Rect, s Step) error {
	pt := func(i int) (int, int) {
		p := s.Points[i]
		if region != nil {
			return region.X + p.X, region.Y + p.Y
		}
		return p.X, p.Y
	}

	switch s.Type {
	case StepClick:
		x, y := pt(0)
		n := max(s.Clicks, 1)
		for i := 0; i < n; i++ {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := in.Move(x, y); err != nil {
				return err
			}
			if err := in.Click(s.Button, false); err != nil {
				return err
			}
			if i < n-1 {
				if err := sleep(ctx, s.ClickInterval); err != nil {
					return err
				}
			}
		}
	case StepDouble:
		x, y := pt(0)
		if err := in.Move(x, y); err != nil {
			return err
		}
		return in.Click(s.Button, true)
	case StepHold:
		x, y := pt(0)
		if err := in.Move(x, y); err != nil {
			return err
		}
		if err := in.Toggle(s.Button, true); err != nil {
			return err
		}
		defer in.Toggle(s.Button, false) // luôn nhả chuột, kể cả khi bị huỷ giữa chừng
		return sleep(ctx, s.HoldMs)
	case StepDrag:
		return drag(ctx, in, s, pt)
	case StepWait:
		// chỉ chờ DelayAfter
	case StepText:
		// Gõ từng ký tự để dừng được giữa chừng với chuỗi dài.
		for _, r := range s.Text {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := in.TypeStr(string(r)); err != nil {
				return err
			}
		}
	case StepKey:
		return in.KeyTap(s.Key, s.Modifiers)
	}
	return nil
}

// drag giữ nút tại A, di dần tới B trong DragMs, rồi nhả.
func drag(ctx context.Context, in Input, s Step, pt func(int) (int, int)) error {
	ax, ay := pt(0)
	bx, by := pt(1)
	if err := in.Move(ax, ay); err != nil {
		return err
	}
	if err := in.Toggle(s.Button, true); err != nil {
		return err
	}
	defer in.Toggle(s.Button, false)

	n := max(s.DragMs/dragStepMs, 1)
	for i := 1; i <= n; i++ {
		if err := sleep(ctx, s.DragMs/n); err != nil {
			return err
		}
		x := ax + (bx-ax)*i/n
		y := ay + (by-ay)*i/n
		if err := in.Drag(x, y, s.Button); err != nil {
			return err
		}
	}
	return nil
}

// sleep chờ ms mili giây hoặc tới khi ctx bị huỷ.
func sleep(ctx context.Context, ms int) error {
	if ms <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(time.Duration(ms) * time.Millisecond)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
