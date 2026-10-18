package humanize

import (
	"math"
	"math/rand/v2"
	"strings"
	"testing"

	"autoclicker/internal/model"
)

func rng(seed uint64) *rand.Rand { return rand.New(rand.NewPCG(seed, seed)) }

func total(path []Move) (ms int) {
	for _, m := range path {
		ms += m.WaitMs
	}
	return ms
}

// Mọi mức, nhiều seed: kết thúc đúng đích, điểm giữa nằm trong màn (cách mép), thời gian theo Fitts trong khoảng,
// không bước nào nhảy quá xa (không "dịch chuyển tức thời").
func TestPathEndsInsideAndTimed(t *testing.T) {
	screen := model.Rect{W: 1440, H: 900}
	cases := []struct{ from, to model.Point }{
		{model.Point{X: 100, Y: 100}, model.Point{X: 1300, Y: 800}},
		{model.Point{X: 1430, Y: 5}, model.Point{X: 1435, Y: 890}}, // sát mép phải
		{model.Point{X: 700, Y: 450}, model.Point{X: 760, Y: 470}},
	}
	for _, lv := range []string{Low, Medium, High} {
		p, _ := ProfileOf(lv)
		for seed := uint64(1); seed <= 200; seed++ {
			for _, c := range cases {
				path := Path(rng(seed), p, c.from, c.to, Target{Size: 30}, screen)
				end := path[len(path)-1]
				if end.X != c.to.X || end.Y != c.to.Y {
					t.Fatalf("%s seed %d: dừng ở %d,%d, muốn %v", lv, seed, end.X, end.Y, c.to)
				}
				prev := c.from
				for _, m := range path[:len(path)-1] {
					if m.X < edgeMargin || m.Y < edgeMargin || m.X > screen.W-1-edgeMargin || m.Y > screen.H-1-edgeMargin {
						t.Fatalf("%s seed %d: điểm %d,%d sát / ra ngoài mép", lv, seed, m.X, m.Y)
					}
					if math.Hypot(float64(m.X-prev.X), float64(m.Y-prev.Y)) > 120 {
						t.Fatalf("%s seed %d: bước nhảy %v → %d,%d", lv, seed, prev, m.X, m.Y)
					}
					prev = model.Point{X: m.X, Y: m.Y}
				}
				d := math.Hypot(float64(c.to.X-c.from.X), float64(c.to.Y-c.from.Y))
				lo, hi := Fitts(d, 30)*p.Speed*0.8, Fitts(d, 30)*p.Speed*1.2+650 // + đi lố và chỉnh lại
				if ms := float64(total(path)); ms < lo || ms > hi {
					t.Fatalf("%s seed %d %v→%v: %v ms, muốn %v–%v", lv, seed, c.from, c.to, ms, lo, hi)
				}
			}
		}
	}
}

// Đường cong thật: điểm giữa lệch khỏi đoạn thẳng ở phần lớn seed; DurMs cố định thì đúng thời gian, không đi lố.
func TestPathCurvedAndFixedDuration(t *testing.T) {
	p, _ := ProfileOf(Medium)
	a, b := model.Point{X: 0, Y: 500}, model.Point{X: 1000, Y: 500}
	curved := 0
	for seed := uint64(1); seed <= 100; seed++ {
		path := Path(rng(seed), p, a, b, Target{DurMs: 600}, model.Rect{})
		if ms := total(path); ms < 590 || ms > 600 {
			t.Fatalf("seed %d: %d ms, muốn ≈ 600", seed, ms)
		}
		maxDev := 0
		for _, m := range path {
			if m.X > b.X+2 {
				t.Fatalf("seed %d: kéo thả không được đi lố (%d)", seed, m.X)
			}
			maxDev = max(maxDev, abs(m.Y-500))
		}
		if maxDev >= 20 {
			curved++
		}
	}
	if curved < 50 {
		t.Errorf("chỉ %d / 100 đường cong rõ (≥ 20 px)", curved)
	}
	if got := Path(rng(1), p, a, model.Point{X: 1, Y: 500}, Target{}, model.Rect{}); len(got) != 1 {
		t.Errorf("quãng 1 px phải đi thẳng: %v", got)
	}
}

func TestAim(t *testing.T) {
	r := rng(7)
	moved := false
	for i := 0; i < 500; i++ {
		x, y := Aim(r, 100, 200, 60, 30)
		if abs(x-100) > 10 || abs(y-200) > 5 {
			t.Fatalf("lệch quá 1/3 giữa đích: %d,%d", x, y)
		}
		moved = moved || x != 100 || y != 200
	}
	if !moved {
		t.Error("không lần nào lệch khỏi tâm")
	}
	if x, y := Aim(r, 5, 6, 0, 30); x != 5 || y != 6 {
		t.Errorf("không rõ cỡ đích phải đúng tâm: %d,%d", x, y)
	}
}

// Gõ: gõ nhầm rồi xoá vẫn ra đúng chuỗi; nhịp không đều; mức Vừa không gõ nhầm.
func TestTyping(t *testing.T) {
	high, _ := ProfileOf(High)
	high.Typo = 0.5
	text := "Xin chao, the gioi! 123"
	typos := 0
	for seed := uint64(1); seed <= 50; seed++ {
		var got []rune
		for _, k := range Typing(rng(seed), high, text) {
			if k.Text == Backspace {
				got, typos = got[:len(got)-1], typos+1
				continue
			}
			got = append(got, []rune(k.Text)...)
		}
		if string(got) != text {
			t.Fatalf("seed %d: ra %q", seed, string(got))
		}
	}
	if typos == 0 {
		t.Error("Typo 0.5 mà không lần nào gõ nhầm")
	}
	med, _ := ProfileOf(Medium)
	keys := Typing(rng(3), med, text)
	waits := map[int]bool{}
	for i, k := range keys {
		if k.Text == Backspace {
			t.Fatal("mức Vừa không gõ nhầm")
		}
		if i > 0 && (k.WaitMs < med.KeyGap[0] || k.WaitMs > med.KeyGap[1]+med.PauseMs[1]) {
			t.Fatalf("ký tự %d chờ %d ms", i, k.WaitMs)
		}
		waits[k.WaitMs] = true
	}
	if keys[0].WaitMs != 0 || len(waits) < 5 || strings.Count(text, "") != len(keys)+1 {
		t.Errorf("nhịp gõ: %v", keys)
	}
}

func TestProfileOf(t *testing.T) {
	for _, lv := range []string{"", Off, "nhanh"} {
		if _, ok := ProfileOf(lv); ok {
			t.Errorf("%q phải là tắt", lv)
		}
	}
	if Effective(Medium, "") != Medium || Effective(Medium, Off) != Off || Effective("", High) != High {
		t.Error("Effective sai")
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
