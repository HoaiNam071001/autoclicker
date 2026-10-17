package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// Node điều khiển của flow (bước 20 F3): Nếu, Switch theo giá trị, Switch theo màn hình, Đếm, Kết thúc, Đặt biến.
// Không thao tác chuột / phím – chỉ chọn cổng ra. Điều kiện xét ở cond.go.

// Cổng ra của node điều khiển (ngoài success / fail ở flow.go).
const (
	PortTrue    = "true"
	PortFalse   = "false"
	PortDefault = "default" // switch: không case nào khớp
	PortNone    = "none"    // switchScreen: hết giờ không case nào hiện
	PortNext    = "next"    // counter: còn lượt
	PortDone    = "done"    // counter: hết lượt
	casePrefix  = "case:"
)

// Step.EndStatus / Step.VarOp
const (
	EndSuccess = "success"
	EndFail    = "fail"
	VarSet     = "set"
	VarAdd     = "add"
)

var caseOps = []string{"eq", "contains", "regex"}

var (
	caseKeyRe = regexp.MustCompile(`^[A-Za-z0-9_]{1,32}$`)
	varNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,31}$`)
)

// casePorts – cổng case:<key> theo thứ tự case, rồi cổng cuối (default / none).
func casePorts(cases []Case, last string) []string {
	out := make([]string, 0, len(cases)+1)
	for _, c := range cases {
		out = append(out, casePrefix+c.Key)
	}
	return append(out, last)
}

// endErr – node Kết thúc báo thất bại (StepError code "end"): dừng kịch bản với lời nhắn của người soạn.
type endErr struct{ msg string }

func (e *endErr) Error() string { return e.msg }

// runControl chạy node điều khiển; trả cổng ra + output. logf ghi chi tiết (từng điều kiện) vào nhật ký node.
func (e *engine) runControl(ctx context.Context, s Step, logf func(level, format string, a ...any)) (string, map[string]any, error) {
	switch s.Type {
	case StepIf:
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
			return PortTrue, out, nil
		}
		return PortFalse, out, nil

	case StepSwitchScreen:
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
			return PortNone, map[string]any{"case": PortNone}, nil
		}
		c := s.Cases[won]
		out := firstHit(parts)
		out["case"] = c.Key
		logf("info", "case %s hiện trước", caseName(c))
		return casePrefix + c.Key, out, nil

	case StepSwitch:
		for _, c := range s.Cases {
			ok, err := matchCase(s.Value, c)
			if err != nil {
				return "", nil, &invalidErr{err.Error()}
			}
			if ok {
				logf("info", "“%s” khớp case %s", s.Value, caseName(c))
				return casePrefix + c.Key, map[string]any{"value": s.Value, "case": c.Key}, nil
			}
		}
		logf("info", "“%s” không khớp case nào", s.Value)
		return PortDefault, map[string]any{"value": s.Value, "case": PortDefault}, nil

	case StepCounter:
		e.counts[s.ID]++
		n := e.counts[s.ID]
		if n <= s.Max {
			return PortNext, map[string]any{"count": n}, nil
		}
		e.counts[s.ID] = 0 // hết lượt thì đếm lại – đi qua vòng thử lần nữa (vd vòng ngoài) vẫn có đủ lượt
		logf("info", "đã đi qua %d lần (tối đa %d) – hết lượt", n, s.Max)
		return PortDone, map[string]any{"count": n}, nil

	case StepEnd:
		if s.EndStatus == EndFail {
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

	case StepSetVar:
		var v any = s.Value
		if s.VarOp == VarAdd {
			cur := 0.0
			if old, ok := e.rv.vars[s.VarName]; ok {
				f, err := parseNumber(fmtValue(old))
				if err != nil {
					return "", nil, &invalidErr{fmt.Sprintf("biến %s đang là “%v”, không cộng được", s.VarName, old)}
				}
				cur = f
			}
			add, err := parseNumber(s.Value)
			if err != nil {
				return "", nil, &invalidErr{fmt.Sprintf("cộng cần một số, có “%s”", s.Value)}
			}
			v = numValue(cur + add)
		}
		e.rv.vars[s.VarName] = v
		return PortSuccess, map[string]any{"name": s.VarName, "value": v}, nil
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
func (e *engine) watch(ctx context.Context, s Step, check func(*snapshot) (bool, error)) (bool, error) {
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

func caseName(c Case) string {
	if c.Label != "" {
		return c.Label
	}
	return c.Key
}

// checkControl kiểm cấu trúc node điều khiển khi mở file (giá trị lạ → lỗi; thiếu dữ liệu để lúc chạy báo).
func checkControl(s *Step) error {
	if s.Cond != nil {
		if err := checkCond(s.Cond, 0); err != nil {
			return err
		}
	}
	keys := map[string]bool{}
	for _, c := range s.Cases {
		if !caseKeyRe.MatchString(c.Key) {
			return fmt.Errorf("mã case %q không hợp lệ – chữ / số / _", c.Key)
		}
		if keys[c.Key] {
			return fmt.Errorf("mã case %q bị trùng", c.Key)
		}
		keys[c.Key] = true
		if c.Op != "" && !slices.Contains(caseOps, c.Op) {
			return fmt.Errorf("cách so của case không hợp lệ: %q", c.Op)
		}
		if c.Cond != nil {
			if err := checkCond(c.Cond, 0); err != nil {
				return err
			}
		}
	}
	if s.EndStatus != "" && s.EndStatus != EndSuccess && s.EndStatus != EndFail {
		return fmt.Errorf("kết quả kết thúc không hợp lệ: %q", s.EndStatus)
	}
	if s.VarOp != "" && s.VarOp != VarSet && s.VarOp != VarAdd {
		return fmt.Errorf("cách đặt biến không hợp lệ: %q", s.VarOp)
	}
	if s.Max < 0 {
		return errors.New("số lần tối đa không âm")
	}
	return nil
}

// validateControl – node điều khiển đủ dữ liệu để chạy (validateScript gọi). inside: điểm nằm trong vùng.
func validateControl(i int, s Step, inside func(Point) bool) error {
	switch s.Type {
	case StepIf:
		if s.Cond == nil {
			return stepErr(i, s, "cond", "chưa có điều kiện")
		}
		if err := validateCond(*s.Cond, inside); err != nil {
			return stepErr(i, s, "cond", "%v", err)
		}
		if s.TimeoutMs < 0 || s.TimeoutMs > 0 && s.IntervalMs <= 0 {
			return stepErr(i, s, "timeoutMs", "thời gian chờ không âm, kiểm tra mỗi phải > 0")
		}
	case StepSwitch, StepSwitchScreen:
		if len(s.Cases) == 0 {
			return stepErr(i, s, "cases", "chưa có case nào")
		}
		for j, c := range s.Cases {
			if s.Type == StepSwitch {
				if !slices.Contains(caseOps, c.Op) {
					return stepErr(i, s, "cases", "case %s: cách so không hợp lệ %q", caseName(c), c.Op)
				}
				if c.Op == "regex" {
					if _, err := regexp.Compile(c.Match); err != nil {
						return stepErr(i, s, "cases", "case %s: mẫu sai – %v", caseName(c), err)
					}
				}
				continue
			}
			if c.Cond == nil || !isLeafScreen(c.Cond.Kind) && c.Cond.Kind != CondAll && c.Cond.Kind != CondAny {
				return stepErr(i, s, "cases", "case %d (%s): cần điều kiện thấy chữ / ảnh / màu", j+1, caseName(c))
			}
			if err := validateCond(*c.Cond, inside); err != nil {
				return stepErr(i, s, "cases", "case %s: %v", caseName(c), err)
			}
		}
		if s.Type == StepSwitchScreen && (s.TimeoutMs < 0 || s.IntervalMs <= 0) {
			return stepErr(i, s, "timeoutMs", "thời gian chờ không âm, kiểm tra mỗi phải > 0")
		}
	case StepCounter:
		if s.Max < 1 {
			return stepErr(i, s, "max", "số lần tối đa phải ≥ 1")
		}
	case StepEnd:
		if s.EndStatus != EndSuccess && s.EndStatus != EndFail {
			return stepErr(i, s, "endStatus", "chọn kết thúc thành công hay thất bại")
		}
	case StepSetVar:
		if !varNameRe.MatchString(s.VarName) {
			return stepErr(i, s, "varName", "tên biến %q không hợp lệ – chữ cái / _ đầu, sau đó chữ / số / _", s.VarName)
		}
		if s.VarOp == VarAdd {
			if _, err := parseNumber(s.Value); err != nil && !nodeVarRe.MatchString(s.Value) && !varRe.MatchString(s.Value) {
				return stepErr(i, s, "value", "cộng cần một số, có “%s”", s.Value)
			}
		} else if s.VarOp != VarSet {
			return stepErr(i, s, "varOp", "chọn gán hay cộng")
		}
	}
	return nil
}
