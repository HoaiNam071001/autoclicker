package model

import "fmt"

// Kiểu input (ScriptInput.Type).
const (
	InputText      = "text"
	InputMultiline = "multiline"
	InputNumber    = "number"
	InputMs        = "ms"
	InputBool      = "bool"
	InputSelect    = "select"
)

// InputTarget – trường của một bước được ghim vào input.
type InputTarget struct {
	StepID string `json:"stepId"`
	Field  string `json:"field"` // tên json của trường trong Step, xem pinFields
}

type ScriptInput struct {
	Key      string        `json:"key"` // chữ không dấu, số, _ – dùng trong {{key}}
	Label    string        `json:"label"`
	Type     string        `json:"type"`              // text | multiline | number | ms | bool | select
	Default  any           `json:"default,omitempty"` // số / chữ / bool theo Type; trống = không có mặc định
	Required bool          `json:"required,omitempty"`
	Min      *int          `json:"min,omitempty"` // number / ms
	Max      *int          `json:"max,omitempty"`
	Options  []string      `json:"options,omitempty"` // select
	Hint     string        `json:"hint,omitempty"`    // gợi ý dưới ô trong form Chạy
	Targets  []InputTarget `json:"targets,omitempty"`
}

// ScriptError – lỗi kiểm tra kịch bản có cấu trúc (bước / trường / input) để UI, CLI và AI (MCP, bước 31) biết chỗ sửa.
type ScriptError struct {
	Step   int    `json:"step,omitempty"` // vị trí node trong mảng steps từ 1, 0 = không thuộc node nào
	StepID string `json:"stepId,omitempty"`
	Ref    string `json:"ref,omitempty"` // mã node (n2) – hiện trong lời báo lỗi
	Field  string `json:"field,omitempty"`
	Edge   string `json:"edgeId,omitempty"` // id cạnh lỗi (Field = edges)
	Input  string `json:"input,omitempty"`  // key input liên quan
	Msg    string `json:"message"`
}

func (e *ScriptError) Error() string {
	if e.Ref != "" {
		return fmt.Sprintf("%s: %s", e.Ref, e.Msg)
	}
	if e.Step > 0 {
		return fmt.Sprintf("bước %d: %s", e.Step, e.Msg)
	}
	return e.Msg
}
