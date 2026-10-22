package script

import (
	"fmt"
	"slices"
	"strings"

	"ghosthand/internal/model"
)

// checkProcess – miền giá trị node Tiến trình khi mở file.
func checkProcess(s model.Step) error {
	if !slices.Contains(model.ProcOps, s.ProcOp) {
		return fmt.Errorf("thao tác tiến trình không hợp lệ: %q (trống = đọc | waitStart | waitQuit | waitIdle)", s.ProcOp)
	}
	if !slices.Contains(model.NameMatches, s.ProcMatch) {
		return fmt.Errorf("cách so khớp tên không hợp lệ: %q (trống = chứa | exact | regex)", s.ProcMatch)
	}
	if s.TimeoutMs < 0 || s.StableMs < 0 || s.CPUBelow < 0 {
		return fmt.Errorf("thời gian chờ / ngưỡng CPU không được âm")
	}
	return nil
}

// validateProcess – đủ dữ liệu để chạy: chỉ thao tác đọc được bỏ trống tên.
func validateProcess(i int, s model.Step) error {
	if s.ProcOp != model.ProcRead && strings.TrimSpace(s.ProcName) == "" {
		return stepErr(i, s, "procName", "chưa nhập tên tiến trình")
	}
	return nil
}
