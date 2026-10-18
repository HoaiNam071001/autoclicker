package engine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script/scripttest"
)

func TestRunClickLoopsAndRegion(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{
		Loops:  2,
		Region: &model.Rect{X: 100, Y: 50},
		Steps: []model.Step{
			modeltest.Step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 10, Y: 20}}; s.Clicks = 2 }),
			modeltest.Step(model.StepDouble, func(s *model.Step) { s.Points = []model.Point{{X: 1, Y: 2}}; s.Button = "right" }),
			modeltest.Step(model.StepClick, func(s *model.Step) { s.Enabled = false }),
		},
	}
	if err := newEngine(in, scripttest.Linear(sc), nil).Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	one := []string{"move 110,70", "click left false", "move 110,70", "click left false", "move 101,52", "click right true"}
	want := append(append([]string{}, one...), one...)
	if got := in.Lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestRunDrag(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		modeltest.Step(model.StepDrag, func(s *model.Step) { s.Points = []model.Point{{X: 0, Y: 0}, {X: 100, Y: 40}}; s.DragMs = 40 }),
	}}
	if err := newEngine(in, scripttest.Linear(sc), nil).Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	got := in.Lines()
	want := []string{"move 0,0", "toggle left true", "drag 25,10 left", "drag 50,20 left", "drag 75,30 left", "drag 100,40 left", "toggle left false"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestRunTypeAndKey(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		modeltest.Step(model.StepText, func(s *model.Step) { s.Text = "Gõ" }),
		modeltest.Step(model.StepKey, func(s *model.Step) { s.Keys = []string{"v"}; s.Modifiers = []string{"cmd", "shift"} }),
	}}
	if err := newEngine(in, scripttest.Linear(sc), nil).Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := []string{"type G", "type õ",
		"down cmd []", "down shift []", "down v [cmd+shift]", "up v", "up shift", "up cmd"}
	if got := in.Lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestStopInfiniteLoopReleasesButton(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 0, Steps: []model.Step{
		modeltest.Step(model.StepHold, func(s *model.Step) { s.HoldMs = 10_000 }),
	}}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error)
	go func() { done <- newEngine(in, scripttest.Linear(sc), nil).Run(ctx) }()
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
	got := in.Lines()
	if got[len(got)-1] != "toggle left false" {
		t.Errorf("chưa nhả chuột khi bị huỷ: %q", got)
	}
}

func TestPauseResumeAndProgress(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 2, Steps: []model.Step{
		modeltest.Step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 1, Y: 1}}; s.Clicks = 3; s.ClickInterval = 20 }),
		modeltest.Step(model.StepWait, func(s *model.Step) { s.DelayAfter = 10 }),
	}}
	var mu sync.Mutex
	var progs []Progress
	eng := newEngine(in, scripttest.Linear(sc), func(p Progress) { mu.Lock(); progs = append(progs, p); mu.Unlock() })

	done := make(chan error)
	go func() { done <- eng.Run(context.Background()) }()
	time.Sleep(30 * time.Millisecond) // đang giữa các click của vòng 1
	if !eng.Pause() || eng.Pause() {
		t.Fatal("Pause lần 1 phải true, lần 2 false")
	}
	time.Sleep(20 * time.Millisecond) // chờ lần click đang dở (nếu có) xong
	n := len(in.Lines())
	time.Sleep(100 * time.Millisecond)
	if got := len(in.Lines()); got != n {
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

func TestKeyChordAndSequence(t *testing.T) {
	cases := map[string][]string{
		// Ctrl+C+T nhấn giữ cùng lúc rồi nhả ngược thứ tự
		model.KeyChord: {"down ctrl []", "down c [ctrl]", "down t [ctrl]", "up t", "up c", "up ctrl"},
		// giữ Ctrl, gõ C rồi T
		model.KeySequence: {"down ctrl []", "down c [ctrl]", "up c", "down t [ctrl]", "up t", "up ctrl"},
	}
	for mode, want := range cases {
		in := &devicetest.FakeInput{}
		sc := model.Script{Loops: 1, Steps: []model.Step{modeltest.Step(model.StepKey, func(s *model.Step) {
			s.Keys = []string{"c", "t"}
			s.Modifiers = []string{"ctrl"}
			s.KeyMode = mode
		})}}
		if err := newEngine(in, scripttest.Linear(sc), nil).Run(context.Background()); err != nil {
			t.Fatal(err)
		}
		if got := in.Lines(); !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %q\nwant %q", mode, got, want)
		}
	}
}

func TestKeyCancelReleasesAll(t *testing.T) {
	in := &devicetest.FakeInput{}
	s := modeltest.Step(model.StepKey, func(s *model.Step) { s.Keys = []string{"c", "t"}; s.Modifiers = []string{"ctrl"} })
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // huỷ trước: lần chờ giữ phím trả lỗi ngay
	if err := newEngine(in, model.Script{}, nil).pressKeys(ctx, s); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
	want := []string{"down ctrl []", "down c [ctrl]", "down t [ctrl]", "up t", "up c", "up ctrl"}
	if got := in.Lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestJitter(t *testing.T) {
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(model.Script{}), nil)
	e.intN = func(n int) int { return 0 } // lệch nhỏ nhất
	if got := e.jitter(100, 5); got != 95 {
		t.Errorf("jitter min = %d, muốn 95", got)
	}
	if got := e.jitter(-1500, 5); got != -1505 {
		t.Errorf("toạ độ âm (màn phụ) = %d, muốn -1505", got)
	}
	if got := e.jitterMs(3, 10); got != 0 {
		t.Errorf("jitterMs không được âm: %d", got)
	}
	e.intN = func(n int) int { return n - 1 } // lệch lớn nhất
	if got := e.jitter(100, 5); got != 105 {
		t.Errorf("jitter max = %d, muốn 105", got)
	}
	if got := e.jitter(100, 0); got != 100 {
		t.Errorf("jitter 0 phải giữ nguyên: %d", got)
	}
}

func TestRunJitterPositionsInBounds(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		modeltest.Step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 500, Y: 300}}; s.Clicks = 200; s.JitterPx = 4 }),
	}}
	if err := newEngine(in, scripttest.Linear(sc), nil).Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, l := range in.Lines() {
		var x, y int
		if _, err := fmt.Sscanf(l, "move %d,%d", &x, &y); err != nil {
			continue
		}
		if x < 496 || x > 504 || y < 296 || y > 304 {
			t.Fatalf("điểm %d,%d ngoài ±4 px", x, y)
		}
		seen[l] = true
	}
	if len(seen) < 10 {
		t.Errorf("200 click mà chỉ có %d vị trí khác nhau – có ngẫu nhiên không?", len(seen))
	}
}

func TestDragFromPrevAndRelative(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 1, Region: &model.Rect{X: 1000, Y: 0}, Steps: []model.Step{
		modeltest.Step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 100, Y: 100}} }),
		// kéo từ chỗ vừa click, sang phải 40, lên 20
		modeltest.Step(model.StepDrag, func(s *model.Step) {
			s.Points = []model.Point{{X: 0, Y: 0}, {X: 40, Y: -20}}
			s.DragFromPrev, s.DragRelative, s.DragMs = true, true, 20
		}),
		// kéo tiếp từ điểm thả vừa rồi tới điểm cố định (trong vùng)
		modeltest.Step(model.StepDrag, func(s *model.Step) {
			s.Points = []model.Point{{X: 0, Y: 0}, {X: 10, Y: 10}}
			s.DragFromPrev, s.DragMs = true, 10
		}),
	}}
	if err := newEngine(in, scripttest.Linear(sc), nil).Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"move 1100,100", "click left false",
		"move 1100,100", "toggle left true", "drag 1120,90 left", "drag 1140,80 left", "toggle left false",
		"move 1140,80", "toggle left true", "drag 1010,10 left", "toggle left false",
	}
	if got := in.Lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestDragFromPrevNeedsPrevPoint(t *testing.T) {
	sc := model.Script{Loops: 1, Steps: []model.Step{
		modeltest.Step(model.StepKey, nil),
		modeltest.Step(model.StepDrag, func(s *model.Step) { s.DragFromPrev = true }),
	}}
	err := newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), nil).Run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "n2: kéo từ điểm trước đó") {
		t.Errorf("muốn lỗi thiếu điểm trước, nhận %v", err)
	}
}
