package model

// Node Bộ gõ (bước 43 S6, engine/inputsource.go): đọc / đổi bộ gõ – bàn phím đang dùng của hệ điều hành.

// Thao tác của node Bộ gõ (Step.InputOp).
const (
	InputRead = ""    // đọc bộ gõ hiện tại + danh sách
	InputSet  = "set" // đổi sang bộ gõ InputLayout
)

var InputOps = []string{InputRead, InputSet}
