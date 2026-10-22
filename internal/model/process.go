package model

// Node Tiến trình (bước 43 S5, engine/process.go): làm với cả tiến trình nền không có cửa sổ (khác node Ứng dụng).

// Thao tác của node Tiến trình (Step.ProcOp).
const (
	ProcRead      = ""          // đọc: tiến trình khớp + CPU / RAM
	ProcWaitStart = "waitStart" // chờ có tiến trình khớp
	ProcWaitQuit  = "waitQuit"  // chờ mọi tiến trình khớp thoát
	ProcWaitIdle  = "waitIdle"  // chờ CPU tổng của tiến trình khớp dưới CPUBelow liên tục StableMs
)

var ProcOps = []string{ProcRead, ProcWaitStart, ProcWaitQuit, ProcWaitIdle}

// ProcIdleCPU – ngưỡng CPU mặc định của waitIdle (%; 100 = đầy một lõi).
const ProcIdleCPU = 5
