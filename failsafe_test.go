package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInCorner(t *testing.T) {
	const w, h = 1440, 900
	cases := []struct {
		x, y int
		want bool
	}{
		{0, 0, true}, {2, 2, true}, {1439, 0, true}, {0, 899, true}, {1437, 897, true},
		{3, 0, false}, {0, 3, false}, {720, 0, false}, {0, 450, false}, {720, 450, false},
	}
	for _, c := range cases {
		if got := inCorner(c.x, c.y, w, h); got != c.want {
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
	go watchCorners(ctx, mouse.loc, 1440, 900, tin, func() { hit <- struct{}{} })

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
