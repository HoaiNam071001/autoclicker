package engine

import (
	"context"
	"errors"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
	"time"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/humanize"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

// humanEngine – engine chạy một node với mức giống người level của kịch bản, seed cố định.
func humanEngine(level string, s model.Step) (*Engine, *devicetest.FakeInput) {
	in := &devicetest.FakeInput{}
	e := newEngine(in, model.Script{Loops: 1, Start: s.ID, Steps: []model.Step{s}, Human: level}, nil)
	e.rng = rand.New(rand.NewPCG(1, 2))
	return e, in
}

// không tính các bước giữa đường (step x,y)
func actions(in *devicetest.FakeInput) []string {
	return slices.DeleteFunc(in.Lines(), func(l string) bool { return strings.HasPrefix(l, "step ") })
}

// Click giống người: đi nhiều bước tới đích, bước cuối là move (chờ app nhận), nhấn – giữ – nhả thay cho click.
func TestHumanClick(t *testing.T) {
	e, in := humanEngine(model.HumanMedium, modeltest.Node("c", "n1", model.StepClick, func(s *model.Step) {
		s.Points = []model.Point{{X: 300, Y: 200}}
	}))
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if steps := len(in.Lines()) - len(actions(in)); steps < 5 {
		t.Errorf("chỉ %d bước giữa đường: %q", steps, in.Lines())
	}
	if want := []string{"move 300,200", "toggle left true", "toggle left false"}; !slices.Equal(actions(in), want) {
		t.Errorf("lệnh %q, muốn %q", actions(in), want)
	}
	// node ghi đè tắt → như cũ
	e, in = humanEngine(model.HumanMedium, modeltest.Node("c", "n1", model.StepClick, func(s *model.Step) {
		s.Points, s.Human = []model.Point{{X: 300, Y: 200}}, model.HumanOff
	}))
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"move 300,200", "click left false"}; !slices.Equal(in.Lines(), want) {
		t.Errorf("node tắt: %q", in.Lines())
	}
}

// Dừng giữa lúc đang di chuột: trả về ngay, chưa bấm.
func TestHumanMoveCancel(t *testing.T) {
	e, in := humanEngine(model.HumanHigh, modeltest.Node("c", "n1", model.StepClick, func(s *model.Step) {
		s.Points = []model.Point{{X: 1400, Y: 900}}
	}))
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(80*time.Millisecond, cancel)
	start := time.Now()
	if err := e.Run(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
	if d := time.Since(start); d > 150*time.Millisecond {
		t.Errorf("dừng sau %v", d)
	}
	if len(actions(in)) != 0 || len(in.Lines()) == 0 {
		t.Errorf("đã di một đoạn rồi dừng, chưa bấm: %q", in.Lines())
	}
}

// Kéo thả giống người: đi cong tới B bằng drag, nhả ở B.
func TestHumanDrag(t *testing.T) {
	e, in := humanEngine(model.HumanMedium, modeltest.Node("d", "n1", model.StepDrag, func(s *model.Step) {
		s.Points, s.DragMs = []model.Point{{X: 100, Y: 100}, {X: 500, Y: 100}}, 300
	}))
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	lines := actions(in)
	i := slices.Index(lines, "toggle left true")
	if i < 1 || lines[i-1] != "move 100,100" || lines[len(lines)-1] != "toggle left false" || lines[len(lines)-2] != "drag 500,100 left" {
		t.Fatalf("lệnh %q", lines)
	}
	curved := false
	for _, l := range lines[i+1 : len(lines)-2] {
		curved = curved || !strings.HasSuffix(l, ",100 left")
	}
	if !curved || len(lines)-i < 10 {
		t.Errorf("đường kéo không cong / quá ít bước: %q", lines[i:])
	}
}

// Gõ giống người: gõ nhầm + xoá lùi (Typo cao, nhịp 0 cho nhanh), chữ cuối cùng vẫn đúng.
func TestHumanTyping(t *testing.T) {
	text := "the quick brown fox"
	e, in := humanEngine(model.HumanHigh, model.Step{})
	e.human = &humanize.Profile{Typo: 0.5}
	if err := e.typeText(context.Background(), text); err != nil {
		t.Fatal(err)
	}
	var got []string
	fixed := 0
	for _, l := range in.Lines() {
		switch {
		case strings.HasPrefix(l, "type "):
			got = append(got, strings.TrimPrefix(l, "type "))
		case l == "down backspace []":
			got, fixed = got[:len(got)-1], fixed+1
		}
	}
	if strings.Join(got, "") != text || fixed == 0 {
		t.Errorf("gõ ra %q, sửa %d lần", strings.Join(got, ""), fixed)
	}
}
