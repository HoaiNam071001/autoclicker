package engine

import (
	"fmt"
	"strconv"
	"strings"

	"autoclicker/internal/launch"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Biến lúc chạy (bước 20): node đọc kết quả node đã chạy trước đó trong cùng vòng.
//   {{n2.status}}  {{n2.port}}  {{n2.error.code}}  {{n2.error.message}}  {{n2.output.<key>}}
//   {{last.…}} = node vừa chạy ngay trước;  {{vars.<tên>}} = biến do node Đặt biến gán (F3)
// Input {{key}} (không dấu chấm) đã thay trước khi chạy (resolveInputs).

// runVars – bảng tra khi thay biến: kết quả theo mã node, node vừa chạy, biến người dùng đặt.
type runVars struct {
	byRef map[string]*model.StepResult
	last  *model.StepResult
	vars  map[string]any
}

// resolveNodeVars trả bản sao node với biến node trong các ô chữ đã thay giá trị. Biến trỏ node chưa chạy / output
// không có → lỗi invalid_input (node thất bại, đi cổng fail) – không lặng lẽ gõ chữ trống.
func resolveNodeVars(s model.Step, rv runVars) (model.Step, error) {
	s = script.CloneStep(s)
	var firstErr error
	for _, f := range script.VarFields(&s) {
		*f.Ptr = script.NodeVarRe.ReplaceAllStringFunc(*f.Ptr, func(m string) string {
			sub := script.NodeVarRe.FindStringSubmatch(m)
			v, err := rv.lookup(sub[1], sub[2])
			if err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("%s: biến %s – %w", f.Name, m, err)
				}
				return m
			}
			if f.URL {
				return launch.EscapeURLValue(v)
			}
			return v
		})
	}
	return s, firstErr
}

func (rv runVars) lookup(name, path string) (string, error) {
	if name == "vars" {
		v, ok := rv.vars[path]
		if !ok {
			return "", fmt.Errorf("chưa có biến %q (chưa node Đặt biến nào gán)", path)
		}
		return fmtValue(v), nil
	}
	r := rv.last
	if name != "last" {
		r = rv.byRef[name]
	}
	if r == nil {
		if name == "last" {
			return "", fmt.Errorf("chưa có node nào chạy trước")
		}
		return "", fmt.Errorf("node %s chưa chạy trong vòng này", name)
	}
	switch {
	case path == "status":
		return r.Status, nil
	case path == "port":
		return r.Port, nil
	case path == "error.code" || path == "error.message":
		if r.Error == nil {
			return "", nil
		}
		if path == "error.code" {
			return r.Error.Code, nil
		}
		return r.Error.Message, nil
	case strings.HasPrefix(path, "output."):
		key := strings.TrimPrefix(path, "output.")
		v, ok := r.Output[key]
		if !ok {
			return "", fmt.Errorf("node %s không có output %q (lần chạy gần nhất: %s)", r.Ref, key, r.Status)
		}
		return fmtValue(v), nil
	}
	return "", fmt.Errorf("không hiểu %q – dùng status, port, error.message, output.<tên>", path)
}

func fmtValue(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	case nil:
		return ""
	}
	return fmt.Sprint(v)
}
