package main

import (
	"context"
	"encoding/json"
	"image"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func typeStep(id, text string) Step {
	return step(StepText, func(s *Step) { s.ID, s.Text = id, text })
}

func textCond(text string) *Cond { return &Cond{Kind: CondText, FindText: text, IntervalMs: 1} }

func ifStep(id string, c *Cond, then, els []Step) Step {
	return step(StepIf, func(s *Step) { s.ID, s.Cond, s.Children, s.Else = id, c, then, els })
}

func groupStep(id string, mod func(*Step), children ...Step) Step {
	return step(StepGroup, func(s *Step) {
		s.ID, s.Children = id, children
		if mod != nil {
			mod(s)
		}
	})
}

// runBlocks chạy các bước với bộ đọc chữ giả (mỗi lần xét điều kiện chữ đọc 1 lần), trả lệnh đã gõ.
func runBlocks(t *testing.T, sc Script, reads ...[]TextLine) ([]string, error) {
	t.Helper()
	if sc.Loops == 0 {
		sc.Loops = 1
	}
	if sc.Region == nil {
		sc.Region = &Rect{W: 400, H: 300}
	}
	if reads == nil {
		reads = [][]TextLine{noText}
	}
	in := &fakeInput{}
	e := newEngine(in, sc, nil)
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}}, &fakeOCR{reads: reads}
	err := e.run(context.Background())
	var typed []string
	for _, l := range in.lines() {
		if strings.HasPrefix(l, "type ") {
			typed = append(typed, strings.TrimPrefix(l, "type "))
		}
	}
	return typed, err
}

func expect(t *testing.T, name string, got []string, err error, want ...string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: lỗi %v", name, err)
	} else if strings.Join(got, "") != strings.Join(want, "") { // chuột giả ghi từng ký tự gõ
		t.Errorf("%s: gõ %q, muốn %q", name, got, want)
	}
}

func TestIfElse(t *testing.T) {
	sc := Script{Steps: []Step{
		ifStep("if", textCond("Gửi"), []Step{typeStep("t1", "đúng")}, []Step{typeStep("e1", "sai")}),
		typeStep("after", "sau"),
	}}
	got, err := runBlocks(t, sc, sendBtn)
	expect(t, "thấy", got, err, "đúng", "sau")
	got, err = runBlocks(t, sc, noText)
	expect(t, "không thấy", got, err, "sai", "sau")
	sc.Steps[0].Cond.Not = true
	got, err = runBlocks(t, sc, sendBtn)
	expect(t, "đảo", got, err, "sai", "sau")
}

func TestGroupLoops(t *testing.T) {
	// lặp 3 lần
	sc := Script{Steps: []Step{groupStep("g", func(s *Step) { s.Repeat = 3 }, typeStep("a", "a"))}}
	got, err := runBlocks(t, sc)
	expect(t, "repeat", got, err, "a", "a", "a")

	// lặp đến khi thấy chữ: 2 lần không thấy → thân chạy 2 lần
	sc = Script{Steps: []Step{groupStep("g", func(s *Step) { s.LoopMode, s.Cond = LoopUntil, textCond("Gửi") }, typeStep("a", "cuộn"))}}
	got, err = runBlocks(t, sc, noText, noText, sendBtn)
	expect(t, "until", got, err, "cuộn", "cuộn")

	// lặp khi còn thấy
	sc.Steps[0].LoopMode = LoopWhile
	got, err = runBlocks(t, sc, sendBtn, noText)
	expect(t, "while", got, err, "cuộn")

	// quá maxLoops → bước thất bại; có onFail next thì đi tiếp
	sc = Script{Steps: []Step{
		groupStep("g", func(s *Step) { s.LoopMode, s.Cond, s.MaxLoops = LoopUntil, textCond("Gửi"), 2 }, typeStep("a", "cuộn")),
		typeStep("b", "tiếp"),
	}}
	_, err = runBlocks(t, sc)
	if err == nil || !strings.Contains(err.Error(), "bước 1: đã lặp 2 lần") {
		t.Errorf("maxLoops: %v", err)
	}
	sc.Steps[0].OnFail = &OnFail{Action: FailNext}
	got, err = runBlocks(t, sc)
	expect(t, "maxLoops + bỏ qua", got, err, "cuộn", "cuộn", "tiếp")
}

func TestBreakAndStop(t *testing.T) {
	// lặp 5 lần, vòng thứ 2 thấy chữ → thoát; khối lặp ngoài vẫn chạy tiếp
	inner := groupStep("in", func(s *Step) { s.Repeat = 5 },
		typeStep("a", "a"),
		ifStep("if", textCond("Gửi"), []Step{step(StepBreak, nil)}, nil),
	)
	sc := Script{Steps: []Step{groupStep("out", func(s *Step) { s.Repeat = 2 }, inner, typeStep("b", "b"))}}
	got, err := runBlocks(t, sc, noText, sendBtn, sendBtn)
	expect(t, "break", got, err, "a", "a", "b", "a", "b")

	// Dừng bình thường: không chạy bước sau, không chạy vòng sau
	sc = Script{Loops: 3, Steps: []Step{typeStep("a", "a"), groupStep("g", nil, step(StepStop, nil)), typeStep("b", "b")}}
	got, err = runBlocks(t, sc)
	expect(t, "stop", got, err, "a")
	// Dừng báo lỗi: kèm lời nhắn và đường dẫn bước
	sc.Steps[1].Children[0].StopError, sc.Steps[1].Children[0].Message = true, "Không có người nhận"
	_, err = runBlocks(t, sc)
	if err == nil || err.Error() != "vòng 1, bước 2.1: Không có người nhận" {
		t.Errorf("stop lỗi: %v", err)
	}
}

func TestOnFailElse(t *testing.T) {
	wait := step(StepWaitImage, func(s *Step) {
		s.ID, s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = "w", FindByText, "Gửi", 1, 1
		s.OnFail, s.Else = &OnFail{Retries: 1, Action: FailElse}, []Step{typeStep("e", "cuộn")}
	})
	sc := Script{Steps: []Step{wait, typeStep("b", "tiếp")}}
	got, err := runBlocks(t, sc)
	expect(t, "else", got, err, "cuộn", "tiếp")
	got, err = runBlocks(t, sc, sendBtn)
	expect(t, "thấy thì không chạy else", got, err, "tiếp")
}

func TestCondInputAndResult(t *testing.T) {
	sc := Script{
		Inputs: []ScriptInput{{Key: "ten", Label: "Tên", Type: InputText}},
		Steps: []Step{
			ifStep("if", &Cond{Kind: CondInput, Left: "{{ten}}", Op: "eq", Right: "lan"}, []Step{typeStep("a", "chào Lan")}, []Step{typeStep("b", "chào bạn")}),
		},
	}
	for v, want := range map[string]string{"Lan": "chào Lan", "Hùng": "chào bạn"} {
		out, err := resolveInputs(sc, map[string]any{"ten": v})
		if err != nil {
			t.Fatal(err)
		}
		got, err := runBlocks(t, out)
		expect(t, "input "+v, got, err, want)
	}
	if sc.Steps[0].Cond.Left != "{{ten}}" {
		t.Error("resolveInputs sửa vào kịch bản gốc")
	}
	for _, c := range []struct {
		l, op, r string
		want     bool
	}{{"10", "gt", "9", true}, {"2,5", "lt", "3", true}, {"Xin chào", "contains", "CHÀO", true}, {" ", "empty", "", true}, {"a", "ne", "A", false}} {
		if got, err := compareValues(c.l, c.op, c.r); err != nil || got != c.want {
			t.Errorf("%q %s %q: %v %v", c.l, c.op, c.r, got, err)
		}
	}
	if _, err := compareValues("abc", "gt", "1"); err == nil {
		t.Error("so lớn với chữ phải báo lỗi")
	}

	// bước trước thành công / thất bại (bỏ qua khi thất bại)
	wait := step(StepWaitImage, func(s *Step) {
		s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs, s.OnFail = FindByText, "Gửi", 1, 1, &OnFail{Action: FailNext}
	})
	sc = Script{Steps: []Step{wait, ifStep("if", &Cond{Kind: CondResult}, []Step{typeStep("a", "đã thấy")}, []Step{typeStep("b", "chưa thấy")})}}
	got, err := runBlocks(t, sc)
	expect(t, "result sai", got, err, "chưa thấy")
	got, err = runBlocks(t, sc, sendBtn)
	expect(t, "result đúng", got, err, "đã thấy")

	// và / hoặc
	both := &Cond{Kind: CondAll, Conds: []Cond{{Kind: CondResult}, {Kind: CondInput, Left: "1", Op: "eq", Right: "2"}}}
	sc.Steps[1].Cond = both
	got, err = runBlocks(t, sc, sendBtn)
	expect(t, "all", got, err, "chưa thấy")
	both.Kind = CondAny
	got, err = runBlocks(t, sc, sendBtn)
	expect(t, "any", got, err, "đã thấy")
}

func TestValidateBlocks(t *testing.T) {
	cases := []struct {
		steps []Step
		want  string
	}{
		{[]Step{step(StepBreak, nil)}, "bước 1: Thoát khối lặp phải nằm trong"},
		{[]Step{ifStep("if", textCond("a"), []Step{step(StepBreak, nil)}, nil)}, "bước 1.1: Thoát khối lặp"},
		{[]Step{ifStep("if", textCond("a"), []Step{typeStep("a", "x")}, []Step{typeStep("b", "")})}, "bước 1.2: chưa nhập nội dung gõ"},
		{[]Step{ifStep("if", &Cond{Kind: CondText}, nil, nil)}, "bước 1: chưa nhập chữ cần tìm"},
		{[]Step{ifStep("if", nil, nil, nil)}, "chưa có điều kiện"},
		{[]Step{groupStep("g", func(s *Step) { s.LoopMode = LoopUntil })}, "cần điều kiện"},
		{[]Step{groupStep("g", nil, typeStep("a", "x")), func() Step {
			s := typeStep("b", "y")
			s.OnFail = &OnFail{Action: FailGoto, Goto: "a"} // bước trong khối khác
			return s
		}()}, "bước 2: khi thất bại nhảy tới bước không còn trong cùng danh sách"},
		{[]Step{typeStep("a", "x"), func() Step { s := typeStep("b", "y"); s.Children = []Step{typeStep("c", "z")}; return s }()}, "chỉ khối"},
	}
	for _, c := range cases {
		err := validateScript(Script{Steps: c.steps}, Rect{})
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("muốn %q, có %v", c.want, err)
		}
	}
	// lỗi bước con trả đúng id + đường dẫn
	err := validateScript(Script{Steps: []Step{typeStep("a", "x"), groupStep("g", nil, typeStep("c", ""))}}, Rect{})
	var se *ScriptError
	if !errorsAs(err, &se) || se.StepID != "c" || se.Path != "2.1" || se.Step != 2 {
		t.Errorf("lỗi bước con: %+v", se)
	}
}

func errorsAs(err error, se **ScriptError) bool {
	e, ok := err.(*ScriptError)
	*se = e
	return ok
}

func TestWalkPathsAndFile(t *testing.T) {
	sc := Script{Steps: []Step{
		typeStep("a", "x"),
		ifStep("if", textCond("a"), []Step{typeStep("t1", "x"), typeStep("t2", "x")},
			[]Step{ifStep("if2", textCond("b"), []Step{typeStep("e1", "x")}, []Step{typeStep("e2", "x")})}),
	}}
	var got []string
	walkSteps(sc.Steps, "", func(s *Step, p string) error {
		got = append(got, s.ID+"="+p)
		return nil
	})
	want := []string{"a=1", "if=2", "t1=2.1", "t2=2.2", "if2=2.3", "e1=2.3.1", "e2=2.3.2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("đánh số %v, muốn %v", got, want)
	}
	// lưu → mở: id trùng trong bước con được đặt lại, version 3
	sc.Steps[1].Else[0].Children[0].ID = "a"
	b, err := encodeScript(sc)
	if err != nil {
		t.Fatal(err)
	}
	back, err := decodeScript(b)
	if err != nil {
		t.Fatal(err)
	}
	if back.Version != 3 || back.Steps[1].Else[0].Children[0].ID == "a" {
		t.Errorf("version %d, id %q", back.Version, back.Steps[1].Else[0].Children[0].ID)
	}
	// file cũ (version 2, không có khối) vẫn mở được
	old, _ := json.Marshal(map[string]any{"version": 2, "name": "cũ", "loops": 1, "steps": []any{map[string]any{"id": "s", "type": "wait", "enabled": true}}})
	if _, err := decodeScript(old); err != nil {
		t.Errorf("file v2: %v", err)
	}
	if err := checkStep(&Step{Type: StepIf, Cond: &Cond{Kind: "lạ"}}); err == nil {
		t.Error("điều kiện lạ phải báo lỗi khi mở file")
	}
}

// Ví dụ khối trong steptypes.go chạy được qua validateScript (AI chép theo ví dụ không bị chặn).
func TestBlockDocExamplesValidate(t *testing.T) {
	for _, d := range stepTypeDocs {
		if !isBlock(d.Type) {
			continue
		}
		var s Step
		if err := json.Unmarshal([]byte(d.Example), &s); err != nil {
			t.Fatal(err)
		}
		if err := validateScript(Script{Steps: []Step{s}, Loops: 1}, Rect{}); err != nil {
			t.Errorf("%s: %v", d.Type, err)
		}
	}
	// thấy chữ → click theo kết quả; khi điều kiện đúng thì click vào chỗ thấy
	var s Step
	json.Unmarshal([]byte(stepTypeDocs[slices.IndexFunc(stepTypeDocs, func(d StepTypeDoc) bool { return d.Type == StepIf })].Example), &s)
	s.Cond.WaitMs = 0
	in := &fakeInput{}
	e := newEngine(in, Script{Loops: 1, Region: &Rect{W: 400, H: 300}, Steps: []Step{s}}, nil)
	e.scr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}}
	e.ocr = &fakeOCR{reads: [][]TextLine{{line("Nhắn tin", 100, 100, 80, 20)}}}
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := in.lines(); len(got) != 2 || got[0] != "move 140,110" {
		t.Errorf("có %q, muốn click vào giữa chữ (140, 110)", got)
	}
}
