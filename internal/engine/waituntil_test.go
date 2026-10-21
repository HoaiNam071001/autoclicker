package engine

import (
	"context"
	"os"
	"testing"
	"time"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

func untilStep(at time.Time, passed string) model.Step {
	return modeltest.Step(model.StepWaitUntil, func(s *model.Step) { s.At, s.IfPassed = at.Format("2006-01-02 15:04:05.000"), passed })
}

func TestWaitUntilOnTime(t *testing.T) {
	e := New(&devicetest.FakeInput{}, model.Script{}, Options{})
	var seen *WaitProgress
	e.onProgress = func(p Progress) {
		if p.Wait != nil {
			seen = p.Wait
		}
	}
	target := time.Now().Add(300 * time.Millisecond).Truncate(time.Millisecond) // giờ hẹn ghi tới ms
	out, err := e.waitUntil(t.Context(), untilStep(target, ""))
	if err != nil {
		t.Fatal(err)
	}
	if late := out["lateMs"].(int); late < 0 || late > 20 {
		t.Errorf("trễ %d ms so với giờ hẹn", late)
	}
	if d := time.Since(target); d < 0 {
		t.Errorf("chạy tiếp sớm hơn giờ hẹn %v", -d)
	}
	if seen == nil || seen.Kind != "time" || seen.TimeoutMs < 250 {
		t.Errorf("tiến độ chờ: %+v", seen)
	}
}

func TestWaitUntilPassed(t *testing.T) {
	e := New(&devicetest.FakeInput{}, model.Script{}, Options{})
	past := time.Now().Add(-time.Minute)
	if out, err := e.waitUntil(t.Context(), untilStep(past, "")); err != nil || out["lateMs"].(int) < 59000 {
		t.Errorf("đã qua, đi tiếp ngay: %v %v", out, err)
	}
	if _, err := e.waitUntil(t.Context(), untilStep(past, model.PassedFail)); errCode(err) != model.ErrTimeout {
		t.Errorf("đã qua, fail: %v (code %s)", err, errCode(err))
	}
	// tomorrow với giờ không ghi ngày: chờ ~24 h – Dừng giữa chừng phải ăn ngay
	s := modeltest.Step(model.StepWaitUntil, func(s *model.Step) { s.At, s.IfPassed = past.Format("15:04:05"), model.PassedTomorrow })
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	var total int64
	e.onProgress = func(p Progress) {
		if p.Wait != nil {
			total = p.Wait.TimeoutMs
		}
	}
	start := time.Now()
	if _, err := e.waitUntil(ctx, s); err == nil || time.Since(start) > 500*time.Millisecond {
		t.Errorf("Dừng giữa lúc chờ: %v sau %v", err, time.Since(start))
	}
	if total < int64(23*time.Hour/time.Millisecond) {
		t.Errorf("tomorrow phải chờ tới ngày mai, tổng chờ %d ms", total)
	}
}

func TestRealNTP(t *testing.T) {
	if os.Getenv("GHOSTHAND_REAL") != "1" {
		t.Skip("đặt GHOSTHAND_REAL=1 để hỏi giờ chuẩn thật")
	}
	off, err := ntpOffset(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("giờ chuẩn − giờ máy = %v", off)
	if off > time.Minute || off < -time.Minute {
		t.Errorf("lệch %v – nghi đọc sai gói NTP", off)
	}
}
