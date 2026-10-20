package vision

import (
	"image"
	"image/color"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/vision/visiontest"
)

// Bước 40: so theo viền (đổi theme), mặt nạ (vùng hay đổi trong mẫu), kiểm tra màu (cùng hình khác màu).

// icon vẽ một "biểu tượng" nét fg trên nền bg cỡ 60×60 tại (x, y): khung + chữ thập + chấm lệch (không đối xứng).
func icon(img *image.RGBA, x, y int, fg, bg color.RGBA) {
	for yy := 0; yy < 60; yy++ {
		for xx := 0; xx < 60; xx++ {
			c := bg
			frame := xx < 4 || xx >= 56 || yy < 4 || yy >= 56
			cross := (xx >= 28 && xx < 32 && yy > 10 && yy < 50) || (yy >= 28 && yy < 32 && xx > 10 && xx < 50)
			dot := xx >= 40 && xx < 50 && yy >= 10 && yy < 20
			if frame || cross || dot {
				c = fg
			}
			img.Set(x+xx, y+yy, c)
		}
	}
}

func fill(img *image.RGBA, c color.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}

func TestMatchEdgeAcrossThemes(t *testing.T) {
	light := image.NewRGBA(image.Rect(0, 0, 900, 600))
	fill(light, color.RGBA{245, 245, 245, 255})
	icon(light, 300, 200, color.RGBA{40, 40, 40, 255}, color.RGBA{245, 245, 245, 255})
	tpl := NewTemplate(visiontest.Crop(light, image.Rect(296, 196, 364, 264)))

	dark := visiontest.Scene(900, 600, 11) // nền khác hẳn + vài khối màu làm nhiễu
	for y := 150; y < 330; y++ {
		for x := 520; x < 720; x++ {
			dark.Set(x, y, color.RGBA{30, 30, 34, 255})
		}
	}
	icon(dark, 600, 220, color.RGBA{225, 225, 230, 255}, color.RGBA{30, 30, 34, 255})
	area := model.Rect{W: 900, H: 600}
	if m, _ := MatchScaled(dark, area, tpl, MatchOptions{}, 1, 1); m.Score > 0.7 && abs(m.X-596) < 3 {
		t.Logf("so xám vẫn thấy ở theme tối (%.2f) – không sao, chỉ cần so viền chắc chắn", m.Score)
	}
	m, _ := MatchScaled(dark, area, tpl, MatchOptions{Mode: model.MatchEdge}, 1, 1)
	if abs(m.X-596) > 1 || abs(m.Y-216) > 1 || m.W != 68 || m.H != 68 || m.Score < 0.85 {
		t.Fatalf("so viền: muốn (596, 216) 68×68 ≥ 0.85, có %+v", m)
	}
	// cùng theme vẫn khớp gần tuyệt đối
	if m, _ := MatchScaled(light, area, tpl, MatchOptions{Mode: model.MatchEdge}, 1, 1); m.X != 296 || m.Y != 196 || m.Score < 0.99 {
		t.Fatalf("so viền cùng theme: %+v", m)
	}
}

// Mẫu có vùng trong suốt (số đếm đổi liên tục): vùng đó không tính → vẫn khớp ~1 dù màn có số khác.
func TestMatchMaskIgnoresTransparent(t *testing.T) {
	img := visiontest.Scene(1000, 700, 12)
	visiontest.Button(img, 400, 300, 160, 60)
	src := visiontest.Crop(img, image.Rect(392, 292, 568, 368)) // khoanh có lề 8 px quanh nút
	masked := image.NewNRGBA(src.Bounds())
	for y := 0; y < 76; y++ {
		for x := 0; x < 176; x++ {
			c := src.RGBAAt(x, y)
			if x >= 38 && x < 138 && y >= 28 && y < 48 {
				c.A = 0 // chỗ chữ / số đổi
			}
			masked.Set(x, y, c)
		}
	}
	// trên màn, chỗ đó đổi thành khối đặc
	for y := 320; y < 340; y++ {
		for x := 430; x < 530; x++ {
			img.Set(x, y, color.RGBA{250, 200, 0, 255})
		}
	}
	area := model.Rect{W: 1000, H: 700}
	plain, _ := MatchScaled(img, area, NewTemplate(src), MatchOptions{}, 1, 1)
	tpl := NewTemplate(masked)
	if !tpl.Masked() {
		t.Fatal("ảnh có alpha 0 mà không ra mặt nạ")
	}
	m, _ := MatchScaled(img, area, tpl, MatchOptions{}, 1, 1)
	if m.X != 392 || m.Y != 292 || m.Score < 0.99 || m.Score <= plain.Score {
		t.Fatalf("mặt nạ: muốn (392, 292) ~1 (> %.2f không mặt nạ), có %+v", plain.Score, m)
	}
	// mẫu Retina có mặt nạ tìm trên màn 1x: mặt nạ thu nhỏ theo
	big := image.NewNRGBA(image.Rect(0, 0, 352, 152))
	for y := 0; y < 152; y++ {
		for x := 0; x < 352; x++ {
			big.Set(x, y, masked.At(x/2, y/2))
		}
	}
	if m, _ := MatchScaled(img, area, NewTemplate(big), MatchOptions{}, 2, 1); abs(m.X-392) > 1 || abs(m.Y-292) > 1 || m.Score < 0.95 {
		t.Fatalf("mặt nạ thu nhỏ: %+v", m)
	}
	// so viền + mặt nạ
	if m, _ := MatchScaled(img, area, tpl, MatchOptions{Mode: model.MatchEdge}, 1, 1); m.X != 392 || m.Y != 292 || m.Score < 0.9 {
		t.Fatalf("viền + mặt nạ: %+v", m)
	}
}

// greyButton – nút cùng hình visiontest.Button nhưng xám cùng độ sáng: so xám không phân biệt được.
func greyButton(img *image.RGBA, x, y, w, h int) {
	visiontest.Button(img, x, y, w, h)
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			if c := img.RGBAAt(x+xx, y+yy); c.B == 242 {
				img.Set(x+xx, y+yy, color.RGBA{105, 105, 105, 255})
			}
		}
	}
}

func TestMatchColorCheck(t *testing.T) {
	img := visiontest.Scene(1000, 700, 13)
	greyButton(img, 150, 150, 120, 40)
	visiontest.Button(img, 700, 500, 120, 40)
	tpl := NewTemplate(visiontest.Crop(img, image.Rect(700, 500, 820, 540)))
	area := model.Rect{W: 1000, H: 700}
	opt := MatchOptions{ColorTol: 10, MinScore: 0.9}
	m, _ := MatchScaled(img, area, tpl, opt, 1, 1)
	if m.X != 700 || m.Y != 500 || !m.ColorOK(opt) || m.ColorDiff > 2 {
		t.Fatalf("muốn nút xanh (700, 500) đúng màu, có %+v", m)
	}
	// chỉ còn nút xám: khớp hình cao nhưng lệch màu → không đạt
	visiontest.Button(img, 700, 500, 120, 40)
	for y := 500; y < 540; y++ {
		for x := 700; x < 820; x++ {
			img.Set(x, y, color.RGBA{235, 235, 235, 255})
		}
	}
	m, _ = MatchScaled(img, area, tpl, opt, 1, 1)
	if m.X != 150 || m.Y != 150 || m.Score < 0.95 || m.ColorOK(opt) || m.ColorDiff < 20 {
		t.Fatalf("chỉ nút xám: muốn (150, 150) khớp cao nhưng lệch màu, có %+v", m)
	}
	// không kiểm tra màu: nút xám được nhận
	if m, _ := MatchScaled(img, area, tpl, MatchOptions{}, 1, 1); m.X != 150 || !m.ColorOK(MatchOptions{}) {
		t.Fatalf("không kiểm tra màu: %+v", m)
	}
}

// So viền cả màn Retina: thêm phần tính Sobel cho cả ảnh chụp.
func BenchmarkMatchEdgeRetina(b *testing.B) {
	img := visiontest.Scene(3024, 1964, 7)
	visiontest.Button(img, 2000, 1500, 240, 80)
	tpl := NewTemplate(visiontest.Crop(img, image.Rect(1992, 1492, 2248, 1588)))
	area := model.Rect{W: 1512, H: 982}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if m, _ := MatchScaled(img, area, tpl, MatchOptions{Mode: model.MatchEdge}, 2, 1); m.X != 1992 || m.Y != 1492 {
			b.Fatalf("%+v", m)
		}
	}
}

// Ảnh mẫu trong suốt hết / rất nhỏ: không panic, không lệch khung.
func TestMatchDegenerateTemplates(t *testing.T) {
	img := visiontest.Scene(400, 300, 14)
	area := model.Rect{W: 400, H: 300}
	clear := image.NewNRGBA(image.Rect(0, 0, 40, 30)) // alpha 0 toàn bộ
	for _, opt := range []MatchOptions{{}, {Mode: model.MatchEdge}, {ColorTol: 10}} {
		MatchScaled(img, area, NewTemplate(clear), opt, 1, 1)
		MatchScaled(img, area, NewTemplate(clear), opt, 2, 1)
	}
	tiny := visiontest.Crop(img, image.Rect(100, 100, 102, 102))
	if m, _ := MatchScaled(img, area, NewTemplate(tiny), MatchOptions{Mode: model.MatchEdge}, 1, 1); m.W != 2 || m.H != 2 || m.X < 0 || m.Y < 0 {
		t.Fatalf("mẫu 2×2 so viền: %+v", m)
	}
}
