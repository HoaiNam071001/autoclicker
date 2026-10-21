package script

import (
	"fmt"
	"slices"

	"ghosthand/internal/model"
)

// Trường số lấy từ biến (bước 20 F5c, Step.NumVars): ô chữ biểu thức (thay biến như ô chữ khác – VarFields), lúc chạy
// tính ra số rồi gán vào trường (setStepField – cùng khoảng hợp lệ như ghim input).

// checkNumVars (mở file): trường nhận biến được với loại node, mỗi trường một dòng.
func checkNumVars(s *model.Step) error {
	allowed := model.NumVarFieldsOf(s.Type)
	seen := map[string]bool{}
	for _, v := range s.NumVars {
		if !slices.Contains(allowed, v.Field) {
			return fmt.Errorf("numVars: trường %q của node %s không nhận biến (được: %v)", v.Field, s.Type, allowed)
		}
		if seen[v.Field] {
			return fmt.Errorf("numVars: trường %q lặp lại", v.Field)
		}
		seen[v.Field] = true
	}
	return nil
}

// checkNumVarPins – trường đã lấy từ biến thì không ghim input được nữa (hai nguồn cho một trường).
func checkNumVarPins(sc model.Script, i int, s model.Step) error {
	for _, in := range sc.Inputs {
		for _, t := range in.Targets {
			if t.StepID == s.ID && s.NumVarOf(t.Field) != "" {
				return stepErr(i, s, t.Field, "trường vừa lấy từ biến vừa ghim input %s – bỏ một trong hai", in.Key)
			}
		}
	}
	return nil
}

// withNumVarPlaceholders – bản sao s với trường lấy từ biến đặt giá trị hợp lệ tạm (1), để kiểm tra tĩnh (amount ≥ 1,
// timeout > 0…) không báo sai khi giá trị chỉ biết lúc chạy.
func withNumVarPlaceholders(s model.Step) model.Step {
	for _, v := range s.NumVars {
		setStepField(&s, v.Field, 1)
	}
	return s
}

// ApplyNumVars tính từng biểu thức (đã thay biến) và gán vào trường số của s. Lỗi: không ra số, số âm, ngoài khoảng.
func ApplyNumVars(s *model.Step) error {
	for _, v := range s.NumVars {
		n, err := EvalCoord(v.Expr)
		if err == nil && n < 0 {
			err = fmt.Errorf("ra %d – không được âm", n)
		}
		if err == nil {
			err = setStepField(s, v.Field, n)
		}
		if err != nil {
			return fmt.Errorf("%s: %w", v.Field, err)
		}
	}
	return nil
}
