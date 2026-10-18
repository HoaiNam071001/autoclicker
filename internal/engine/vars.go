package engine

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
			return script.EscapeValue(f, v)
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
		key, rest, _ := strings.Cut(strings.TrimPrefix(path, "output."), ".")
		v, ok := r.Output[key]
		if !ok {
			return "", fmt.Errorf("node %s không có output %q (lần chạy gần nhất: %s)", r.Ref, key, r.Status)
		}
		if rest != "" { // đọc sâu: {{n2.output.json.data.items.0.id}} (node HTTP)
			var err error
			if v, err = dig(v, rest); err != nil {
				return "", fmt.Errorf("node %s output %s: %w", r.Ref, key, err)
			}
		}
		return fmtValue(v), nil
	}
	return "", fmt.Errorf("không hiểu %q – dùng status, port, error.message, output.<tên>", path)
}

// dig đi theo đường "a.b.0.c" trong object / mảng (JSON đã giải mã).
func dig(v any, path string) (any, error) {
	for _, k := range strings.Split(path, ".") {
		switch x := v.(type) {
		case map[string]any:
			next, ok := x[k]
			if !ok {
				return nil, fmt.Errorf("không có khoá %q", k)
			}
			v = next
		case []any:
			i, err := strconv.Atoi(k)
			if err != nil || i < 0 || i >= len(x) {
				return nil, fmt.Errorf("mảng có %d phần tử, không có phần tử %q (đếm từ 0)", len(x), k)
			}
			v = x[i]
		default:
			return nil, fmt.Errorf("không đọc sâu được vào %q – giá trị không phải object / mảng", k)
		}
	}
	return v, nil
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
	case map[string]any, []any: // object / mảng (output json của node HTTP) – in ra JSON để dùng tiếp
		b, _ := json.Marshal(x)
		return string(b)
	}
	return fmt.Sprint(v)
}
