package script

import (
	"fmt"

	"ghosthand/internal/model"
)

// Lỗi kiểm tra có cấu trúc (model.ScriptError): node / trường / input nào sai – UI và AI chỉ đúng chỗ.

func stepErr(i int, s model.Step, field, format string, a ...any) *model.ScriptError {
	return &model.ScriptError{Step: i + 1, StepID: s.ID, Ref: s.Ref, Field: field, Msg: fmt.Sprintf(format, a...)}
}
func inputErr(in model.ScriptInput, format string, a ...any) *model.ScriptError {
	return &model.ScriptError{Input: in.Key, Msg: fmt.Sprintf("%s: %s", inputName(in), fmt.Sprintf(format, a...))}
}
