package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"autoclicker/internal/model"
)

// Giá trị input: ép theo kiểu (số / ms / bool / chọn), chữ hiển thị, gán vào trường được ghim của node.

// pinFields – trường ghim được và kiểu giá trị: "int" (input number / ms) hoặc "bool".
// Không ghim được trường tự tính (kéo thả "từ điểm trước") hay danh sách (điểm, phím).
var pinFields = map[string]string{
	"clicks":        "int",
	"clickInterval": "int",
	"holdMs":        "int",
	"dragMs":        "int",
	"delayAfter":    "int",
	"jitterMs":      "int",
	"jitterPx":      "int",
	"amount":        "int",
	"timeoutMs":     "int",
	"threshold":     "int",
	"enabled":       "bool",
}

// IsEmptyValue: nil hoặc chữ chỉ có khoảng trắng.
func IsEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}

// coerceInput đổi giá trị (từ form, JSON, dòng lệnh – số có thể là chuỗi) sang kiểu của input và kiểm tra khoảng / lựa chọn.
func coerceInput(in model.ScriptInput, v any) (any, error) {
	switch inputKind(in.Type) {
	case "int":
		n, err := toInt(v)
		if err != nil {
			return nil, inputErr(in, "%v", err)
		}
		if n < 0 {
			return nil, inputErr(in, "không được âm")
		}
		if in.Min != nil && n < *in.Min || in.Max != nil && n > *in.Max {
			return nil, inputErr(in, "%d nằm ngoài khoảng %s", n, rangeText(in))
		}
		return n, nil
	case "bool":
		switch b := v.(type) {
		case bool:
			return b, nil
		case string:
			if p, err := strconv.ParseBool(strings.TrimSpace(b)); err == nil {
				return p, nil
			}
		}
		return nil, inputErr(in, "cần bật / tắt (true / false), nhận %v", v)
	}
	s := ValueString(v)
	if in.Type == model.InputSelect && len(in.Options) > 0 && !slices.Contains(in.Options, s) {
		return nil, inputErr(in, "%q không có trong danh sách chọn (%s)", s, strings.Join(in.Options, ", "))
	}
	return s, nil
}
func toInt(v any) (int, error) {
	switch n := v.(type) {
	case int:
		return n, nil
	case float64:
		if n != math.Trunc(n) {
			return 0, fmt.Errorf("cần số nguyên, nhận %v", n)
		}
		return int(n), nil
	case json.Number:
		return toInt(string(n))
	case string:
		i, err := strconv.Atoi(strings.TrimSpace(n))
		if err != nil {
			return 0, fmt.Errorf("cần số nguyên, nhận %q", n)
		}
		return i, nil
	}
	return 0, fmt.Errorf("cần số nguyên, nhận %v", v)
}

// ValueString – giá trị input dạng chữ (chèn vào {{key}}, hiện trong lỗi).
func ValueString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
	return fmt.Sprint(v)
}
func rangeText(in model.ScriptInput) string {
	switch {
	case in.Min != nil && in.Max != nil:
		return fmt.Sprintf("%d–%d", *in.Min, *in.Max)
	case in.Min != nil:
		return fmt.Sprintf("≥ %d", *in.Min)
	default:
		return fmt.Sprintf("≤ %d", *in.Max)
	}
}
func setStepField(s *model.Step, field string, v any) error {
	if field == "enabled" {
		s.Enabled = v.(bool)
		return nil
	}
	n, ok := v.(int)
	if !ok {
		return fmt.Errorf("trường %q cần số", field)
	}
	switch field {
	case "clicks":
		if n < 1 {
			return errors.New("số click phải ≥ 1")
		}
		s.Clicks = n
	case "clickInterval":
		s.ClickInterval = n
	case "holdMs":
		s.HoldMs = n
	case "dragMs":
		s.DragMs = n
	case "delayAfter":
		s.DelayAfter = n
	case "jitterMs":
		s.JitterMs = n
	case "jitterPx":
		s.JitterPx = n
	case "timeoutMs":
		if n < 1 {
			return errors.New("thời gian chờ tối đa phải > 0")
		}
		s.TimeoutMs = n
	case "threshold":
		if n < 1 || n > 100 {
			return errors.New("độ khớp tối thiểu phải trong 1–100%")
		}
		s.Threshold = n
	case "amount":
		if n < 1 {
			return errors.New("số nấc cuộn phải ≥ 1")
		}
		s.Amount = n
	default:
		return fmt.Errorf("không ghim được trường %q", field)
	}
	return nil
}
