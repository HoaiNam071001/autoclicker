package main

import (
	"context"
	"image"
	"image/color"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/device"
	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

// Test node điều khiển (bước 20 F3): Nếu, Switch, Switch theo màn hình, Đếm, Kết thúc, Đặt biến.

// typed – chữ đã gõ (fakeInput ghi từng ký tự).
func typed(in *devicetest.FakeInput) string {
	var b strings.Builder
	for _, l := range in.Lines() {
		if t, ok := strings.CutPrefix(l, "type "); ok {
			b.WriteString(t)
		}
	}
	return b.String()
}

// ocrSeq – bộ đọc chữ giả: lần đọc thứ i trả frames[i] (lần cuối lặp lại); đếm số lần đọc.
type ocrSeq struct {
	frames [][]vision.TextLine
	reads  int
}

func (o *ocrSeq) ReadText(image.Image) ([]vision.TextLine, error) {
	o.reads++
	return o.frames[min(o.reads-1, len(o.frames)-1)], nil
}

func ctlEngine(sc model.Script, ocr device.TextReader) (*engine, *devicetest.FakeInput) {
	if sc.Loops == 0 {
		sc.Loops = 1
	}
	if sc.Region == nil {
		sc.Region = &model.Rect{X: 100, Y: 50, W: 400, H: 300}
	}
	in := &devicetest.FakeInput{}
	e := newEngine(in, sc, nil)
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}}, ocr
	return e, in
}

func textCond(t string) *model.Cond { return &model.Cond{Kind: model.CondText, FindText: t} }

func trace(e *engine) *[]string {
	var got []string
	e.onNode = func(r model.StepResult) { got = append(got, r.Ref+">"+r.Port) }
	return &got
}

// Ví dụ Facebook: không thấy nút → Đếm → cuộn → tìm lại; hết lượt → Kết thúc báo lỗi kèm số lần.
func TestCounterEndFail(t *testing.T) {
	sc := model.Script{Start: "w", Steps: []model.Step{
		node("w", "n1", model.StepWaitImage, func(s *model.Step) { s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = model.FindByText, "Gửi", 1, 1 }),
		node("s", "n2", model.StepScroll, func(s *model.Step) { s.Amount = 1 }),
		node("k", "n3", model.StepCounter, func(s *model.Step) { s.Max = 2 }),
		node("x", "n4", model.StepEnd, func(s *model.Step) {
			s.EndStatus, s.Message = model.EndFail, "không thấy sau {{n3.output.count}} lần"
		}),
	}, Edges: []model.Edge{edge("w", model.PortFail, "k"), edge("k", model.PortNext, "s"), edge("k", model.PortDone, "x"), edge("s", model.PortSuccess, "w")}}
	e, in := ctlEngine(sc, &ocrSeq{frames: [][]vision.TextLine{noText}})
	got := trace(e)
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "n4: không thấy sau 3 lần") || e.failed != "x" {
		t.Fatalf("lỗi %v, failed %q", err, e.failed)
	}
	if n := strings.Count(strings.Join(in.Lines(), "|"), "scroll"); n != 2 {
		t.Errorf("cuộn %d lần, muốn 2", n)
	}
	want := []string{"n1>fail", "n3>next", "n2>success", "n1>fail", "n3>next", "n2>success", "n1>fail", "n3>done", "n4>fail"}
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("đường đi %q", *got)
	}
	if r := e.results["x"]; r.Error == nil || r.Error.Code != model.ErrEnd {
		t.Errorf("mã lỗi kết thúc: %+v", r.Error)
	}
}

// Kết thúc thành công: dừng vòng, vòng sau chạy lại từ đầu; Đếm đếm lại mỗi vòng.
func TestEndSuccessNextLoop(t *testing.T) {
	sc := model.Script{Loops: 2, Start: "k", Steps: []model.Step{
		node("k", "n1", model.StepCounter, func(s *model.Step) { s.Max = 5 }),
		node("t", "n2", model.StepText, func(s *model.Step) { s.Text = "{{n1.output.count}}" }),
		node("x", "n3", model.StepEnd, func(s *model.Step) { s.EndStatus = model.EndSuccess }),
		node("z", "n4", model.StepText, func(s *model.Step) { s.Text = "không tới" }),
	}, Edges: []model.Edge{edge("k", model.PortNext, "t"), edge("t", model.PortSuccess, "x")}}
	e, in := ctlEngine(sc, &ocrSeq{frames: [][]vision.TextLine{noText}})
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := typed(in); got != "11" {
		t.Errorf("gõ %q, muốn 11 (đếm lại mỗi vòng)", got)
	}
}

// Nếu: thấy chữ → true, output chỗ thấy để Click theo ảnh dùng; không thấy (xem một lần) → false.
func TestIfTextAndClickFromIf(t *testing.T) {
	flow := func(lines []vision.TextLine) (*engine, *devicetest.FakeInput) {
		sc := model.Script{Start: "i", Steps: []model.Step{
			node("i", "n1", model.StepIf, func(s *model.Step) { s.Cond = textCond("Gửi") }),
			node("c", "n2", model.StepClickImage, func(s *model.Step) { s.ImageFrom = "i" }),
			node("f", "n3", model.StepText, func(s *model.Step) { s.Text = "không" }),
		}, Edges: []model.Edge{edge("i", model.PortTrue, "c"), edge("i", model.PortFalse, "f")}}
		return ctlEngine(sc, &ocrSeq{frames: [][]vision.TextLine{lines}})
	}
	e, in := flow(sendBtn)
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"move 220,160", "click left false"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("lệnh %q", in.Lines())
	}
	if r := e.results["i"]; r.Output["result"] != true || r.Output["text"] != "Gửi" || len(r.Logs) == 0 {
		t.Errorf("output Nếu: %+v, log %v", r.Output, r.Logs)
	}
	e, in = flow(noText)
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if typed(in) != "không" || e.results["i"].Port != model.PortFalse {
		t.Errorf("nhánh sai: %q, cổng %q", typed(in), e.results["i"].Port)
	}
}

// Nếu có thời gian chờ: xem lại tới khi đúng; đảo + nhóm và / hoặc + so giá trị biến node.
func TestIfWaitGroupsAndValues(t *testing.T) {
	ocr := &ocrSeq{frames: [][]vision.TextLine{noText, noText, sendBtn}}
	sc := model.Script{Start: "v", Steps: []model.Step{
		node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "lan", model.VarAdd, "2,5" }),
		node("i", "n2", model.StepIf, func(s *model.Step) {
			s.TimeoutMs, s.IntervalMs = 20, 1
			s.Cond = &model.Cond{Kind: model.CondAll, Conds: []model.Cond{
				*textCond("Gửi"),
				{Kind: model.CondValue, Left: "{{vars.lan}}", Op: "gt", Right: "2"},
				{Kind: model.CondAny, Not: true, Conds: []model.Cond{{Kind: model.CondValue, Left: "{{n1.output.name}}", Op: "eq", Right: "LAN"}}},
			}}
		}),
		node("t", "n3", model.StepText, func(s *model.Step) { s.Text = "{{n2.output.result}}" }),
	}, Edges: []model.Edge{edge("v", model.PortSuccess, "i"), edge("i", model.PortTrue, "t"), edge("i", model.PortFalse, "t")}}
	e, in := ctlEngine(sc, ocr)
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	// điều kiện 3: tên biến "lan" = "LAN" (không phân biệt hoa thường) rồi đảo → sai → cả nhóm sai dù đã thấy chữ
	if typed(in) != "false" {
		t.Errorf("gõ %q", typed(in))
	}
	if ocr.reads < 3 {
		t.Errorf("Nếu có chờ phải xem lại tới khi hết giờ / thấy, đọc %d lần", ocr.reads)
	}
}

// Nếu theo ảnh mẫu / màu pixel – cùng cách tìm với Chờ ảnh / Chờ pixel.
func TestIfImageAndPixel(t *testing.T) {
	with, tpl := screenWithButton(1, true)
	sc := model.Script{Region: &model.Rect{X: 0, Y: 0, W: 400, H: 300}, Start: "i", Steps: []model.Step{
		node("i", "n1", model.StepIf, func(s *model.Step) {
			s.Cond = &model.Cond{Kind: model.CondAll, Conds: []model.Cond{
				{Kind: model.CondImage, Template: tpl, TemplateScale: 1, Threshold: 90},
				{Kind: model.CondPixel, Points: []model.Point{{X: 10, Y: 10}}, Color: "#1E73F0", Tolerance: 5},
			}}
		}),
	}}
	e, _ := ctlEngine(sc, nil)
	e.scr = &fakeScreen{frames: []image.Image{with}, colors: []color.RGBA{{30, 115, 240, 255}}}
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if r := e.results["i"]; r.Port != model.PortTrue || r.Output["x"] != 280 || r.Output["y"] != 130 {
		t.Errorf("cổng %q, output %v", r.Port, r.Output)
	}
}

// Switch theo giá trị: case đầu khớp thắng; regex; không khớp → default.
func TestSwitchValue(t *testing.T) {
	for value, want := range map[string]string{"Đã GỬI xong": "case:ok", "Error 500": "case:err", "gì đó": model.PortDefault} {
		sc := model.Script{Start: "s", Steps: []model.Step{
			node("s", "n1", model.StepSwitch, func(s *model.Step) {
				s.Value = value
				s.Cases = []model.Case{{Key: "ok", Op: "contains", Match: "đã gửi"}, {Key: "err", Op: "regex", Match: `(?i)error \d+`}, {Key: "x", Op: "eq", Match: "Error 500"}}
			}),
			node("t", "n2", model.StepText, func(s *model.Step) { s.Text = "{{n1.output.case}}" }),
		}, Edges: []model.Edge{edge("s", "case:ok", "t"), edge("s", "case:err", "t"), edge("s", model.PortDefault, "t")}}
		e, in := ctlEngine(sc, nil)
		if err := e.run(context.Background()); err != nil {
			t.Fatal(err)
		}
		if r := e.results["s"]; r.Port != want || typed(in) != strings.TrimPrefix(want, "case:") {
			t.Errorf("%q → cổng %q, gõ %q; muốn %q", value, r.Port, typed(in), want)
		}
	}
}

// Switch theo màn hình: case nào hiện trước thắng; mỗi lần xét chỉ chụp + đọc chữ một lần cho mọi case cùng vùng;
// hết giờ → none.
func TestSwitchScreen(t *testing.T) {
	cases := []model.Case{{Key: "login", Cond: textCond("Đăng nhập")}, {Key: "send", Cond: textCond("Gửi")}}
	ocr := &ocrSeq{frames: [][]vision.TextLine{noText, sendBtn}}
	sc := model.Script{Start: "s", Steps: []model.Step{
		node("s", "n1", model.StepSwitchScreen, func(s *model.Step) { s.Cases, s.TimeoutMs, s.IntervalMs = cases, 1000, 1 }),
		node("c", "n2", model.StepClickImage, func(s *model.Step) { s.ImageFrom = "s" }),
	}, Edges: []model.Edge{edge("s", "case:send", "c")}}
	e, in := ctlEngine(sc, ocr)
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	scr := e.scr.(*fakeScreen)
	if r := e.results["s"]; r.Port != "case:send" || r.Output["case"] != "send" || len(scr.shots) != 2 || ocr.reads != 2 {
		t.Errorf("cổng %q, output %v, chụp %d, đọc %d", r.Port, r.Output, len(scr.shots), ocr.reads)
	}
	if want := []string{"move 220,160", "click left false"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("lệnh %q", in.Lines())
	}

	e, _ = ctlEngine(model.Script{Start: "s", Steps: []model.Step{
		node("s", "n1", model.StepSwitchScreen, func(s *model.Step) { s.Cases, s.TimeoutMs, s.IntervalMs = cases, 5, 1 }),
	}}, &ocrSeq{frames: [][]vision.TextLine{noText}})
	if err := e.run(context.Background()); err != nil || e.results["s"].Port != model.PortNone {
		t.Errorf("hết giờ: %v, cổng %q", err, e.results["s"].Port)
	}
}

// Đặt biến cộng dồn qua vòng lặp; cộng vào chữ → lỗi invalid_input.
func TestSetVarAdd(t *testing.T) {
	sc := model.Script{Start: "v", Steps: []model.Step{
		node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "n", model.VarAdd, "1" }),
		node("k", "n2", model.StepCounter, func(s *model.Step) { s.Max = 2 }),
		node("t", "n3", model.StepText, func(s *model.Step) { s.Text = "{{vars.n}}" }),
	}, Edges: []model.Edge{edge("v", model.PortSuccess, "k"), edge("k", model.PortNext, "v"), edge("k", model.PortDone, "t")}}
	e, in := ctlEngine(sc, nil)
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if typed(in) != "3" {
		t.Errorf("gõ %q, muốn 3", typed(in))
	}
	e, _ = ctlEngine(model.Script{Start: "a", Steps: []model.Step{
		node("a", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "n", model.VarSet, "chữ" }),
		node("b", "n2", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "n", model.VarAdd, "1" }),
	}, Edges: []model.Edge{edge("a", model.PortSuccess, "b")}}, nil)
	if err := e.run(context.Background()); err == nil || e.results["b"].Error.Code != model.ErrInvalidInput {
		t.Errorf("cộng vào chữ phải lỗi: %v", err)
	}
}

func TestValidateControl(t *testing.T) {
	bad := map[string]model.Step{
		"chưa nhập chữ cần thấy": node("a", "n1", model.StepIf, nil),
		"mẫu sai":                node("a", "n1", model.StepSwitch, func(s *model.Step) { s.Cases = []model.Case{{Key: "a", Op: "regex", Match: "("}} }),
		"chưa có case":           node("a", "n1", model.StepSwitchScreen, func(s *model.Step) { s.Cases = nil }),
		"cần điều kiện thấy":     node("a", "n1", model.StepSwitchScreen, func(s *model.Step) { s.Cases = []model.Case{{Key: "a"}} }),
		"tối đa phải ≥ 1":        node("a", "n1", model.StepCounter, func(s *model.Step) { s.Max = 0 }),
		"tên biến":               node("a", "n1", model.StepSetVar, func(s *model.Step) { s.VarName = "1x" }),
		"cộng cần một số":        node("a", "n1", model.StepSetVar, func(s *model.Step) { s.VarOp, s.Value = model.VarAdd, "abc" }),
		"không có cổng":          node("a", "n1", model.StepCounter, nil),
	}
	for want, s := range bad {
		sc := model.Script{Start: "a", Steps: []model.Step{s, node("b", "n2", model.StepWait, nil)}}
		if want == "không có cổng" {
			sc.Edges = []model.Edge{edge("a", model.PortSuccess, "b")}
		}
		err := validateScript(sc, model.Rect{W: 1920, H: 1080})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi %q, có %v", want, err)
		}
	}
	// cổng case theo key; node Kết thúc không có cổng
	ok := model.Script{Start: "a", Steps: []model.Step{
		node("a", "n1", model.StepSwitch, nil), node("b", "n2", model.StepEnd, nil),
	}, Edges: []model.Edge{edge("a", "case:c2", "b"), edge("a", model.PortDefault, "b")}}
	if err := validateScript(ok, model.Rect{W: 1920, H: 1080}); err != nil {
		t.Errorf("flow hợp lệ bị từ chối: %v", err)
	}
	ok.Edges = append(ok.Edges, edge("b", model.PortSuccess, "a"))
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
		node("i", "n1", model.StepIf, func(s *model.Step) {
			s.Cond = &model.Cond{Kind: model.CondAny, Conds: []model.Cond{{Kind: model.CondValue, Left: "{{ten}}", Op: "eq", Right: "x"}}}
		}),
		node("s", "n2", model.StepSwitchScreen, func(s *model.Step) { s.Cases = []model.Case{{Key: "a", Cond: textCond("chào {{ten}}")}} }),
	}}
	out, err := resolveInputs(sc, nil)
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

// Thử node: Nếu xét trên màn hình lúc này; Switch có biến node → chưa biết.
func TestTestNode(t *testing.T) {
	sc := model.Script{Region: &model.Rect{X: 100, Y: 50, W: 400, H: 300}, Steps: []model.Step{
		node("i", "n1", model.StepIf, func(s *model.Step) {
			s.Cond = &model.Cond{Kind: model.CondAll, Conds: []model.Cond{*textCond("Gửi"), *textCond("Hủy")}}
		}),
		node("s", "n2", model.StepSwitch, func(s *model.Step) { s.Value = "{{n1.output.text}}" }),
	}}
	e, _ := ctlEngine(sc, &ocrSeq{frames: [][]vision.TextLine{sendBtn}})
	r, err := e.testNode(sc.Steps[0])
	if err != nil || r.Port != model.PortFalse || len(r.Parts) != 2 || !r.Parts[0].Ok || r.Parts[1].Ok || r.Parts[0].Hit == nil {
		t.Errorf("Nếu: %+v %v", r, err)
	}
	r, err = e.testNode(sc.Steps[1])
	if err != nil || r.Known || r.Port != "" || len(r.Parts) != 2 {
		t.Errorf("Switch biến node: %+v %v", r, err)
	}
}
