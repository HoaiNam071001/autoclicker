package script

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"ghosthand/internal/model"
)

// Kiểm tra node Chạy lệnh (bước 43 S1): miền giá trị (checkShell, lúc mở file) và dữ liệu đủ để chạy (validateShell).

// envName – tên biến môi trường hợp lệ cho cả shell lẫn PowerShell.
var envName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func checkShell(s model.Step) error {
	if !slices.Contains(model.ShellLangs, s.ShellLang) {
		return fmt.Errorf("ngôn ngữ script không hợp lệ: %q (trống = shell của máy | powershell | applescript | jxa)", s.ShellLang)
	}
	if s.TimeoutMs < 0 {
		return fmt.Errorf("thời gian chờ không được âm")
	}
	return nil
}

func validateShell(i int, s model.Step) error {
	if strings.TrimSpace(s.Script) == "" {
		return stepErr(i, s, "script", "chưa nhập nội dung lệnh")
	}
	for j, kv := range s.ShellEnv {
		if !envName.MatchString(strings.TrimSpace(kv.Key)) {
			return stepErr(i, s, fmt.Sprintf("shellEnv[%d].key", j), "tên biến môi trường %q không hợp lệ (chữ, số, gạch dưới; không bắt đầu bằng số)", kv.Key)
		}
	}
	return nil
}
