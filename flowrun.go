package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Engine chạy flow (bước 20): từ node start, chạy node → StepResult → đi theo cạnh của cổng ra → node kế.
// Hết cạnh: cổng fail = dừng kịch bản báo lỗi, cổng khác = xong vòng. Lặp Loops vòng như trước.

// maxNodeRuns – số lần chạy node tối đa mỗi vòng (chặn vòng lặp không có node Đếm).
const maxNodeRuns = 10000

// notFoundErr – chờ ảnh / chữ / màu hết giờ không thấy (StepError code not_found).
type notFoundErr struct{ msg string }

func (e *notFoundErr) Error() string { return e.msg }

// run chạy tới khi đủ số vòng, hết MaxMinutes hoặc ctx bị huỷ.
// Hết thời gian là kết thúc bình thường (nil); bị huỷ trả ctx.Err().
// MaxMinutes tính theo đồng hồ thật, kể cả lúc tạm dừng.
func (e *engine) run(ctx context.Context) error {
	sc := e.sc
	if err := validateScript(sc, e.screen); err != nil {
		return err
	}
	runCtx := ctx
	if sc.MaxMinutes > 0 {
		var cancel context.CancelFunc
		runCtx, cancel = context.WithTimeout(ctx, time.Duration(sc.MaxMinutes)*time.Minute)
		defer cancel()
	}

	e.byID = make(map[string]Step, len(sc.Steps))
	e.order = map[string]int{}
	total := 0
	for _, id := range flowOrder(sc) {
		for _, s := range sc.Steps {
			if s.ID == id && s.Enabled {
				total++
				e.order[id] = total
			}
		}
	}
	for _, s := range sc.Steps {
		e.byID[s.ID] = s
	}
	e.mu.Lock()
	e.start = time.Now()
	e.prog = Progress{Loops: sc.Loops, Total: total}
	e.mu.Unlock()

	for loop := 1; sc.Loops == 0 || loop <= sc.Loops; loop++ {
		err := e.runLoop(runCtx, loop)
		if err == nil {
			continue
		}
		if ctx.Err() == nil && errors.Is(err, context.DeadlineExceeded) {
			return nil // hết MaxMinutes
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}

// runLoop chạy một vòng flow. Kết quả / biến đặt lại mỗi vòng (Đếm cũng đếm lại từ 0).
func (e *engine) runLoop(ctx context.Context, loop int) error {
	e.rv = runVars{byRef: map[string]*StepResult{}, vars: map[string]any{}}
	e.results = map[string]*StepResult{}
	var trail []string // các node gần nhất – báo lỗi lặp vô hạn
	cur := e.sc.Start
	for n := 1; cur != ""; n++ {
		s := e.byID[cur]
		if n > maxNodeRuns {
			e.failed = s.ID
			return fmt.Errorf("vòng %d: đã chạy %d node mà chưa xong – có thể lặp vô hạn ở %s; thêm node Đếm để giới hạn số vòng lặp",
				loop, maxNodeRuns, strings.Join(trail, " → "))
		}
		res := e.execNode(ctx, s, loop)
		e.results[s.ID], e.rv.byRef[s.Ref], e.rv.last = &res, &res, &res
		if e.onNode != nil {
			e.onNode(res)
		}
		if res.Status == StatusCancelled {
			return ctx.Err()
		}
		if trail = append(trail, s.Ref); len(trail) > 6 {
			trail = trail[1:]
		}
		next := edgeOut(e.sc.Edges, s.ID, res.Port)
		if res.Status == StatusFail && next == "" {
			e.failed = s.ID
			return fmt.Errorf("vòng %d, %s: %s", loop, s.Ref, res.Error.Message)
		}
		if res.Status == StatusSuccess {
			if err := sleep(ctx, e.jitterMs(s.DelayAfter, s.JitterMs)); err != nil {
				return err
			}
		}
		cur = next
	}
	return nil
}

// execNode chạy một node, thử lại theo Retry; không bao giờ trả lỗi – mọi kết quả nằm trong StepResult.
func (e *engine) execNode(ctx context.Context, s Step, loop int) StepResult {
	res := StepResult{NodeID: s.ID, Ref: s.Ref, StartedAt: time.Now()}
	defer func() { res.DurationMs = time.Since(res.StartedAt).Milliseconds() }()
	logf := func(level, format string, a ...any) {
		res.Logs = append(res.Logs, LogLine{AtMs: time.Since(res.StartedAt).Milliseconds(), Level: level, Msg: fmt.Sprintf(format, a...)})
	}
	if !s.Enabled {
		res.Status, res.Port = StatusSkipped, PortSuccess
		logf("info", "node đang tắt – đi tiếp cổng success")
		return res
	}
	e.miss = nil
	tries := 1
	delay := 0
	if s.Retry != nil {
		tries += s.Retry.Count
		delay = s.Retry.DelayMs
	}
	var err error
	for t := 1; t <= tries; t++ {
		if t > 1 {
			if err = sleep(ctx, delay); err != nil {
				break
			}
		}
		if err = e.gate.wait(ctx); err != nil {
			break
		}
		res.Attempts = t
		e.update(true, func(p *Progress) {
			p.Loop, p.Step, p.NodeID, p.Ref, p.Retry, p.Retries = loop, e.order[s.ID], s.ID, s.Ref, t-1, tries-1
		})
		var rs Step
		if rs, err = resolveNodeVars(s, e.rv); err != nil {
			res.Status, res.Port = StatusFail, PortFail
			res.Error = &StepError{Code: ErrInvalidInput, Message: err.Error()}
			logf("error", "%v", err)
			return res
		}
		var out map[string]any
		if out, err = e.runStep(ctx, rs); err == nil {
			res.Status, res.Port, res.Output = StatusSuccess, PortSuccess, out
			if tries > 1 {
				e.update(true, func(p *Progress) { p.Retry, p.Retries = 0, 0 })
			}
			return res
		}
		if ctx.Err() != nil {
			break
		}
		if t < tries {
			logf("warn", "lần %d / %d: %v – thử lại", t, tries, err)
		}
	}
	if tries > 1 {
		e.update(true, func(p *Progress) { p.Retry, p.Retries = 0, 0 })
	}
	if ctx.Err() != nil {
		res.Status = StatusCancelled
		res.Error = &StepError{Code: ErrCancelled, Message: "đã dừng"}
		return res
	}
	res.Status, res.Port = StatusFail, PortFail
	res.Error = &StepError{Code: errCode(err), Message: err.Error()}
	logf("error", "%v", err)
	return res
}

func errCode(err error) string {
	var nf *notFoundErr
	switch {
	case errors.As(err, &nf):
		return ErrNotFound
	case errors.Is(err, errNoPermission):
		return ErrPermission
	}
	return ErrInternal
}
