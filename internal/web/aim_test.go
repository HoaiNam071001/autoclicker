package web

import "testing"

func TestSameBox(t *testing.T) {
	a := findResult{Sel: 0, X: 10, Y: 20, W: 100, H: 30}
	moved, other, zoom := a, a, a
	moved.X += 1.5
	other.Sel = 1
	zoom.Geom.DPR = 2.5
	near := a
	near.Y += 0.8 // làm tròn subpixel giữa 2 lần đọc – vẫn tính là đứng yên
	if !sameBox(a, near) || sameBox(a, moved) || sameBox(a, other) || sameBox(a, zoom) {
		t.Error("sameBox sai")
	}
}
