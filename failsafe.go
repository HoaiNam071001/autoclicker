package main

import (
	"context"
	"sync"
	"time"
)

// Dừng khẩn cấp (bước 9): người dùng rê chuột vào một góc màn hình chính thì dừng kịch bản.
const (
	cornerSize    = 3                     // cách mép ≤ cornerSize-1 px theo cả 2 chiều = chạm góc
	engineMoveTol = 2                     // lệch ≤ tol px so với điểm engine vừa đưa tới = do engine
	failsafeEvery = 50 * time.Millisecond // chu kỳ đọc vị trí chuột
)

func inCorner(x, y, w, h int) bool {
	nearX := x < cornerSize || x >= w-cornerSize
	nearY := y < cornerSize || y >= h-cornerSize
	return nearX && nearY
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

// watchCorners đọc vị trí chuột định kỳ tới khi ctx bị huỷ; gọi onHit (một lần) khi chuột
// ở góc màn hình (w × h) mà không phải do engine đưa tới.
func watchCorners(ctx context.Context, location func() (int, int), w, h int, t *trackingInput, onHit func()) {
	tick := time.NewTicker(failsafeEvery)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		x, y := location()
		if inCorner(x, y, w, h) && !t.movedByEngine(x, y) {
			onHit()
			return
		}
	}
}
