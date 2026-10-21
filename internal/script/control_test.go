package script

import (
	"strings"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

func TestValidateControl(t *testing.T) {
	bad := map[string]model.Step{
		"chưa nhập chữ cần thấy": modeltest.Node("a", "n1", model.StepIf, nil),
		"mẫu sai":                modeltest.Node("a", "n1", model.StepSwitch, func(s *model.Step) { s.Cases = []model.Case{{Key: "a", Op: "regex", Match: "("}} }),
		"chưa có case":           modeltest.Node("a", "n1", model.StepSwitchScreen, func(s *model.Step) { s.Cases = nil }),
		"cần điều kiện thấy":     modeltest.Node("a", "n1", model.StepSwitchScreen, func(s *model.Step) { s.Cases = []model.Case{{Key: "a"}} }),
		"tối đa phải ≥ 1":        modeltest.Node("a", "n1", model.StepCounter, func(s *model.Step) { s.Max = 0 }),
		"tên biến":               modeltest.Node("a", "n1", model.StepSetVar, func(s *model.Step) { s.VarName = "1x" }),
		"cộng cần một số":        modeltest.Node("a", "n1", model.StepSetVar, func(s *model.Step) { s.VarOp, s.Value = model.VarAdd, "abc" }),
		"chưa nhập biểu thức":    modeltest.Node("a", "n1", model.StepCalc, nil),
		"thiếu dấu )":            modeltest.Node("a", "n1", model.StepCalc, func(s *model.Step) { s.Value = "round({{n1.output.x}}" }),
		"chưa có chữ cần xử lý":  modeltest.Node("a", "n1", model.StepTextSplit, nil),
		"chưa nhập mẫu regex":    modeltest.Node("a", "n1", model.StepTextMatch, func(s *model.Step) { s.Value = "x" }),
		"mẫu regex không hợp lệ": modeltest.Node("a", "n1", model.StepTextClean, func(s *model.Step) { s.Value, s.ReplaceFrom, s.ReplaceRegex = "x", "(", true }),
		"không có cổng":          modeltest.Node("a", "n1", model.StepCounter, nil),
	}
	for want, s := range bad {
		sc := model.Script{Start: "a", Steps: []model.Step{s, modeltest.Node("b", "n2", model.StepWait, nil)}}
		if want == "không có cổng" {
			sc.Edges = []model.Edge{modeltest.Edge("a", model.PortSuccess, "b")}
		}
		err := ValidateScript(sc, model.Rect{W: 1920, H: 1080})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi %q, có %v", want, err)
		}
	}
	// cổng case theo key; node Kết thúc không có cổng
	ok := model.Script{Start: "a", Steps: []model.Step{
		modeltest.Node("a", "n1", model.StepSwitch, nil), modeltest.Node("b", "n2", model.StepEnd, nil),
	}, Edges: []model.Edge{modeltest.Edge("a", "case:c2", "b"), modeltest.Edge("a", model.PortDefault, "b")}}
	if err := ValidateScript(ok, model.Rect{W: 1920, H: 1080}); err != nil {
		t.Errorf("flow hợp lệ bị từ chối: %v", err)
	}
	ok.Edges = append(ok.Edges, modeltest.Edge("b", model.PortSuccess, "a"))
	if err := validateFlow(ok); err == nil {
		t.Error("node Kết thúc không được có cạnh ra")
	}
}

func TestCheckControl(t *testing.T) {
	for want, s := range map[string]model.Step{
		"loại điều kiện": {Type: model.StepIf, Cond: &model.Cond{Kind: "xyz"}},
		"phép so":        {Type: model.StepIf, Cond: &model.Cond{Kind: model.CondAll, Conds: []model.Cond{{Kind: model.CondValue, Op: "~"}}}},
		"bị trùng":       {Type: model.StepSwitch, Cases: []model.Case{{Key: "a"}, {Key: "a"}}},
		"mã case":        {Type: model.StepSwitch, Cases: []model.Case{{Key: "a b"}}},
		"kết thúc":       {Type: model.StepEnd, EndStatus: "maybe"},
	} {
		if err := checkStep(&s); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi %q, có %v", want, err)
		}
	}
}

// Input {{key}} thay trong điều kiện / case trên bản sao – kịch bản gốc giữ nguyên.
func TestResolveInputsInConditions(t *testing.T) {
	sc := model.Script{Start: "i", Inputs: []model.ScriptInput{{Key: "ten", Type: "text", Default: "Lan"}}, Steps: []model.Step{
		modeltest.Node("i", "n1", model.StepIf, func(s *model.Step) {
			s.Cond = &model.Cond{Kind: model.CondAny, Conds: []model.Cond{{Kind: model.CondValue, Left: "{{ten}}", Op: "eq", Right: "x"}}}
		}),
		modeltest.Node("s", "n2", model.StepSwitchScreen, func(s *model.Step) { s.Cases = []model.Case{{Key: "a", Cond: modeltest.TextCond("chào {{ten}}")}} }),
	}}
	out, err := ResolveInputs(sc, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out.Steps[0].Cond.Conds[0].Left != "Lan" || out.Steps[1].Cases[0].Cond.FindText != "chào Lan" {
		t.Errorf("chưa thay: %+v / %+v", out.Steps[0].Cond, out.Steps[1].Cases[0].Cond)
	}
	if sc.Steps[0].Cond.Conds[0].Left != "{{ten}}" || sc.Steps[1].Cases[0].Cond.FindText != "chào {{ten}}" {
		t.Error("resolveInputs sửa vào kịch bản gốc")
	}
}
