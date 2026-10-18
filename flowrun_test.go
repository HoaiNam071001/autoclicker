package main

import (
	"context"
	"image"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

// node gắn id + ref cố định cho test flow
func node(id, ref, typ string, mod func(*model.Step)) model.Step {
	return step(typ, func(s *model.Step) {
		s.ID, s.Ref = id, ref
		if mod != nil {
			mod(s)
		}
	})
}

func edge(from, port, to string) model.Edge {
	return model.Edge{ID: from + "." + port, From: from, Port: port, To: to}
}

var (
	noText  = []vision.TextLine{line("Đang tải", 0, 0, 80, 20)}
	sendBtn = []vision.TextLine{line("Gửi", 100, 100, 40, 20)}
)

// ocrAfter – bộ đọc chữ giả: thấy nút "Gửi" khi ready() đúng (vd sau N lần cuộn), không thì chỉ thấy "Đang tải".
type ocrAfter func() bool

func (f ocrAfter) ReadText(image.Image) ([]vision.TextLine, error) {
	if f() {
		return sendBtn, nil
	}
	return noText, nil
}

// Flow: chờ chữ "Gửi" (w) → click chỗ thấy (c) → gõ (t); không thấy → cuộn (s) → chờ lại.
// Nút hiện sau `scrolls` lần cuộn (-1 = không bao giờ); ready (nếu có) thay điều kiện đó.
func sendFlow(scrolls int, ready func() bool) (*engine, *fakeInput) {
	sc := model.Script{Loops: 1, Region: &model.Rect{X: 100, Y: 50, W: 400, H: 300}, Start: "w",
		Steps: []model.Step{
			node("w", "n1", model.StepWaitImage, func(s *model.Step) { s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = model.FindByText, "Gửi", 1, 1 }),
			node("c", "n2", model.StepClickImage, func(s *model.Step) { s.ImageFrom = "w" }),
			node("t", "n3", model.StepText, func(s *model.Step) { s.Text = "ở {{n1.output.x}},{{n1.output.y}}" }),
			node("s", "n4", model.StepScroll, func(s *model.Step) { s.Amount = 1 }),
		},
		Edges: []model.Edge{edge("w", model.PortSuccess, "c"), edge("w", model.PortFail, "s"), edge("s", model.PortSuccess, "w"), edge("c", model.PortSuccess, "t")},
	}
	in := &fakeInput{}
	if ready == nil {
		ready = func() bool {
			n := 0
			for _, l := range in.lines() {
				if strings.HasPrefix(l, "scroll") {
					n++
				}
			}
			return scrolls >= 0 && n >= scrolls
		}
	}
	e := newEngine(in, sc, nil)
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}}, ocrAfter(ready)
	return e, in
}

func TestFlowFailPortLoopsBack(t *testing.T) {
	e, in := sendFlow(2, nil)
	var got []string
	e.onNode = func(r model.StepResult) { got = append(got, r.Ref+":"+r.Status+">"+r.Port) }
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	// chữ "Gửi" ở px (100, 100) 40×20 → tâm 100 + 120, 50 + 110
	want := []string{"scroll down", "scroll down", "move 220,160", "click left false", "type ở", "type  ", "type 2", "type 2", "type 0", "type ,", "type 1", "type 6", "type 0"}
	if !reflect.DeepEqual(in.lines(), want) {
		t.Errorf("lệnh %q\nmuốn %q", in.lines(), want)
	}
	if w := []string{"n1:fail>fail", "n4:success>success", "n1:fail>fail", "n4:success>success", "n1:success>success", "n2:success>success", "n3:success>success"}; !reflect.DeepEqual(got, w) {
		t.Errorf("kết quả %q", got)
	}
	if r := e.results["w"]; r.Output["x"] != 220 || r.Output["text"] != "Gửi" {
		t.Errorf("output chờ chữ: %+v", r.Output)
	}
}

func TestFlowFailUnconnectedStops(t *testing.T) {
	e, in := sendFlow(-1, nil)
	e.sc.Edges = e.sc.Edges[:1] // bỏ cạnh fail
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "n1: không thấy chữ") || e.failed != "w" || len(in.lines()) != 0 {
		t.Fatalf("lỗi %v, failed %q, lệnh %q", err, e.failed, in.lines())
	}
	if r := e.results["w"]; r.Error == nil || r.Error.Code != model.ErrNotFound {
		t.Errorf("mã lỗi: %+v", r.Error)
	}
}

func TestFlowRetry(t *testing.T) {
	tries := 0 // thấy ở lần thử thứ 3 (mỗi lần thử: chờ đọc nhiều lần tới hết giờ)
	e, _ := sendFlow(0, func() bool { return tries >= 3 })
	e.onProgress = func(p Progress) {
		if p.NodeID == "w" {
			tries = p.Retry + 1
		}
	}
	e.sc.Edges = e.sc.Edges[:1]
	e.sc.Steps[0].Retry = &model.Retry{Count: 2}
	e.sc.Steps[1].Enabled = false // click tắt → bỏ qua, đi tiếp
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	r := e.results["w"]
	if r.Attempts != 3 || r.Status != model.StatusSuccess || len(r.Logs) != 2 || !strings.Contains(r.Logs[0].Msg, "lần 1 / 3") {
		t.Errorf("thử lại: %+v", r)
	}
	if e.results["c"].Status != model.StatusSkipped {
		t.Errorf("node tắt: %+v", e.results["c"])
	}
}

func TestFlowVarFromNodeNotRun(t *testing.T) {
	sc := model.Script{Loops: 1, Start: "a", Steps: []model.Step{
		node("a", "n1", model.StepText, func(s *model.Step) { s.Text = "{{n2.output.text}}" }),
		node("b", "n2", model.StepText, func(s *model.Step) { s.Text = "x" }),
	}, Edges: []model.Edge{edge("a", model.PortSuccess, "b")}}
	e := newEngine(&fakeInput{}, sc, nil)
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "n2 chưa chạy") || e.results["a"].Error.Code != model.ErrInvalidInput {
		t.Fatalf("lỗi %v", err)
	}
}

func TestFlowClickFromFailedNode(t *testing.T) {
	// chờ chữ thất bại → đi cổng fail sang click lấy chỗ thấy của chính node đó → không có chỗ thấy → dừng
	e, in := sendFlow(-1, nil)
	e.sc.Edges = []model.Edge{edge("w", model.PortFail, "c")}
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "n2: node n1 chưa tìm thấy") || len(in.lines()) != 0 {
		t.Fatalf("lỗi %v, lệnh %q", err, in.lines())
	}
}

func TestFlowInfiniteLoopLimit(t *testing.T) {
	sc := model.Script{Loops: 1, Start: "a", Steps: []model.Step{
		node("a", "n1", model.StepWait, nil), node("b", "n2", model.StepWait, nil),
	}, Edges: []model.Edge{edge("a", model.PortSuccess, "b"), edge("b", model.PortSuccess, "a")}}
	sc.Steps[0].DelayAfter, sc.Steps[1].DelayAfter = 0, 0
	err := newEngine(&fakeInput{}, sc, nil).run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "lặp vô hạn") || !strings.Contains(err.Error(), "n1 → n2") {
		t.Fatalf("lỗi %v", err)
	}
}

func TestValidateFlow(t *testing.T) {
	base := func() model.Script {
		return model.Script{Start: "a", Steps: []model.Step{node("a", "n1", model.StepWait, nil), node("b", "n2", model.StepWait, nil)},
			Edges: []model.Edge{edge("a", model.PortSuccess, "b")}}
	}
	cases := map[string]func(*model.Script){
		"chưa chọn node bắt đầu": func(s *model.Script) { s.Start = "x" },
		"trùng với node khác":    func(s *model.Script) { s.Steps[1].Ref = "n1" },
		"không hợp lệ":           func(s *model.Script) { s.Steps[1].Ref = "vars" },
		"không có cổng":          func(s *model.Script) { s.Edges[0].Port = "true" },
		"nối tới node không còn": func(s *model.Script) { s.Edges[0].To = "x" },
		"nối vào chính node này": func(s *model.Script) { s.Edges[0].To = "a" },
		"có hơn một cạnh": func(s *model.Script) {
			s.Edges = append(s.Edges, model.Edge{ID: "z", From: "a", Port: model.PortSuccess, To: "b"})
		},
		"không có trong kịch bản": func(s *model.Script) {
			s.Steps[1] = node("b", "n2", model.StepText, func(x *model.Step) { x.Text = "{{n9.output.x}}" })
		},
		"chưa chọn node lấy chỗ": func(s *model.Script) { s.Steps[1] = node("b", "n2", model.StepClickImage, nil) },
		"không tìm ảnh / chữ": func(s *model.Script) {
			s.Steps[1] = node("b", "n2", model.StepClickImage, func(x *model.Step) { x.ImageFrom = "a" })
		},
		"kịch bản chưa có node nào": func(s *model.Script) { s.Steps = nil },
	}
	for want, mod := range cases {
		sc := base()
		mod(&sc)
		if err := validateFlow(sc); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi %q, có %v", want, err)
		}
	}
	if err := validateFlow(base()); err != nil {
		t.Errorf("flow hợp lệ báo lỗi: %v", err)
	}
	// biến vars / last không cần node
	sc := base()
	sc.Steps[1] = node("b", "n2", model.StepText, func(x *model.Step) { x.Text = "{{last.status}} {{vars.dem}}" })
	if err := validateFlow(sc); err != nil {
		t.Errorf("vars / last: %v", err)
	}
}

func TestUnreachableNodes(t *testing.T) {
	sc := model.Script{Start: "a", Steps: []model.Step{{ID: "a"}, {ID: "b"}, {ID: "c"}}, Edges: []model.Edge{edge("a", model.PortFail, "c")}}
	if got := unreachableNodes(sc); !reflect.DeepEqual(got, []string{"b"}) {
		t.Errorf("%q", got)
	}
}
