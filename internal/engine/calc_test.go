package engine

import (
	"context"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

// Node Tính toán: biến số tính được, biến chữ thành chuỗi (không chèn được phép tính), sai kiểu → invalid_input.
func TestCalcNode(t *testing.T) {
	run := func(val, expr string) (*Engine, string, error) {
		sc := model.Script{Start: "v", Steps: []model.Step{
			modeltest.Node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "x", model.VarSet, val }),
			modeltest.Node("c", "n2", model.StepCalc, func(s *model.Step) { s.Value = expr }),
			modeltest.Node("t", "n3", model.StepText, func(s *model.Step) { s.Text = "{{n2.output.value}}" }),
		}, Edges: []model.Edge{modeltest.Edge("v", model.PortSuccess, "c"), modeltest.Edge("c", model.PortSuccess, "t")}}
		e, in := ctlEngine(sc, nil)
		err := e.Run(context.Background())
		return e, typed(in), err
	}
	for _, c := range []struct{ val, expr, want string }{
		{"1500", "{{vars.x}} * 2 + 1", "3001"},
		{"2.5", "round({{vars.x}})", "3"},
		{"10", "{{vars.x}} / 4", "2.5"},
		{"Hết hàng", "{{vars.x}} == \"Hết hàng\"", "true"},
		{"1 + 1", "{{vars.x}} == 2", "false"}, // chữ OCR "1 + 1" là chuỗi, không phải phép cộng
	} {
		if _, got, err := run(c.val, c.expr); err != nil || got != c.want {
			t.Errorf("x = %q, %s → %q %v, muốn %q", c.val, c.expr, got, err, c.want)
		}
	}
	e, _, err := run("1.234.000 đ", "{{vars.x}} * 2")
	if err == nil || e.results["c"].Error.Code != model.ErrInvalidInput {
		t.Errorf("nhân chữ phải lỗi invalid_input: %v", err)
	}
}
