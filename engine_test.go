package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeInput ghi lại các lệnh dưới dạng chuỗi thay vì điều khiển chuột thật.
type fakeInput struct {
	mu  sync.Mutex
	log []string
}

func (f *fakeInput) add(format string, a ...any) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.log = append(f.log, fmt.Sprintf(format, a...))
	return nil
}

func (f *fakeInput) Move(x, y int) error               { return f.add("move %d,%d", x, y) }
func (f *fakeInput) Drag(x, y int, b string) error     { return f.add("drag %d,%d %s", x, y, b) }
func (f *fakeInput) Toggle(b string, down bool) error  { return f.add("toggle %s %v", b, down) }
func (f *fakeInput) Click(b string, double bool) error { return f.add("click %s %v", b, double) }
func (f *fakeInput) TypeStr(s string) error            { return f.add("type %s", s) }
func (f *fakeInput) KeyTap(k string, mods []string) error {
	return f.add("key %s %s", strings.Join(mods, "+"), k)
}

func (f *fakeInput) lines() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.log...)
}

func step(t string, mod func(*Step)) Step {
	s, err := (&App{}).NewStep(t)
	if err != nil {
		panic(err)
	}
	s.DelayAfter = 0
	if mod != nil {
		mod(&s)
	}
	return s
}

func TestRunClickLoopsAndRegion(t *testing.T) {
	in := &fakeInput{}
	sc := Script{
		Loops:  2,
		Region: &Rect{X: 100, Y: 50},
		Steps: []Step{
			step(StepClick, func(s *Step) { s.Points = []Point{{10, 20}}; s.Clicks = 2 }),
			step(StepDouble, func(s *Step) { s.Points = []Point{{1, 2}}; s.Button = "right" }),
			step(StepClick, func(s *Step) { s.Enabled = false }),
		},
	}
	if err := newEngine(in, sc, nil).run(context.Background()); err != nil {
		t.Fatal(err)
	}
	one := []string{"move 110,70", "click left false", "move 110,70", "click left false", "move 101,52", "click right true"}
	want := append(append([]string{}, one...), one...)
	if got := in.lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestRunDrag(t *testing.T) {
	in := &fakeInput{}
	sc := Script{Loops: 1, Steps: []Step{
		step(StepDrag, func(s *Step) { s.Points = []Point{{0, 0}, {100, 40}}; s.DragMs = 40 }),
	}}
	if err := newEngine(in, sc, nil).run(context.Background()); err != nil {
		t.Fatal(err)
	}
	got := in.lines()
	want := []string{"move 0,0", "toggle left true", "drag 25,10 left", "drag 50,20 left", "drag 75,30 left", "drag 100,40 left", "toggle left false"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestRunTypeAndKey(t *testing.T) {
	in := &fakeInput{}
	sc := Script{Loops: 1, Steps: []Step{
		step(StepText, func(s *Step) { s.Text = "Gõ" }),
		step(StepKey, func(s *Step) { s.Key = "v"; s.Modifiers = []string{"cmd", "shift"} }),
	}}
	if err := newEngine(in, sc, nil).run(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := []string{"type G", "type õ", "key cmd+shift v"}
	if got := in.lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestStopInfiniteLoopReleasesButton(t *testing.T) {
	in := &fakeInput{}
	sc := Script{Loops: 0, Steps: []Step{
		step(StepHold, func(s *Step) { s.HoldMs = 10_000 }),
	}}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error)
	go func() { done <- newEngine(in, sc, nil).run(ctx) }()
	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v, muốn context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("không dừng sau khi huỷ")
	}
	got := in.lines()
	if got[len(got)-1] != "toggle left false" {
		t.Errorf("chưa nhả chuột khi bị huỷ: %q", got)
	}
}

func TestValidateScript(t *testing.T) {
	cases := map[string]Script{
		"không có bước nào được bật": {Steps: []Step{step(StepWait, func(s *Step) { s.Enabled = false })}},
		"cần 2 điểm":                 {Steps: []Step{step(StepDrag, func(s *Step) { s.Points = s.Points[:1] })}},
		"chưa nhập nội dung gõ":      {Steps: []Step{step(StepText, nil)}},
		"chưa chọn phím":             {Steps: []Step{step(StepKey, func(s *Step) { s.Key = "" })}},
		"không hợp lệ":               {Steps: []Step{{Type: "scroll", Enabled: true}}},
	}
	for want, sc := range cases {
		err := validateScript(sc)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi chứa %q, nhận %v", want, err)
		}
	}
}

func TestPauseResumeAndProgress(t *testing.T) {
	in := &fakeInput{}
	sc := Script{Loops: 2, Steps: []Step{
		step(StepClick, func(s *Step) { s.Points = []Point{{1, 1}}; s.Clicks = 3; s.ClickInterval = 20 }),
		step(StepWait, func(s *Step) { s.DelayAfter = 10 }),
	}}
	var mu sync.Mutex
	var progs []Progress
	eng := newEngine(in, sc, func(p Progress) { mu.Lock(); progs = append(progs, p); mu.Unlock() })

	done := make(chan error)
	go func() { done <- eng.run(context.Background()) }()
	time.Sleep(30 * time.Millisecond) // đang giữa các click của vòng 1
	if !eng.Pause() || eng.Pause() {
		t.Fatal("Pause lần 1 phải true, lần 2 false")
	}
	time.Sleep(20 * time.Millisecond) // chờ lần click đang dở (nếu có) xong
	n := len(in.lines())
	time.Sleep(100 * time.Millisecond)
	if got := len(in.lines()); got != n {
		t.Fatalf("vẫn chạy khi tạm dừng: %d → %d lệnh", n, got)
	}
	if p := eng.Snapshot(); !p.Paused || p.ElapsedMs > 80 {
		t.Errorf("snapshot khi dừng = %+v (elapsed không được tính lúc tạm dừng)", p)
	}
	eng.Resume()
	if err := <-done; err != nil {
		t.Fatal(err)
	}

	last := eng.Snapshot()
	if last.Clicks != 6 || last.Loop != 2 || last.Step != 2 || last.Total != 2 || last.Paused {
		t.Errorf("tiến độ cuối = %+v", last)
	}
	mu.Lock()
	defer mu.Unlock()
	sawPaused := false
	for _, p := range progs {
		sawPaused = sawPaused || p.Paused
	}
	if !sawPaused {
		t.Error("không báo tiến độ khi tạm dừng")
	}
}
