package samples

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// Ảnh của trang image-plus.html (bước 40) ghép như trang vẽ ở Retina: so viền thấy icon theme tối (không nhầm hình
// tròn), kiểm tra màu chọn nút xanh (không nhầm nút xám), mặt nạ bỏ qua số đè lên huy hiệu.
func TestImagePlusAssets(t *testing.T) {
	load := func(name string) image.Image {
		f, err := os.Open(filepath.Join("pages", name))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			t.Fatal(err)
		}
		return img
	}
	put := func(dst *image.RGBA, src image.Image, x, y int) {
		draw.Draw(dst, src.Bounds().Add(image.Pt(x, y)), src, image.Point{}, draw.Over)
	}
	scr := image.NewRGBA(image.Rect(0, 0, 1600, 1000))
	draw.Draw(scr, scr.Bounds(), &image.Uniform{color.RGBA{246, 247, 249, 255}}, image.Point{}, draw.Src)
	draw.Draw(scr, image.Rect(0, 0, 800, 500), &image.Uniform{color.RGBA{17, 24, 39, 255}}, image.Point{}, draw.Src)
	put(scr, load("theme-decoy-dark.png"), 100, 100)
	put(scr, load("theme-icon-dark.png"), 500, 300)
	put(scr, load("btn-grey.png"), 900, 100)
	put(scr, load("btn-blue.png"), 1200, 400)
	put(scr, load("badge.png"), 300, 700)
	draw.Draw(scr, image.Rect(340, 744, 390, 780), &image.Uniform{color.White}, image.Point{}, draw.Src) // "số" đè giữa huy hiệu

	area := model.Rect{W: 800, H: 500} // Retina: 2 px ảnh / điểm
	cases := []struct {
		tpl  string
		opt  vision.MatchOptions
		x, y int
		min  float64
	}{
		{"theme-icon-light.png", vision.MatchOptions{Mode: model.MatchEdge}, 500, 300, 0.85},
		{"btn-blue.png", vision.MatchOptions{ColorTol: 15, MinScore: 0.9}, 1200, 400, 0.95},
		{"badge-mask.png", vision.MatchOptions{}, 300, 700, 0.97},
	}
	for _, c := range cases {
		m, _ := vision.MatchScaled(scr, area, vision.NewTemplate(load(c.tpl)), c.opt, 2, 1)
		if m.X != c.x || m.Y != c.y || m.Score < c.min || !m.ColorOK(c.opt) {
			t.Errorf("%s: muốn (%d, %d) ≥ %.2f, có %+v", c.tpl, c.x, c.y, c.min, m)
		}
	}
	// không có các bước 40 thì nhầm / khớp thấp – mẫu thử mới có ý nghĩa
	if m, _ := vision.MatchScaled(scr, area, vision.NewTemplate(load("theme-icon-light.png")), vision.MatchOptions{}, 2, 1); m.X == 500 && m.Score >= 0.8 {
		t.Logf("so xám vẫn thấy icon theme tối (%.2f)", m.Score)
	}
	if m, _ := vision.MatchScaled(scr, area, vision.NewTemplate(load("badge.png")), vision.MatchOptions{}, 2, 1); m.Score >= 0.95 {
		t.Errorf("huy hiệu không mặt nạ mà vẫn khớp %.2f – số đè chưa đủ to để thấy khác biệt", m.Score)
	}
}
