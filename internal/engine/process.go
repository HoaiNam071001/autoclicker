package engine

import (
	"context"
	"fmt"
	"slices"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/procinfo"
)

// Node Tiến trình (bước 43 S5): đọc / chờ tiến trình khớp ProcName (model.MatchName trên tên tiến trình) qua
// procinfo.Backend – thấy cả tiến trình nền không cửa sổ. Không thoát tiến trình (việc của node Ứng dụng).

const procItemsMax = 200

func (e *Engine) process(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	switch s.ProcOp {
	case model.ProcWaitStart:
		var hit []procinfo.Info
		err := e.pollUntil(ctx, s, "process", func() (bool, error) {
			var err error
			hit, err = e.findProcs(ctx, s)
			return len(hit) > 0, err
		}, func() error {
			return &timeoutErr{msg: fmt.Sprintf("tiến trình %s chưa chạy sau %s", s.ProcName, fmtSeconds(s.TimeoutMs))}
		})
		return procOutput(hit), err
	case model.ProcWaitQuit:
		var hit []procinfo.Info
		err := e.pollUntil(ctx, s, "process", func() (bool, error) {
			var err error
			hit, err = e.findProcs(ctx, s)
			return len(hit) == 0, err
		}, func() error {
			return &timeoutErr{msg: fmt.Sprintf("tiến trình %s vẫn chạy sau %s", s.ProcName, fmtSeconds(s.TimeoutMs))}
		})
		return procOutput(hit), err
	case model.ProcWaitIdle:
		return e.waitProcIdle(ctx, s)
	}
	hit, err := e.findProcs(ctx, s)
	return procOutput(hit), err
}

// waitProcIdle – tổng CPU dưới ngưỡng liên tục StableMs; chạm ngưỡng thì đếm lại từ đầu.
func (e *Engine) waitProcIdle(ctx context.Context, s model.Step) (map[string]any, error) {
	below := float64(s.CPUBelow)
	if below <= 0 {
		below = model.ProcIdleCPU
	}
	var hit []procinfo.Info
	var quietSince time.Time
	var cpu float64
	err := e.pollUntil(ctx, s, "process", func() (bool, error) {
		var err error
		if hit, err = e.findProcs(ctx, s); err != nil {
			return false, err
		}
		if len(hit) == 0 {
			return false, &notFoundErr{msg: fmt.Sprintf("tiến trình %s không chạy", s.ProcName)}
		}
		cpu = procOutput(hit)["cpuPercent"].(float64)
		if cpu >= below {
			quietSince = time.Time{}
			return false, nil
		}
		if quietSince.IsZero() {
			quietSince = time.Now()
		}
		return time.Since(quietSince) >= time.Duration(s.StableMs)*time.Millisecond, nil
	}, func() error {
		return &timeoutErr{msg: fmt.Sprintf("tiến trình %s vẫn dùng CPU %.0f%% (cần dưới %.0f%%) sau %s", s.ProcName, cpu, below, fmtSeconds(s.TimeoutMs))}
	})
	return procOutput(hit), err
}

func (e *Engine) findProcs(ctx context.Context, s model.Step) ([]procinfo.Info, error) {
	all, err := e.procs.List(ctx)
	if err != nil {
		return nil, err
	}
	out := all[:0:0]
	for _, p := range all {
		ok, err := model.MatchName(p.Name, s.ProcName, s.ProcMatch)
		if err != nil {
			return nil, &invalidErr{err.Error()}
		}
		if ok {
			out = append(out, p)
		}
	}
	return out, nil
}

// procOutput – gộp các tiến trình khớp: pid / name của tiến trình nặng nhất, tổng CPU / RAM, items nặng → nhẹ.
func procOutput(hit []procinfo.Info) map[string]any {
	hit = slices.Clone(hit)
	slices.SortStableFunc(hit, func(a, b procinfo.Info) int {
		switch {
		case a.CPU > b.CPU:
			return -1
		case a.CPU < b.CPU:
			return 1
		}
		return 0
	})
	out := map[string]any{"running": len(hit) > 0, "count": len(hit), "pid": 0, "name": "", "cpuPercent": 0.0, "memMB": 0.0}
	var cpu, mem float64
	items := make([]any, 0, min(len(hit), procItemsMax))
	for i, p := range hit {
		cpu += p.CPU
		mem += p.MemMB
		if i == 0 {
			out["pid"], out["name"] = p.PID, p.Name
		}
		if i < procItemsMax {
			items = append(items, fmt.Sprintf("%s · %d · CPU %.1f%% · RAM %.0f MB", p.Name, p.PID, p.CPU, p.MemMB))
		}
	}
	out["cpuPercent"], out["memMB"], out["items"] = round1(cpu), round1(mem), items
	return out
}

func round1(f float64) float64 { return float64(int(f*10+0.5)) / 10 }
