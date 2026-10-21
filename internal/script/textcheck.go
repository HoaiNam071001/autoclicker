package script

import (
	"fmt"
	"slices"
	"strings"

	"ghosthand/internal/model"
)

// Kiểm tra node xử lý chữ (bước 42 G4b): Lọc mẫu, Tách chữ, Làm sạch chữ, Số & ngày.

// checkText – miền giá trị khi mở file (gọi từ checkControl).
func checkText(s *model.Step) error {
	if !slices.Contains(model.CaseModes, s.CaseMode) {
		return fmt.Errorf("kiểu hoa / thường không hợp lệ: %q (trống | lower | upper | title)", s.CaseMode)
	}
	if !slices.Contains(model.NumOps, s.NumOp) {
		return fmt.Errorf("việc của Số & ngày không hợp lệ: %q (trống = lấy số | format | date)", s.NumOp)
	}
	if !slices.Contains(model.NumStyles, s.NumStyle) {
		return fmt.Errorf("kiểu viết số không hợp lệ: %q (trống = 1.234,5 | en | plain)", s.NumStyle)
	}
	if s.Decimals < 0 || s.Decimals > 10 {
		return fmt.Errorf("số chữ số thập phân phải 0–10, có %d", s.Decimals)
	}
	return nil
}

// validateText – dữ liệu đủ để chạy: có chữ vào, regex đúng cú pháp (regex có biến chỉ biết khi chạy).
func validateText(i int, s model.Step) error {
	if strings.TrimSpace(s.Value) == "" && !(s.Type == model.StepTextNumber && s.NumOp == model.NumDate) {
		return stepErr(i, s, "value", "chưa có chữ cần xử lý – thường là biến, vd {{n3.output.text}}")
	}
	switch s.Type {
	case model.StepTextMatch:
		if s.Pattern == "" {
			return stepErr(i, s, "pattern", "chưa nhập mẫu regex")
		}
		return checkRegex(i, s, "pattern", s.Pattern)
	case model.StepTextClean:
		if s.ReplaceRegex {
			return checkRegex(i, s, "replaceFrom", s.ReplaceFrom)
		}
	}
	return nil
}
