package vision

import (
	"image"
	"testing"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/vision/visiontest"
)

// Nhiều nút giống nhau (2 cột × 6 hàng, sát nhau) → thấy đủ, mỗi nút một lần, theo thứ tự đọc.
func TestMatchAllGrid(t *testing.T) {
	img := visiontest.Scene(1400, 900, 7)
	var want []model.Rect
	for row := 0; row < 6; row++ {
		for col := 0; col < 2; col++ {
			x, y := 200+col*500, 100+row*110
			visiontest.Button(img, x, y, 120, 40)
			want = append(want, model.Rect{X: x, Y: y})
		}
	}
	tpl := NewTemplate(visiontest.Crop(img, image.Rect(200, 100, 320, 140)))
	area := model.Rect{W: 1400, H: 900}
	start := time.Now()
	ms, _ := MatchAll(img, area, tpl, MatchOptions{MinScore: 0.9}, 0, 1)
	t.Logf("1400×900, 12 nút: %v", time.Since(start))
	if len(ms) != len(want) {
		t.Fatalf("muốn %d chỗ, có %d: %+v", len(want), len(ms), ms)
	}
	for i, m := range ms {
		if abs(m.X-want[i].X) > 1 || abs(m.Y-want[i].Y) > 1 {
			t.Errorf("chỗ %d: muốn (%d, %d), có (%d, %d)", i+1, want[i].X, want[i].Y, m.X, m.Y)
		}
	}
}

// Không có nút → rỗng; ngưỡng cao loại chỗ na ná.
func TestMatchAllNone(t *testing.T) {
	other := visiontest.Scene(600, 400, 4)
	visiontest.Button(other, 100, 100, 120, 40)
	tpl := NewTemplate(visiontest.Crop(other, image.Rect(100, 100, 220, 140)))
	ms, _ := MatchAll(visiontest.Scene(600, 400, 3), model.Rect{W: 600, H: 400}, tpl, MatchOptions{MinScore: 0.9}, 0, 1)
	if len(ms) != 0 {
		t.Fatalf("muốn rỗng, có %+v", ms)
	}
}

func TestReadingLess(t *testing.T) {
	a, b := model.Rect{X: 300, Y: 10, W: 50, H: 20}, model.Rect{X: 10, Y: 14, W: 50, H: 20} // cùng hàng
	if !ReadingLess(b, a) || ReadingLess(a, b) {
		t.Error("cùng hàng: trái trước")
	}
	c := model.Rect{X: 0, Y: 60, W: 50, H: 20}
	if !ReadingLess(a, c) {
		t.Error("hàng trên trước hàng dưới")
	}
}
