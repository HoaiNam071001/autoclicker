package model

// Node Bộ gõ (bước 43 S6, engine/inputsource.go): đọc / đổi bộ gõ – bàn phím đang dùng của hệ điều hành.

// Thao tác của node Bộ gõ (Step.InputOp).
const (
	InputRead = ""    // đọc bộ gõ hiện tại + danh sách
	InputSet  = "set" // đổi sang bộ gõ InputLayout
	// InputToggleIME bấm phím tắt bật / tắt tiếng Việt của phần mềm gõ (Unikey, EVKey, OpenKey…) – bước 43 S6b. Phần mềm
	// đó không có API đọc trạng thái nên đây là BẤM MỘT LẦN (đảo trạng thái), không phải "bật" hay "tắt" chắc chắn.
	InputToggleIME = "toggleIME"
)

var InputOps = []string{InputRead, InputSet, InputToggleIME}

// InputDefaultHotkey – phím tắt đổi chế độ mặc định của Unikey / EVKey / OpenKey (người dùng có thể đã đổi → InputHotkey).
const InputDefaultHotkey = "ctrl+shift"
