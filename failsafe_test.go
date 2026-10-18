package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"autoclicker/internal/model"
)

func TestHardCorners(t *testing.T) {
	// A 1440×900 bên trái, B 1920×1080 bên phải, mép trên thẳng hàng.
	ds := []model.Rect{{X: 0, Y: 0, W: 1440, H: 900}, {X: 1440, Y: 0, W: 1920, H: 1080}}
	got := map[model.Point]bool{}
	for _, c := range hardCorners(ds) {
		got[c] = true
	}
	want := []model.Point{{X: 0, Y: 0}, {X: 0, Y: 899}, {X: 3359, Y: 0}, {X: 3359, Y: 1079}, {X: 1440, Y: 1079}}
	for _, p := range want {
		if !got[p] {
			t.Errorf("thiếu góc %v", p)
		}
	}
	// Góc tiếp giáp 2 màn: chuột đi xuyên qua được → không tính.
	for _, p := range []model.Point{{X: 1439, Y: 0}, {X: 1440, Y: 0}, {X: 1439, Y: 899}} {
		if got[p] {
			t.Errorf("góc tiếp giáp %v không được tính", p)
		}
	}
	if len(got) != len(want) {
		t.Errorf("có %d góc, muốn %d: %v", len(got), len(want), got)
	}
}

func TestInCorner(t *testing.T) {
	corners := hardCorners([]model.Rect{{X: 0, Y: 0, W: 1440, H: 900}})
	cases := []struct {
		x, y int
		want bool
	}{
		{0, 0, true}, {2, 2, true}, {1439, 0, true}, {0, 899, true}, {1437, 897, true},
		{3, 0, false}, {0, 3, false}, {720, 0, false}, {0, 450, false}, {720, 450, false},
	}
	for _, c := range cases {
		if got := inCorner(c.x, c.y, corners); got != c.want {
			t.Errorf("inCorner(%d, %d) = %v, muốn %v", c.x, c.y, got, c.want)
		}
	}
}

// fakeMouse – vị trí chuột giả, đổi được trong lúc watchCorners đang chạy.
type fakeMouse struct {
	mu   sync.Mutex
	x, y int
}

func (m *fakeMouse) set(x, y int) { m.mu.Lock(); m.x, m.y = x, y; m.mu.Unlock() }
func (m *fakeMouse) loc() (int, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.x, m.y
}

func TestWatchCorners(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mouse := &fakeMouse{x: 500, y: 500}
	tin := &trackingInput{Input: &fakeInput{}}
	hit := make(chan struct{}, 1)
	go watchCorners(ctx, mouse.loc, hardCorners([]model.Rect{{X: 0, Y: 0, W: 1440, H: 900}}), tin, func() { hit <- struct{}{} })

	// Engine tự đưa chuột tới góc (điểm của kịch bản) → không dừng.
	tin.Move(0, 0)
	mouse.set(0, 1) // lệch 1 px do làm tròn
	select {
	case <-hit:
		t.Fatal("dừng nhầm khi engine đưa chuột vào góc")
	case <-time.After(3 * failsafeEvery):
	}

	// Người dùng rê chuột sang góc khác → dừng.
	mouse.set(1439, 899)
	select {
	case <-hit:
	case <-time.After(5 * failsafeEvery):
		t.Fatal("không dừng khi người dùng rê chuột vào góc")
	}
}
