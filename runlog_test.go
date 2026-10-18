package main

import (
	"context"
	"image"
	"strings"
	"sync"
	"testing"
	"time"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

func TestShootOnFail(t *testing.T) {
	e, _ := sendFlow(-1, nil)
	e.sc.Edges = e.sc.Edges[:1] // không thấy → dừng
	e.shots = map[string]image.Image{}
	var got []model.StepResult
	e.onNode = func(r model.StepResult) { got = append(got, r) }
	if err := e.run(context.Background()); err == nil {
		t.Fatal("phải lỗi")
	}
	if len(got) != 1 || got[0].Shot != "1.png" || e.shots["1.png"] == nil || got[0].Loop != 1 {
		t.Errorf("kết quả %+v, ảnh %v", got, e.shots)
	}
	if img := vision.Shrink(image.NewRGBA(image.Rect(0, 0, 3200, 2000)), shotMaxW); img.Bounds().Dx() != 1600 || img.Bounds().Dy() != 1000 {
		t.Errorf("thu nhỏ: %v", img.Bounds())
	}
}

// Debug: breakpoint chỉ dừng khi debug bật; chạy từng node dừng trước mọi node; Tiếp tắt chạy từng node.
func TestDebugBreakpointAndStep(t *testing.T) {
	sc := model.Script{Loops: 1, Start: "a", Steps: []model.Step{
		node("a", "n1", model.StepText, func(s *model.Step) { s.Text = "a" }),
		node("b", "n2", model.StepText, func(s *model.Step) { s.Text = "b"; s.Breakpoint = true }),
		node("c", "n3", model.StepText, func(s *model.Step) { s.Text = "c" }),
	}, Edges: []model.Edge{edge("a", model.PortSuccess, "b"), edge("b", model.PortSuccess, "c")}}

	// không debug: chạy thẳng
	if err := newEngine(&devicetest.FakeInput{}, sc, nil).run(context.Background()); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var breaks []string
	in := &devicetest.FakeInput{}
	e := newEngine(in, sc, func(p Progress) {
		if p.Break != "" {
			mu.Lock()
			breaks = append(breaks, p.Ref+":"+p.Break)
			mu.Unlock()
		}
	})
	e.debug = true
	done := make(chan error)
	go func() { done <- e.run(context.Background()) }()
	waitFor := func(n int) {
		t.Helper()
		for i := 0; i < 100; i++ {
			mu.Lock()
			k := len(breaks)
			mu.Unlock()
			if k >= n && e.gate.paused() {
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
		t.Fatalf("không dừng lần %d (%v)", n, breaks)
	}
	waitFor(1) // dừng trước n2 – n1 đã chạy
	if l := in.Lines(); len(l) != 1 || l[0] != "type a" {
		t.Errorf("trước breakpoint: %q", l)
	}
	e.Step() // chạy n2 rồi dừng trước n3
	waitFor(2)
	if l := in.Lines(); len(l) != 2 {
		t.Errorf("sau F10: %q", l)
	}
	e.Resume()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if w := "n2:breakpoint n3:step"; strings.Join(breaks, " ") != w {
		t.Errorf("dừng ở %v, muốn %s", breaks, w)
	}
}
