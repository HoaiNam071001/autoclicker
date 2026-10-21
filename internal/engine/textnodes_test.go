package engine

import (
	"context"
	"testing"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

// Chuỗi xử lý chữ như mẫu xu-ly-chu: đọc giá (biến) → Số & ngày lấy số → Tính toán so ngưỡng; Tách chữ → Lặp qua danh sách.
func TestTextNodesChain(t *testing.T) {
	setX := func(v string) model.Step {
		return modeltest.Node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "x", model.VarSet, v })
	}
	sc := model.Script{Start: "v", Steps: []model.Step{
		setX("Giá: 1.234.000 đ"),
		modeltest.Node("n", "n2", model.StepTextNumber, func(s *model.Step) { s.Value = "{{vars.x}}" }),
		modeltest.Node("c", "n3", model.StepCalc, func(s *model.Step) { s.Value = "{{n2.output.value}} > 1000000" }),
		modeltest.Node("t", "n4", model.StepText, func(s *model.Step) { s.Text = "{{n3.output.value}}|" }),
		modeltest.Node("s", "n5", model.StepTextSplit, func(s *model.Step) { s.Value, s.Separator = "Táo, Cam,,Xoài", "," }),
		modeltest.Node("f", "n6", model.StepForEach, func(s *model.Step) { s.Value = "{{n5.output.items}}" }),
		modeltest.Node("i", "n7", model.StepText, func(s *model.Step) { s.Text = "{{n6.output.item}};" }),
	}, Edges: []model.Edge{
		modeltest.Edge("v", model.PortSuccess, "n"), modeltest.Edge("n", model.PortSuccess, "c"), modeltest.Edge("c", model.PortSuccess, "t"),
		modeltest.Edge("t", model.PortSuccess, "s"), modeltest.Edge("s", model.PortSuccess, "f"),
		modeltest.Edge("f", model.PortItem, "i"), modeltest.Edge("i", model.PortSuccess, "f"),
	}}
	e, in := ctlEngine(sc, nil)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := typed(in); got != "true|Táo;Cam;Xoài;" {
		t.Errorf("gõ ra %q", got)
	}

	// Không khớp mẫu → cổng fail (nối thì đi tiếp, không nối thì dừng kịch bản với not_found).
	sc = model.Script{Start: "v", Steps: []model.Step{
		setX("Hết hàng"),
		modeltest.Node("m", "n2", model.StepTextMatch, func(s *model.Step) { s.Value, s.Pattern = "{{vars.x}}", `\d+` }),
	}, Edges: []model.Edge{modeltest.Edge("v", model.PortSuccess, "m")}}
	e, _ = ctlEngine(sc, nil)
	if err := e.Run(context.Background()); err == nil || e.results["m"].Error.Code != model.ErrNotFound || e.results["m"].Port != model.PortFail {
		t.Errorf("không khớp phải fail not_found: %v %+v", err, e.results["m"])
	}
}

func TestRunTextNode(t *testing.T) {
	now := time.Date(2026, 10, 21, 9, 0, 0, 0, time.Local)
	run := func(s model.Step) map[string]any {
		t.Helper()
		_, out, err := RunTextNode(s, now)
		if err != nil {
			t.Fatalf("%s: %v", s.Type, err)
		}
		return out
	}
	if out := run(model.Step{Type: model.StepTextMatch, Value: "#12 và #34", Pattern: `#(\d+)`, MatchAll: true}); out["text"] != "12" || out["count"] != 2 {
		t.Errorf("Lọc mẫu: %v", out)
	}
	if out := run(model.Step{Type: model.StepTextClean, Value: " Đà  Nẵng ", StripAccents: true, Trim: true}); out["text"] != "Da Nang" {
		t.Errorf("Làm sạch: %v", out)
	}
	if out := run(model.Step{Type: model.StepTextNumber, NumOp: model.NumFormat, Value: "1.234.567,891", Decimals: 2}); out["text"] != "1.234.567,89" {
		t.Errorf("Định dạng số: %v", out)
	}
	if out := run(model.Step{Type: model.StepTextNumber, NumOp: model.NumDate, DateLayout: "dd/MM HH:mm"}); out["text"] != "21/10 09:00" || out["timestamp"] != now.Unix() {
		t.Errorf("Ngày bây giờ: %v", out)
	}
	if _, _, err := RunTextNode(model.Step{Type: model.StepTextNumber, Value: "Hết hàng"}, now); errCode(err) != model.ErrNotFound {
		t.Errorf("không có số phải not_found: %v", err)
	}
	if _, _, err := RunTextNode(model.Step{Type: model.StepTextMatch, Value: "x", Pattern: "("}, now); errCode(err) != model.ErrInvalidInput {
		t.Errorf("regex sai phải invalid_input: %v", err)
	}
}
