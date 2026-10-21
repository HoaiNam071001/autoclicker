package script

import (
	"strings"
	"testing"

	"ghosthand/internal/model"
)

func inputScript() model.Script {
	return model.Script{
		Name: "Nhắn tin",
		Steps: []model.Step{
			{ID: "s1", Type: model.StepClick, Points: []model.Point{{X: 10, Y: 10}}, Clicks: 1, DelayAfter: 500, Enabled: true},
			{ID: "s2", Type: model.StepText, Text: "xin chào {{M1}} và {{M2}}", DelayAfter: 300, Enabled: true},
			{ID: "s3", Type: model.StepClick, Points: []model.Point{{X: 20, Y: 20}}, Clicks: 1, DelayAfter: 500, Enabled: true},
		},
		Inputs: []model.ScriptInput{
			{Key: "click_delay", Label: "Chờ sau khi click", Type: model.InputMs, Default: 500.0,
				Targets: []model.InputTarget{{StepID: "s1", Field: "delayAfter"}, {StepID: "s3", Field: "delayAfter"}}},
			{Key: "M1", Label: "Tên người nhận", Type: model.InputText, Required: true},
			{Key: "M2", Label: "Tên thứ hai", Type: model.InputText, Default: "Lan"},
		},
	}
}

func TestResolveInputs(t *testing.T) {
	sc := inputScript()
	out, err := ResolveInputs(sc, map[string]any{"click_delay": "1000", "M1": "Hoà", "M2": ""})
	if err != nil {
		t.Fatal(err)
	}
	if out.Steps[0].DelayAfter != 1000 || out.Steps[2].DelayAfter != 1000 {
		t.Errorf("một input ghim hai bước: muốn 1000 / 1000, có %d / %d", out.Steps[0].DelayAfter, out.Steps[2].DelayAfter)
	}
	if got := out.Steps[1].Text; got != "xin chào Hoà và Lan" {
		t.Errorf("thay biến (M2 trống lấy mặc định): %q", got)
	}
	if sc.Steps[0].DelayAfter != 500 || sc.Steps[1].Text != "xin chào {{M1}} và {{M2}}" {
		t.Error("resolveInputs sửa vào kịch bản gốc")
	}

	// thiếu giá trị → mặc định; số từ JSON (float64)
	out, err = ResolveInputs(sc, map[string]any{"M1": "A", "click_delay": 250.0})
	if err != nil || out.Steps[0].DelayAfter != 250 {
		t.Fatalf("số float64: %v %d", err, out.Steps[0].DelayAfter)
	}
}

func TestResolveInputsErrors(t *testing.T) {
	cases := []struct {
		name   string
		edit   func(*model.Script)
		values map[string]any
		want   string
	}{
		{"thiếu bắt buộc", nil, map[string]any{}, "Tên người nhận (M1): chưa nhập"},
		{"sai kiểu", nil, map[string]any{"M1": "a", "click_delay": "nhanh"}, "cần số nguyên"},
		{"số lẻ", nil, map[string]any{"M1": "a", "click_delay": 1.5}, "cần số nguyên"},
		{"số âm", nil, map[string]any{"M1": "a", "click_delay": -1.0}, "không được âm"},
		{"input lạ", nil, map[string]any{"M1": "a", "M9": "x"}, `không có input "M9"`},
		{"biến không có", func(sc *model.Script) { sc.Steps[1].Text = "chào {{ten}}" }, map[string]any{"M1": "a"}, "bước 2: biến {{ten}} không có"},
		{"không bước nào dùng", func(sc *model.Script) { sc.Steps[1].Text = "chào {{M1}}" }, map[string]any{"M1": "a"}, "(M2): không bước nào dùng"},
		{"ngoài khoảng", func(sc *model.Script) { m := 100; sc.Inputs[0].Max = &m }, map[string]any{"M1": "a"}, "ngoài khoảng ≤ 100"},
		{"không có trong danh sách", func(sc *model.Script) {
			sc.Inputs[2].Type, sc.Inputs[2].Options = model.InputSelect, []string{"Lan", "Mai"}
		}, map[string]any{"M1": "a", "M2": "Hùng"}, "không có trong danh sách"},
		{"chỉ ghim vào bước đã xoá", func(sc *model.Script) {
			sc.Steps = sc.Steps[1:2]
		}, map[string]any{"M1": "a"}, "(click_delay): không bước nào dùng"},
	}
	for _, c := range cases {
		sc := inputScript()
		if c.edit != nil {
			c.edit(&sc)
		}
		_, err := ResolveInputs(sc, c.values)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: muốn lỗi chứa %q, có %v", c.name, c.want, err)
		}
	}
}

func TestResolveInputsStructuredError(t *testing.T) {
	sc := inputScript()
	sc.Steps[1].Text = "{{lạ}} {{ten}}"
	_, err := ResolveInputs(sc, map[string]any{"M1": "a"})
	se, ok := err.(*model.ScriptError)
	if !ok || se.Step != 2 || se.StepID != "s2" || se.Field != "text" {
		t.Fatalf("muốn ScriptError bước 2 / s2 / text, có %#v", err)
	}
}

func TestExpandVarsEscape(t *testing.T) {
	got, err := ExpandVars(`\{{M1}} là {{M1}}, {{ không phải biến`, func(k string) (string, bool) { return "Tú", k == "M1" })
	if err != nil || got != "{{M1}} là Tú, {{ không phải biến" {
		t.Fatalf("%q %v", got, err)
	}
}

func TestResolveBoolInput(t *testing.T) {
	sc := inputScript()
	sc.Inputs = append(sc.Inputs, model.ScriptInput{Key: "on3", Type: model.InputBool, Default: true, Targets: []model.InputTarget{{StepID: "s3", Field: "enabled"}}})
	out, err := ResolveInputs(sc, map[string]any{"M1": "a", "on3": "false"})
	if err != nil || out.Steps[2].Enabled {
		t.Fatalf("bool từ chuỗi: %v enabled=%v", err, out.Steps[2].Enabled)
	}
}

func TestNoInputsKeepsScript(t *testing.T) {
	sc := sampleScript()
	out, err := ResolveInputs(sc, nil)
	if err != nil || out.Steps[3].Text != "xin chào" {
		t.Fatalf("%v %q", err, out.Steps[3].Text)
	}
}

func TestCheckInputs(t *testing.T) {
	bad := map[string]string{
		"key có dấu":      `{"version":4,"steps":[],"inputs":[{"key":"tên","type":"text"}]}`,
		"key trùng":       `{"version":4,"steps":[],"inputs":[{"key":"a","type":"text"},{"key":"a","type":"text"}]}`,
		"kiểu lạ":         `{"version":4,"steps":[],"inputs":[{"key":"a","type":"date"}]}`,
		"trường lạ":       `{"version":4,"steps":[{"id":"s","type":"wait"}],"inputs":[{"key":"a","type":"ms","targets":[{"stepId":"s","field":"points"}]}]}`,
		"sai kiểu trường": `{"version":4,"steps":[{"id":"s","type":"wait"}],"inputs":[{"key":"a","type":"text","targets":[{"stepId":"s","field":"delayAfter"}]}]}`,
	}
	for name, in := range bad {
		if _, err := DecodeScript([]byte(in)); err == nil {
			t.Errorf("%s: muốn lỗi", name)
		}
	}
	sc, err := DecodeScript([]byte(`{"version":4,"steps":[{"id":"s","type":"wait"}],"inputs":[{"key":"a","type":"ms","targets":[{"stepId":"s","field":"delayAfter"},{"stepId":"x","field":"delayAfter"}]}]}`))
	if err != nil || len(sc.Inputs[0].Targets) != 1 {
		t.Fatalf("bỏ ghim vào bước không còn: %v %+v", err, sc.Inputs)
	}
	if _, err := DecodeScript([]byte(`{"version":3,"steps":[]}`)); err == nil || !strings.Contains(err.Error(), "bản cũ") {
		t.Errorf("file version cũ (danh sách bước) phải bị từ chối: %v", err)
	}
}
