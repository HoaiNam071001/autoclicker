package device

import (
	"os"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// tess-page.tsv: Tesseract vie đọc trang thử 40 O0 (nửa sáng, ảnh 700×372 phóng ×2 qua vision.PrepOCR).
func TestParseTessTSV(t *testing.T) {
	raw, err := os.ReadFile("testdata/tess-page.tsv")
	if err != nil {
		t.Fatal(err)
	}
	lines := parseTessTSV(string(raw), 2, 0)
	if len(lines) != 11 {
		t.Fatalf("%d dòng", len(lines))
	}
	first := lines[0]
	if first.Text != "Đăng nhập bằng tài khoản của bạn" || first.Approx || len(first.Words) != 7 {
		t.Errorf("dòng đầu: %+v", first)
	}
	if w := first.Words[1]; w.Start != 5 || w.Len != 4 || w.Box != (model.Rect{X: 53, Y: 19, W: 25, H: 11}) {
		t.Errorf("từ “nhập”: %+v", w)
	}
	// khung theo px ảnh gốc: nút “Mở rộng” ở hàng cuối (y ~ 315–345 trên trang)
	hits := vision.FindText(lines, vision.TextQuery{Text: "Mở rộng"})
	if len(hits) != 1 || hits[0].Box.Y < 310 || hits[0].Box.Y > 345 || hits[0].Box.X < 240 || hits[0].Box.X > 260 {
		t.Errorf("Mở rộng: %+v", hits)
	}
	if n := len(vision.FindText(lines, vision.TextQuery{Text: "Xác nhận"})); n != 2 {
		t.Errorf("Xác nhận thấy %d chỗ, cần 2 (dòng chữ + nút)", n)
	}
}
