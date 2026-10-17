package main

import (
	"context"
	"image"
	"reflect"
	"strings"
	"testing"
)

// Kịch bản: chờ chữ "Gửi" → click theo kết quả → gõ "a"; bộ đọc chữ giả trả lần lượt reads.
func failScript(f *OnFail, reads ...[]TextLine) (*engine, *fakeInput, *fakeOCR) {
	typeA := step(StepText, func(s *Step) { s.ID, s.Text = "typeA", "a" })
	sc := textScript("Gửi", step(StepClickImage, nil), typeA)
	sc.Steps[0].TimeoutMs, sc.Steps[0].OnFail = 1, f
	in, ocr := &fakeInput{}, &fakeOCR{reads: reads}
	e := newEngine(in, sc, nil)
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}}, ocr
	return e, in, ocr
}

var (
	noText  = []TextLine{line("Đang tải", 0, 0, 80, 20)}
	sendBtn = []TextLine{line("Gửi", 100, 100, 40, 20)}
)

func TestOnFailRetry(t *testing.T) {
	// 2 lần đầu không thấy (mỗi lần chờ đọc ít nhất 1 lần), lần thử lại thứ 2 thấy
	e, in, _ := failScript(&OnFail{Retries: 2, Action: FailStop}, noText, noText, sendBtn)
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"move 220,160", "click left false", "type a"}; !reflect.DeepEqual(in.lines(), want) {
		t.Errorf("có %q, muốn %q", in.lines(), want)
	}
	// hết số lần thử lại → dừng
	e, _, _ = failScript(&OnFail{Retries: 1, Action: FailStop}, noText)
	if err := e.run(context.Background()); err == nil || !strings.Contains(err.Error(), "không thấy chữ") {
		t.Errorf("muốn lỗi không thấy, có %v", err)
	}
}

func TestOnFailNextSkipsStaleClick(t *testing.T) {
	// không thấy → bỏ qua; click theo kết quả bước trước không được click chỗ cũ → cũng lỗi (không có OnFail) → dừng
	e, in, _ := failScript(&OnFail{Action: FailNext}, noText)
	e.found, e.hasFound = Point{X: 1, Y: 1}, true // chỗ thấy từ lần trước
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "bước 2") || len(in.lines()) != 0 {
		t.Errorf("lỗi %v, lệnh %q", err, in.lines())
	}
	// click cũng bỏ qua khi thất bại → chỉ còn gõ chữ
	e, in, _ = failScript(&OnFail{Action: FailNext}, noText)
	e.sc.Steps[1].OnFail = &OnFail{Action: FailNext}
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"type a"}; !reflect.DeepEqual(in.lines(), want) {
		t.Errorf("có %q, muốn %q", in.lines(), want)
	}
}

func TestOnFailGoto(t *testing.T) {
	e, in, _ := failScript(&OnFail{Action: FailGoto, Goto: "typeA"}, noText)
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"type a"}; !reflect.DeepEqual(in.lines(), want) {
		t.Errorf("có %q, muốn %q", in.lines(), want)
	}
	// nhảy ngược: bước 3 (gõ) lỗi → về bước 1; lần 2 bước 1 thấy chữ
	e, in, _ = failScript(nil, sendBtn)
	e.sc.Steps = append(e.sc.Steps, step(StepWaitImage, func(s *Step) {
		s.ID, s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = "w2", FindByText, "Xong", 1, 1
		s.OnFail = &OnFail{Action: FailGoto, Goto: e.sc.Steps[0].ID}
	}))
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "quá 1000 lần") {
		t.Errorf("muốn chặn vòng lặp không dừng, có %v", err)
	}
	if len(in.lines()) < 6 {
		t.Errorf("phải chạy lại từ bước 1 nhiều lần, có %d lệnh", len(in.lines()))
	}
}

func TestValidateOnFail(t *testing.T) {
	sc := textScript("Gửi", step(StepText, func(s *Step) { s.ID, s.Text, s.Enabled = "t", "a", false }))
	for _, c := range []struct {
		f    OnFail
		want string
	}{
		{OnFail{Action: FailGoto, Goto: "mất"}, "không còn"},
		{OnFail{Action: FailGoto, Goto: sc.Steps[0].ID}, "chính bước này"},
		{OnFail{Action: FailGoto, Goto: "t"}, "đang tắt"},
		{OnFail{Action: FailStop, Retries: 101}, "0–100"},
	} {
		sc.Steps[0].OnFail = &c.f
		if err := validateScript(sc, Rect{}); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%+v: %v", c.f, err)
		}
	}
	if err := checkStep(&Step{Type: StepWait, OnFail: &OnFail{Action: "lạ"}}); err == nil {
		t.Error("action lạ phải báo lỗi khi mở file")
	}
}
