package script

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

var (
	caseKeyRe = regexp.MustCompile(`^[A-Za-z0-9_]{1,32}$`)
	varNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,31}$`)
)

var condOps = []string{"eq", "ne", "contains", "gt", "lt", "empty"}

func isLeafScreen(kind string) bool {
	return kind == model.CondImage || kind == model.CondText || kind == model.CondElement || kind == model.CondPixel
}

// ParseNumber đọc số, nhận cả dấu phẩy thập phân ("1,5").
func ParseNumber(s string) (float64, error) {
	return strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(s), ",", "."), 64)
}

// checkCond kiểm cấu trúc điều kiện (mở file): loại, phép so, lồng nhóm.
func checkCond(c *model.Cond, depth int) error {
	switch c.Kind {
	case model.CondImage:
		if err := checkMatch(c.MatchMode, c.ColorTol); err != nil {
			return err
		}
	case model.CondText, model.CondElement, model.CondPixel:
	case model.CondValue:
		if !slices.Contains(condOps, c.Op) {
			return fmt.Errorf("phép so không hợp lệ: %q", c.Op)
		}
	case model.CondAll, model.CondAny:
		if depth > 8 {
			return errors.New("điều kiện lồng quá sâu")
		}
		for i := range c.Conds {
			if err := checkCond(&c.Conds[i], depth+1); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("loại điều kiện không hợp lệ: %q", c.Kind)
	}
	for _, p := range c.Points {
		if !validAnchor(p.Ax) || !validAnchor(p.Ay) {
			return fmt.Errorf("neo điểm không hợp lệ: %q / %q", p.Ax, p.Ay)
		}
	}
	return nil
}

// validateCond – điều kiện đủ dữ liệu để chạy (trước khi chạy): chữ, ảnh mẫu, điểm + màu, nhóm không rỗng.
// inside: điểm nằm trong vùng.
func validateCond(c model.Cond, inside func(model.Point) bool) error {
	switch c.Kind {
	case model.CondText:
		if strings.TrimSpace(c.FindText) == "" {
			return errors.New("điều kiện chưa nhập chữ cần thấy")
		}
	case model.CondElement:
		if err := checkElementSelectors(c.Selectors); err != nil {
			return fmt.Errorf("điều kiện phần tử: %w", err)
		}
	case model.CondImage:
		if c.Template == "" {
			return errors.New("điều kiện chưa có ảnh mẫu")
		}
		if _, err := vision.DecodePNG(c.Template); err != nil {
			return fmt.Errorf("ảnh mẫu của điều kiện hỏng: %v", err)
		}
		if c.Threshold < 1 || c.Threshold > 100 {
			return errors.New("độ khớp của điều kiện phải trong 1–100%")
		}
	case model.CondPixel:
		if len(c.Points) == 0 {
			return errors.New("điều kiện màu chưa chọn điểm")
		}
		if !inside(c.Points[0]) {
			return errors.New("điểm của điều kiện màu nằm ngoài vùng")
		}
		if _, err := vision.ParseColor(c.Color); err != nil {
			return err
		}
	case model.CondValue:
		if !slices.Contains(condOps, c.Op) {
			return fmt.Errorf("phép so không hợp lệ %q", c.Op)
		}
	case model.CondAll, model.CondAny:
		if len(c.Conds) == 0 {
			return errors.New("nhóm điều kiện chưa có điều kiện nào")
		}
		for i, sub := range c.Conds {
			if err := validateCond(sub, inside); err != nil {
				return fmt.Errorf("điều kiện %d: %w", i+1, err)
			}
		}
	}
	return nil
}

// cloneCond chép sâu (thay biến trên bản sao không đụng kịch bản gốc).
func cloneCond(c *model.Cond) *model.Cond {
	if c == nil {
		return nil
	}
	x := *c
	x.Points = slices.Clone(c.Points)
	if c.Conds != nil {
		x.Conds = make([]model.Cond, len(c.Conds))
		for i := range c.Conds {
			x.Conds[i] = *cloneCond(&c.Conds[i])
		}
	}
	return &x
}

var caseOps = []string{"eq", "contains", "regex"}

// CaseName – tên hiện của case: nhãn, không có thì key.
func CaseName(c model.Case) string {
	if c.Label != "" {
		return c.Label
	}
	return c.Key
}

// checkControl kiểm cấu trúc node điều khiển khi mở file (giá trị lạ → lỗi; thiếu dữ liệu để lúc chạy báo).
func checkControl(s *model.Step) error {
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
	if s.EndStatus != "" && s.EndStatus != model.EndSuccess && s.EndStatus != model.EndFail {
		return fmt.Errorf("kết quả kết thúc không hợp lệ: %q", s.EndStatus)
	}
	if s.VarOp != "" && s.VarOp != model.VarSet && s.VarOp != model.VarAdd {
		return fmt.Errorf("cách đặt biến không hợp lệ: %q", s.VarOp)
	}
	if s.Max < 0 {
		return errors.New("số lần tối đa không âm")
	}
	return nil
}

// validateControl – node điều khiển đủ dữ liệu để chạy (validateScript gọi). inside: điểm nằm trong vùng.
func validateControl(i int, s model.Step, inside func(model.Point) bool) error {
	switch s.Type {
	case model.StepIf:
		if s.Cond == nil {
			return stepErr(i, s, "cond", "chưa có điều kiện")
		}
		if err := validateCond(*s.Cond, inside); err != nil {
			return stepErr(i, s, "cond", "%v", err)
		}
		if s.TimeoutMs < 0 || s.TimeoutMs > 0 && s.IntervalMs <= 0 {
			return stepErr(i, s, "timeoutMs", "thời gian chờ không âm, kiểm tra mỗi phải > 0")
		}
	case model.StepSwitch, model.StepSwitchScreen:
		if len(s.Cases) == 0 {
			return stepErr(i, s, "cases", "chưa có case nào")
		}
		for j, c := range s.Cases {
			if s.Type == model.StepSwitch {
				if !slices.Contains(caseOps, c.Op) {
					return stepErr(i, s, "cases", "case %s: cách so không hợp lệ %q", CaseName(c), c.Op)
				}
				if c.Op == "regex" {
					if _, err := regexp.Compile(c.Match); err != nil {
						return stepErr(i, s, "cases", "case %s: mẫu sai – %v", CaseName(c), err)
					}
				}
				continue
			}
			if c.Cond == nil || !isLeafScreen(c.Cond.Kind) && c.Cond.Kind != model.CondAll && c.Cond.Kind != model.CondAny {
				return stepErr(i, s, "cases", "case %d (%s): cần điều kiện thấy chữ / ảnh / màu", j+1, CaseName(c))
			}
			if err := validateCond(*c.Cond, inside); err != nil {
				return stepErr(i, s, "cases", "case %s: %v", CaseName(c), err)
			}
		}
		if s.Type == model.StepSwitchScreen && (s.TimeoutMs < 0 || s.IntervalMs <= 0) {
			return stepErr(i, s, "timeoutMs", "thời gian chờ không âm, kiểm tra mỗi phải > 0")
		}
	case model.StepCounter:
		if s.Max < 1 {
			return stepErr(i, s, "max", "số lần tối đa phải ≥ 1")
		}
	case model.StepEnd:
		if s.EndStatus != model.EndSuccess && s.EndStatus != model.EndFail {
			return stepErr(i, s, "endStatus", "chọn kết thúc thành công hay thất bại")
		}
	case model.StepSetVar:
		if !varNameRe.MatchString(s.VarName) {
			return stepErr(i, s, "varName", "tên biến %q không hợp lệ – chữ cái / _ đầu, sau đó chữ / số / _", s.VarName)
		}
		if s.VarOp == model.VarAdd {
			if _, err := ParseNumber(s.Value); err != nil && !NodeVarRe.MatchString(s.Value) && !varRe.MatchString(s.Value) {
				return stepErr(i, s, "value", "cộng cần một số, có “%s”", s.Value)
			}
		} else if s.VarOp != model.VarSet {
			return stepErr(i, s, "varOp", "chọn gán hay cộng")
		}
	}
	return nil
}
