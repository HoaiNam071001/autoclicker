package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ghosthand/internal/model"
	"ghosthand/internal/script"
)

// Node điều khiển của flow (bước 20 F3): Nếu, Switch theo giá trị, Switch theo màn hình, Đếm, Kết thúc, Đặt biến.
// Không thao tác chuột / phím – chỉ chọn cổng ra. Điều kiện xét ở cond.go.

// endErr – node Kết thúc báo thất bại (StepError code "end"): dừng kịch bản với lời nhắn của người soạn.
type endErr struct{ msg string }

func (e *endErr) Error() string { return e.msg }

// runControl chạy node điều khiển; trả cổng ra + output. logf ghi chi tiết (từng điều kiện) vào nhật ký node.
func (e *Engine) runControl(ctx context.Context, s model.Step, logf func(level, format string, a ...any)) (string, map[string]any, error) {
	switch s.Type {
	case model.StepIf:
		var parts []CondPart
		ok, err := e.watch(ctx, s, func(snap *snapshot) (bool, error) {
			parts = parts[:0]
			v, _, err := e.evalCond(*s.Cond, s.ID+"#cond", nil, snap, &parts)
			return v, err
		})
		if err != nil {
			return "", nil, err
		}
		logParts(parts, logf)
		out := firstHit(parts)
		out["result"] = ok
		if ok {
			return model.PortTrue, out, nil
		}
		return model.PortFalse, out, nil

	case model.StepSwitchScreen:
		var parts []CondPart
		won := -1
		_, err := e.watch(ctx, s, func(snap *snapshot) (bool, error) {
			parts = parts[:0]
			for i, c := range s.Cases {
				p0 := len(parts)
				v, _, err := e.evalCond(*c.Cond, fmt.Sprintf("%s#case%d", s.ID, i), []int{i}, snap, &parts)
				if err != nil {
					return false, err
				}
				for j := p0; j < len(parts); j++ {
					parts[j].Case = c.Key
				}
				if v {
					won = i
					return true, nil
				}
			}
			return false, nil
		})
		if err != nil {
			return "", nil, err
		}
		logParts(parts, logf)
		if won < 0 {
			return model.PortNone, map[string]any{"case": model.PortNone}, nil
		}
		c := s.Cases[won]
		out := firstHit(parts)
		out["case"] = c.Key
		logf("info", "case %s hiện trước", script.CaseName(c))
		return model.CasePrefix + c.Key, out, nil

	case model.StepSwitch:
		for _, c := range s.Cases {
			ok, err := matchCase(s.Value, c)
			if err != nil {
				return "", nil, &invalidErr{err.Error()}
			}
			if ok {
				logf("info", "“%s” khớp case %s", s.Value, script.CaseName(c))
				return model.CasePrefix + c.Key, map[string]any{"value": s.Value, "case": c.Key}, nil
			}
		}
		logf("info", "“%s” không khớp case nào", s.Value)
		return model.PortDefault, map[string]any{"value": s.Value, "case": model.PortDefault}, nil

	case model.StepCounter:
		e.counts[s.ID]++
		n := e.counts[s.ID]
		if n <= s.Max {
			return model.PortNext, map[string]any{"count": n}, nil
		}
		e.counts[s.ID] = 0 // hết lượt thì đếm lại – đi qua vòng thử lần nữa (vd vòng ngoài) vẫn có đủ lượt
		logf("info", "đã đi qua %d lần (tối đa %d) – hết lượt", n, s.Max)
		return model.PortDone, map[string]any{"count": n}, nil

	case model.StepEnd:
		if s.EndStatus == model.EndFail {
			msg := s.Message
			if strings.TrimSpace(msg) == "" {
				msg = "kết thúc: thất bại"
			}
			return "", nil, &endErr{msg}
		}
		if s.Message == "" {
			return "", nil, nil // không có cổng ra → xong vòng
		}
		logf("info", "%s", s.Message)
		return "", map[string]any{"message": s.Message}, nil // lời nhắn hiện ở thanh trạng thái / nhật ký (F5)

	case model.StepSetVar:
		var v any = s.Value
		if s.VarOp == model.VarAdd {
			cur := 0.0
			if old, ok := e.rv.vars[s.VarName]; ok {
				f, err := script.ParseNumber(fmtValue(old))
				if err != nil {
					return "", nil, &invalidErr{fmt.Sprintf("biến %s đang là “%v”, không cộng được", s.VarName, old)}
				}
				cur = f
			}
			add, err := script.ParseNumber(s.Value)
			if err != nil {
				return "", nil, &invalidErr{fmt.Sprintf("cộng cần một số, có “%s”", s.Value)}
			}
			v = numValue(cur + add)
		}
		e.rv.vars[s.VarName] = v
		return model.PortSuccess, map[string]any{"name": s.VarName, "value": v}, nil
	}
	return "", nil, fmt.Errorf("không phải node điều khiển: %q", s.Type)
}

// invalidErr – dữ liệu node sai lúc chạy (StepError code invalid_input).
type invalidErr struct{ msg string }

func (e *invalidErr) Error() string { return e.msg }

// numValue – số nguyên giữ dạng int (in ra "3", không phải "3.0").
func numValue(f float64) any {
	if f == float64(int(f)) {
		return int(f)
	}
	return f
}

// watch xét check trên một ảnh chụp mới mỗi IntervalMs tới khi đúng, tối đa TimeoutMs (0 = xem một lần).
// Hết giờ không phải lỗi – trả false.
func (e *Engine) watch(ctx context.Context, s model.Step, check func(*snapshot) (bool, error)) (bool, error) {
	if s.TimeoutMs <= 0 {
		return check(newSnapshot(e.scr, e.ocr))
	}
	found := false
	_, err := e.waitLoop(ctx, s, "cond", func() (bool, int, error) {
		ok, err := check(newSnapshot(e.scr, e.ocr))
		found = ok
		return ok, 0, err
	})
	if errors.Is(err, errTimeout) {
		return false, nil
	}
	return found, err
}

func logParts(parts []CondPart, logf func(level, format string, a ...any)) {
	for _, p := range parts {
		mark := map[bool]string{true: "✓", false: "✗"}[p.Ok]
		logf("info", "%s %s", mark, p.Detail)
	}
}
