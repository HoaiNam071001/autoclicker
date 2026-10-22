package script

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"ghosthand/internal/model"
)

// Node Nguồn (bước 43 S3) và Màn hình / âm thanh (bước 43 S7).

func checkPower(s model.Step) error {
	if !slices.Contains(model.PowerOps, s.PowerOp) {
		return fmt.Errorf("hành động nguồn không hợp lệ: %q (lock | displaySleep | sleep | logout | shutdown | restart)", s.PowerOp)
	}
	if s.CountdownSec < 0 {
		return fmt.Errorf("thời gian đếm ngược không được âm")
	}
	return nil
}

func checkDisplay(s model.Step) error {
	if !slices.Contains(model.DisplayOps, s.DisplayOp) {
		return fmt.Errorf("thao tác màn hình / âm thanh không hợp lệ: %q (trống = đọc | theme | brightness | audioOut | audioIn | sound)", s.DisplayOp)
	}
	return nil
}

// validateDisplay – đủ dữ liệu để chạy. Ô còn biến thì phần còn lại kiểm tra lúc chạy.
func validateDisplay(i int, s model.Step) error {
	v := strings.TrimSpace(s.DisplayValue)
	if s.DisplayOp == model.DisplayRead {
		return nil
	}
	if v == "" {
		what := map[string]string{
			model.DisplayTheme: "chế độ (light | dark | toggle)", model.DisplayBrightness: "độ sáng 0–100",
			model.DisplayAudioOut: "tên loa", model.DisplayAudioIn: "tên micro", model.DisplaySound: "file hoặc tên âm",
		}[s.DisplayOp]
		return stepErr(i, s, "displayValue", "chưa nhập %s", what)
	}
	if strings.Contains(v, "{{") {
		return nil
	}
	switch s.DisplayOp {
	case model.DisplayTheme:
		if !slices.Contains(model.DisplayThemes, v) {
			return stepErr(i, s, "displayValue", "chế độ “%s” không hợp lệ – chọn light, dark hoặc toggle", v)
		}
	case model.DisplayBrightness:
		if n, err := strconv.Atoi(v); err != nil || n < 0 || n > 100 {
			return stepErr(i, s, "displayValue", "độ sáng “%s” không hợp lệ – số nguyên 0–100", v)
		}
	}
	return nil
}
