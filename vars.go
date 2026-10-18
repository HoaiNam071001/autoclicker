package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"autoclicker/internal/launch"
	"autoclicker/internal/model"
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
	s = cloneStep(s)
	var firstErr error
	for _, f := range varFields(&s) {
		*f.p = nodeVarRe.ReplaceAllStringFunc(*f.p, func(m string) string {
			sub := nodeVarRe.FindStringSubmatch(m)
			v, err := rv.lookup(sub[1], sub[2])
			if err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("%s: biến %s – %w", f.name, m, err)
				}
				return m
			}
			if f.url {
				return launch.EscapeURLValue(v)
			}
			return v
		})
	}
	return s, firstErr
}

// varField – ô chữ nhận {{biến}} của node (kể cả trong điều kiện / case). url: giá trị chèn vào phải mã hoá.
type varField struct {
	name string // tên trường (json) cho lời báo lỗi, vd "cond.conds[1].left"
	p    *string
	url  bool
}

// varFields – các ô chữ nhận biến của s (input {{key}} ở resolveInputs, biến node ở resolveNodeVars, kiểm tra ở
// validateFlow). Con trỏ trỏ vào s – muốn sửa thì gọi trên bản sao (cloneStep).
func varFields(s *model.Step) []varField {
	var out []varField
	add := func(name string, p *string, url bool) { out = append(out, varField{name, p, url}) }
	switch s.Type {
	case model.StepText:
		add("text", &s.Text, false)
	case model.StepWaitImage, model.StepClickImage:
		if s.FindBy == model.FindByText {
			add("findText", &s.FindText, false)
		}
	case model.StepOpen:
		add("target", &s.Target, launch.TargetKind(s.Target) == launch.TargetURL)
	case model.StepIf:
		if s.Cond != nil {
			condFields(s.Cond, "cond", add)
		}
	case model.StepSwitch:
		add("value", &s.Value, false)
		for i := range s.Cases {
			add(fmt.Sprintf("cases[%d].match", i), &s.Cases[i].Match, false)
		}
	case model.StepSwitchScreen:
		for i := range s.Cases {
			if c := s.Cases[i].Cond; c != nil {
				condFields(c, fmt.Sprintf("cases[%d].cond", i), add)
			}
		}
	case model.StepEnd:
		add("message", &s.Message, false)
	case model.StepSetVar:
		add("value", &s.Value, false)
	}
	return out
}

func condFields(c *model.Cond, name string, add func(string, *string, bool)) {
	switch c.Kind {
	case model.CondText:
		add(name+".findText", &c.FindText, false)
	case model.CondValue:
		add(name+".left", &c.Left, false)
		add(name+".right", &c.Right, false)
	case model.CondAll, model.CondAny:
		for i := range c.Conds {
			condFields(&c.Conds[i], fmt.Sprintf("%s.conds[%d]", name, i), add)
		}
	}
}

// cloneStep chép sâu phần lồng nhau sửa được (điểm, điều kiện, case) – thay biến trên bản sao không đụng kịch bản gốc.
func cloneStep(s model.Step) model.Step {
	s.Points = slices.Clone(s.Points)
	s.Cond = cloneCond(s.Cond)
	if s.Cases != nil {
		cs := make([]model.Case, len(s.Cases))
		for i, c := range s.Cases {
			c.Cond = cloneCond(c.Cond)
			cs[i] = c
		}
		s.Cases = cs
	}
	return s
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
