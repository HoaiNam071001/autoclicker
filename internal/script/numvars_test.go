package script

import (
	"strings"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

// numVarsFlow: Đọc phần tử web (output count) → Cuộn với numVars nv.
func numVarsFlow(nv ...model.NumVar) model.Script {
	return model.Script{Loops: 1, Start: "r", Steps: []model.Step{
		modeltest.Node("r", "n1", model.StepWebRead, func(s *model.Step) {
			s.Selectors, s.TimeoutMs, s.IntervalMs = []string{"css=li"}, 1000, 100
		}),
		modeltest.Node("s", "n2", model.StepScroll, func(s *model.Step) { s.Direction, s.NumVars = "down", nv }),
	}, Edges: []model.Edge{modeltest.Edge("r", model.PortSuccess, "s")}}
}

// Trường số lấy từ biến (F5c): giá trị tĩnh không cần hợp lệ (amount = 0), biến phải là output số có thật.
func TestValidateNumVars(t *testing.T) {
	for _, expr := range []string{"{{n1.output.count}} + 1", "{{n1.output.durationMs}}", "{{vars.n}}", "{{so_nac}}", "3"} {
		if err := ValidateScript(numVarsFlow(model.NumVar{Field: "amount", Expr: expr}), model.Rect{}); err != nil {
			t.Errorf("%q: %v", expr, err)
		}
	}
	bad := map[string]string{
		`không có output "x" – có: text, count`: "{{n1.output.x}}",
		"không phải số – ô số dùng":             "{{n1.status}}",
		"sai dạng":                              "{{n1.output.count}} *",
		"chưa nhập giá trị":                     " ",
	}
	for want, expr := range bad {
		if err := ValidateScript(numVarsFlow(model.NumVar{Field: "amount", Expr: expr}), model.Rect{}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: muốn lỗi %q, có %v", expr, want, err)
		}
	}
	// ghim input + biến cho cùng trường
	sc := numVarsFlow(model.NumVar{Field: "amount", Expr: "2"})
	sc.Inputs = []model.ScriptInput{{Key: "nac", Type: model.InputNumber, Targets: []model.InputTarget{{StepID: "s", Field: "amount"}}}}
	if err := ValidateScript(sc, model.Rect{}); err == nil || !strings.Contains(err.Error(), "vừa ghim input nac") {
		t.Errorf("ghim + biến: %v", err)
	}
}

// Mở file: trường không nhận biến với loại node / lặp lại → lỗi.
func TestCheckNumVars(t *testing.T) {
	for want, nv := range map[string][]model.NumVar{
		`"holdMs" của node scroll không nhận biến`: {{Field: "holdMs", Expr: "1"}},
		`"amount" lặp lại`:                         {{Field: "amount", Expr: "1"}, {Field: "amount", Expr: "2"}},
	} {
		s := numVarsFlow(nv...).Steps[1]
		if err := checkStep(&s); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn %q, có %v", want, err)
		}
	}
}

func TestApplyNumVars(t *testing.T) {
	s := model.Step{Type: model.StepClick, NumVars: []model.NumVar{{Field: "clicks", Expr: "2 + 1"}, {Field: "delayAfter", Expr: "1500.4"}}}
	if err := ApplyNumVars(&s); err != nil || s.Clicks != 3 || s.DelayAfter != 1500 {
		t.Errorf("%v %+v", err, s)
	}
	s.NumVars = []model.NumVar{{Field: "timeoutMs", Expr: "0"}}
	if err := ApplyNumVars(&s); err == nil || !strings.Contains(err.Error(), "timeoutMs: thời gian chờ tối đa phải > 0") {
		t.Errorf("%v", err)
	}
}
