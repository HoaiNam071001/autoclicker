package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Input của kịch bản (bước 21b): giá trị người soạn cho phép đổi khi bấm Chạy.
// Nguồn: trường được ghim (Targets) và biến {{key}} trong trường chữ (bước gõ chữ).
// resolveInputs áp giá trị vào bản sao kịch bản trước khi chạy – engine không biết gì về input.

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
	"enabled":       "bool",
}

func inputKind(typ string) string {
	switch typ {
	case InputNumber, InputMs:
		return "int"
	case InputBool:
		return "bool"
	}
	return "string"
}

func validInputType(t string) bool {
	return slices.Contains([]string{InputText, InputMultiline, InputNumber, InputMs, InputBool, InputSelect}, t)
}

var inputKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// varRe: \{{ là "{{" thật; {{key}} là biến.
var varRe = regexp.MustCompile(`\\\{\{|\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

// ScriptError – lỗi kiểm tra kịch bản có cấu trúc (bước / trường / input) để UI, CLI và AI (MCP, bước 31) biết chỗ sửa.
type ScriptError struct {
	Step   int    `json:"step,omitempty"` // thứ tự bước từ 1, 0 = không thuộc bước nào
	StepID string `json:"stepId,omitempty"`
	Field  string `json:"field,omitempty"`
	Input  string `json:"input,omitempty"` // key input liên quan
	Msg    string `json:"message"`
}

func (e *ScriptError) Error() string {
	if e.Step > 0 {
		return fmt.Sprintf("bước %d: %s", e.Step, e.Msg)
	}
	return e.Msg
}

func stepErr(i int, s Step, field, format string, a ...any) *ScriptError {
	return &ScriptError{Step: i + 1, StepID: s.ID, Field: field, Msg: fmt.Sprintf(format, a...)}
}

func inputErr(in ScriptInput, format string, a ...any) *ScriptError {
	return &ScriptError{Input: in.Key, Msg: fmt.Sprintf("%s: %s", inputName(in), fmt.Sprintf(format, a...))}
}

func inputName(in ScriptInput) string {
	if in.Label != "" && in.Label != in.Key {
		return fmt.Sprintf("%s (%s)", in.Label, in.Key)
	}
	return in.Key
}

// checkInputs (mở file): key hợp lệ, không trùng, kiểu đúng, trường ghim tồn tại và hợp kiểu.
// Ghim vào bước không còn (sửa tay file) thì bỏ ghim đó. Không bắt buộc đã có bước dùng – bản soạn dở vẫn mở được.
func checkInputs(sc *Script) error {
	ids := map[string]bool{}
	for _, s := range sc.Steps {
		ids[s.ID] = true
	}
	seen := map[string]bool{}
	for i := range sc.Inputs {
		in := &sc.Inputs[i]
		if !inputKeyRe.MatchString(in.Key) {
			return fmt.Errorf("input %q: key chỉ gồm chữ không dấu, số, _ và không bắt đầu bằng số", in.Key)
		}
		if seen[in.Key] {
			return fmt.Errorf("input %q bị trùng", in.Key)
		}
		seen[in.Key] = true
		if !validInputType(in.Type) {
			return fmt.Errorf("input %q: kiểu không hợp lệ %q", in.Key, in.Type)
		}
		targets := []InputTarget{}
		for _, t := range in.Targets {
			kind, ok := pinFields[t.Field]
			if !ok {
				return fmt.Errorf("input %q: không ghim được trường %q", in.Key, t.Field)
			}
			if kind != inputKind(in.Type) {
				return fmt.Errorf("input %q: kiểu %s không ghim được vào trường %q", in.Key, in.Type, t.Field)
			}
			if ids[t.StepID] {
				targets = append(targets, t)
			}
		}
		in.Targets = targets
	}
	return nil
}

// resolveInputs trả bản sao kịch bản đã áp giá trị input: values[key] (thiếu / trống thì lấy mặc định),
// ghi vào các trường ghim và thay {{key}} trong ô chữ. Lỗi: thiếu bắt buộc, sai kiểu / ngoài khoảng, biến không có
// trong inputs, input không bước nào dùng (ghim vào bước đã xoá không tính), giá trị cho input không tồn tại. Sau đó vẫn cần validateScript.
func resolveInputs(sc Script, values map[string]any) (Script, error) {
	byKey := map[string]ScriptInput{}
	for _, in := range sc.Inputs {
		byKey[in.Key] = in
	}
	for k := range values {
		if _, ok := byKey[k]; !ok {
			return Script{}, &ScriptError{Input: k, Msg: fmt.Sprintf("kịch bản không có input %q", k)}
		}
	}

	vals := map[string]any{} // giá trị đã đổi kiểu; không có key = trống (không bắt buộc)
	for _, in := range sc.Inputs {
		v, ok := values[in.Key]
		if !ok || isEmptyValue(v) {
			v = in.Default
		}
		if isEmptyValue(v) {
			if in.Required {
				return Script{}, inputErr(in, "chưa nhập (bắt buộc)")
			}
			continue
		}
		cv, err := coerceInput(in, v)
		if err != nil {
			return Script{}, err
		}
		vals[in.Key] = cv
	}

	out := sc
	out.Steps = slices.Clone(sc.Steps)
	index := map[string]int{}
	for i, s := range out.Steps {
		index[s.ID] = i
	}
	used := map[string]bool{}
	for _, in := range sc.Inputs {
		for _, t := range in.Targets {
			i, ok := index[t.StepID]
			if !ok {
				continue // bước đã xoá (vẫn giữ ghim để hoàn tác xoá bước thì ghim quay lại)
			}
			used[in.Key] = true
			if v, ok := vals[in.Key]; ok {
				if err := setStepField(&out.Steps[i], t.Field, v); err != nil {
					return Script{}, stepErr(i, out.Steps[i], t.Field, "%v", err)
				}
			}
		}
	}
	for i := range out.Steps {
		s := &out.Steps[i]
		if s.Type != StepText {
			continue
		}
		text, err := expandVars(s.Text, func(key string) (string, bool) {
			if _, ok := byKey[key]; !ok {
				return "", false
			}
			used[key] = true
			return valueString(vals[key]), true
		})
		if err != nil {
			return Script{}, stepErr(i, *s, "text", "%v", err)
		}
		s.Text = text
	}
	for _, in := range sc.Inputs {
		if !used[in.Key] {
			return Script{}, inputErr(in, "không bước nào dùng – ghim vào một trường, chèn {{%s}} vào ô chữ, hoặc xoá input", in.Key)
		}
	}
	return out, nil
}

// expandVars thay {{key}} bằng lookup(key); \{{ thành {{. Biến lookup không biết → lỗi.
func expandVars(text string, lookup func(key string) (string, bool)) (string, error) {
	var err error
	out := varRe.ReplaceAllStringFunc(text, func(m string) string {
		if m == `\{{` {
			return "{{"
		}
		key := m[2 : len(m)-2]
		v, ok := lookup(key)
		if !ok && err == nil {
			err = fmt.Errorf("biến {{%s}} không có trong tab Input", key)
		}
		return v
	})
	return out, err
}

func isEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}

// coerceInput đổi giá trị (từ form, JSON, dòng lệnh – số có thể là chuỗi) sang kiểu của input và kiểm tra khoảng / lựa chọn.
func coerceInput(in ScriptInput, v any) (any, error) {
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
	s := valueString(v)
	if in.Type == InputSelect && len(in.Options) > 0 && !slices.Contains(in.Options, s) {
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

func valueString(v any) string {
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

func rangeText(in ScriptInput) string {
	switch {
	case in.Min != nil && in.Max != nil:
		return fmt.Sprintf("%d–%d", *in.Min, *in.Max)
	case in.Min != nil:
		return fmt.Sprintf("≥ %d", *in.Min)
	default:
		return fmt.Sprintf("≤ %d", *in.Max)
	}
}

func setStepField(s *Step, field string, v any) error {
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
	default:
		return fmt.Errorf("không ghim được trường %q", field)
	}
	return nil
}

// ---- Giá trị lần trước + bộ giá trị (preset) – lưu trên máy theo từng kịch bản, không ghi vào file kịch bản ----

const inputMemoryFile = "inputs.json"

type InputPreset struct {
	Name   string         `json:"name"`
	Values map[string]any `json:"values"`
}

type InputMemory struct {
	Last       map[string]any `json:"last,omitempty"`       // giá trị lần chạy trước
	LastPreset string         `json:"lastPreset,omitempty"` // bộ giá trị đã chọn lần trước ("" = tự điền)
	Presets    []InputPreset  `json:"presets,omitempty"`
	AutoRun    bool           `json:"autoRun,omitempty"` // "Lần sau chạy luôn với giá trị này" – bấm Chạy không hỏi
}

func loadInputMemory(dir string) map[string]InputMemory {
	all := map[string]InputMemory{}
	if b, err := os.ReadFile(filepath.Join(dir, inputMemoryFile)); err == nil {
		_ = json.Unmarshal(b, &all) // hỏng thì coi như chưa có
	}
	return all
}

// GetInputMemory trả giá trị lần trước / preset của kịch bản id (đường dẫn file, hoặc tên khi chưa lưu).
func (a *App) GetInputMemory(id string) InputMemory {
	dir, err := configDir()
	if err != nil {
		return InputMemory{}
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	return loadInputMemory(dir)[id]
}

// SaveInputMemory ghi giá trị lần trước / preset của kịch bản id.
func (a *App) SaveInputMemory(id string, m InputMemory) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	all := loadInputMemory(dir)
	all[id] = m
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(filepath.Join(dir, inputMemoryFile), b)
}
