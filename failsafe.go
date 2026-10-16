package main

import (
	"context"
	"sync"
	"time"
)

// Dừng khẩn cấp (bước 9, 17): người dùng rê chuột vào một góc "cứng" của bất kỳ màn hình nào thì dừng.
const (
	cornerSize    = 3                     // cách góc < cornerSize px theo cả 2 chiều = chạm góc
	engineMoveTol = 2                     // lệch ≤ tol px so với điểm engine vừa đưa tới = do engine
	failsafeEvery = 50 * time.Millisecond // chu kỳ đọc vị trí chuột
)

// hardCorners trả các góc màn hình mà chuột bị chặn lại – bỏ góc tiếp giáp màn khác (chuột đi xuyên
// qua được, tính vào sẽ dừng nhầm khi rê chuột giữa 2 màn).
func hardCorners(displays []Rect) []Point {
	inAny := func(x, y int) bool {
		for _, d := range displays {
			if x >= d.X && x < d.X+d.W && y >= d.Y && y < d.Y+d.H {
				return true
			}
		}
		return false
	}
	var out []Point
	for _, d := range displays {
		for _, c := range []struct{ x, y, dx, dy int }{
			{d.X, d.Y, -1, -1},
			{d.X + d.W - 1, d.Y, 1, -1},
			{d.X, d.Y + d.H - 1, -1, 1},
			{d.X + d.W - 1, d.Y + d.H - 1, 1, 1},
		} {
			if !inAny(c.x+c.dx, c.y) && !inAny(c.x, c.y+c.dy) && !inAny(c.x+c.dx, c.y+c.dy) {
				out = append(out, Point{X: c.x, Y: c.y})
			}
		}
	}
	return out
}

func inCorner(x, y int, corners []Point) bool {
	for _, c := range corners {
		if abs(x-c.X) < cornerSize && abs(y-c.Y) < cornerSize {
			return true
		}
	}
	return false
}

// trackingInput bọc Input, ghi lại điểm engine vừa đưa chuột tới – để kịch bản có điểm ở góc
// màn hình không tự kích hoạt dừng khẩn cấp.
type trackingInput struct {
	Input
	mu      sync.Mutex
	lastX   int
	lastY   int
	hasLast bool
}

func (t *trackingInput) record(x, y int) {
	t.mu.Lock()
	t.lastX, t.lastY, t.hasLast = x, y, true
	t.mu.Unlock()
}

func (t *trackingInput) Move(x, y int) error {
	t.record(x, y)
	return t.Input.Move(x, y)
}

func (t *trackingInput) Drag(x, y int, button string) error {
	t.record(x, y)
	return t.Input.Drag(x, y, button)
}

// movedByEngine – (x, y) có phải điểm engine vừa đưa chuột tới không.
func (t *trackingInput) movedByEngine(x, y int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.hasLast && abs(x-t.lastX) <= engineMoveTol && abs(y-t.lastY) <= engineMoveTol
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// watchCorners đọc vị trí chuột định kỳ tới khi ctx bị huỷ; gọi onHit (một lần) khi chuột
// ở một trong các góc corners mà không phải do engine đưa tới.
func watchCorners(ctx context.Context, location func() (int, int), corners []Point, t *trackingInput, onHit func()) {
	tick := time.NewTicker(failsafeEvery)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		x, y := location()
		if inCorner(x, y, corners) && !t.movedByEngine(x, y) {
			onHit()
			return
		}
	}
}
