package script

import (
	"fmt"
	"regexp"
	"slices"

	"autoclicker/internal/model"
)

// Input của kịch bản (bước 21b): giá trị người soạn cho phép đổi khi bấm Chạy.
// Nguồn: trường được ghim (Targets) và biến {{key}} trong trường chữ (gõ chữ, URL / app của bước Mở).
// ResolveInputs áp giá trị vào bản sao kịch bản trước khi chạy – engine không biết gì về input.

func inputKind(typ string) string {
	switch typ {
	case model.InputNumber, model.InputMs:
		return "int"
	case model.InputBool:
		return "bool"
	}
	return "string"
}

func validInputType(t string) bool {
	return slices.Contains([]string{model.InputText, model.InputMultiline, model.InputNumber, model.InputMs, model.InputBool, model.InputSelect}, t)
}

var inputKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// varRe: \{{ là "{{" thật; {{key}} là biến.
var varRe = regexp.MustCompile(`\\\{\{|\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

func inputName(in model.ScriptInput) string {
	if in.Label != "" && in.Label != in.Key {
		return fmt.Sprintf("%s (%s)", in.Label, in.Key)
	}
	return in.Key
}

// checkInputs (mở file): key hợp lệ, không trùng, kiểu đúng, trường ghim tồn tại và hợp kiểu.
// Ghim vào bước không còn (sửa tay file) thì bỏ ghim đó. Không bắt buộc đã có bước dùng – bản soạn dở vẫn mở được.
func checkInputs(sc *model.Script) error {
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
		targets := []model.InputTarget{}
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

// ResolveInputs trả bản sao kịch bản đã áp giá trị input: values[key] (thiếu / trống thì lấy mặc định),
// ghi vào các trường ghim và thay {{key}} trong ô chữ. Lỗi: thiếu bắt buộc, sai kiểu / ngoài khoảng, biến không có
// trong inputs, input không bước nào dùng (ghim vào bước đã xoá không tính), giá trị cho input không tồn tại. Sau đó vẫn cần validateScript.
func ResolveInputs(sc model.Script, values map[string]any) (model.Script, error) {
	byKey := map[string]model.ScriptInput{}
	for _, in := range sc.Inputs {
		byKey[in.Key] = in
	}
	for k := range values {
		if _, ok := byKey[k]; !ok {
			return model.Script{}, &model.ScriptError{Input: k, Msg: fmt.Sprintf("kịch bản không có input %q", k)}
		}
	}

	vals := map[string]any{} // giá trị đã đổi kiểu; không có key = trống (không bắt buộc)
	for _, in := range sc.Inputs {
		v, ok := values[in.Key]
		if !ok || IsEmptyValue(v) {
			v = in.Default
		}
		if IsEmptyValue(v) {
			if in.Required {
				return model.Script{}, inputErr(in, "chưa nhập (bắt buộc)")
			}
			continue
		}
		cv, err := coerceInput(in, v)
		if err != nil {
			return model.Script{}, err
		}
		vals[in.Key] = cv
	}

	out := sc
	out.Steps = make([]model.Step, len(sc.Steps))
	for i, s := range sc.Steps {
		out.Steps[i] = CloneStep(s)
	}
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
					return model.Script{}, stepErr(i, out.Steps[i], t.Field, "%v", err)
				}
			}
		}
	}
	for i := range out.Steps {
		s := &out.Steps[i]
		for _, f := range VarFields(s) {
			text, err := ExpandVars(*f.Ptr, func(key string) (string, bool) {
				if _, ok := byKey[key]; !ok {
					return "", false
				}
				used[key] = true
				// URL: mã hoá (dấu cách, tiếng Việt, &); thân JSON: thoát dấu nháy – không làm hỏng ô
				return EscapeValue(f, ValueString(vals[key])), true
			})
			if err != nil {
				return model.Script{}, stepErr(i, *s, f.Name, "%v", err)
			}
			*f.Ptr = text
		}
	}
	for _, in := range sc.Inputs {
		if !used[in.Key] {
			return model.Script{}, inputErr(in, "không bước nào dùng – ghim vào một trường, chèn {{%s}} vào ô chữ, hoặc xoá input", in.Key)
		}
	}
	return out, nil
}

// ExpandVars thay {{key}} bằng lookup(key); \{{ thành {{. Biến lookup không biết → lỗi.
func ExpandVars(text string, lookup func(key string) (string, bool)) (string, error) {
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

// WithInputDefaults – bản sao node với {{key}} của input thay bằng giá trị mặc định (Thử tìm / Thử node).
func WithInputDefaults(sc model.Script, s model.Step) (model.Step, error) {
	s = CloneStep(s)
	for _, f := range VarFields(&s) {
		t, err := ExpandVars(*f.Ptr, func(key string) (string, bool) {
			for _, in := range sc.Inputs {
				if in.Key == key && !IsEmptyValue(in.Default) {
					return ValueString(in.Default), true
				}
			}
			return "", false
		})
		if err != nil {
			return s, fmt.Errorf("thử dùng giá trị mặc định của input: %w", err)
		}
		*f.Ptr = t
	}
	return s, nil
}
