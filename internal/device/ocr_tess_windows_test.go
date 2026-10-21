package device

import (
	"image"
	"image/png"
	"os"
	"testing"
	"time"

	"autoclicker/internal/vision"
)

// Trang thử 40 O0 (chụp bằng Chrome headless, DPI 100 % / 125 %): chữ 11–20 px thường / đậm + nút chữ trắng nền xanh,
// cột trái nền sáng, cột phải nền tối. Tesseract phải đọc đúng cả dấu mọi chuỗi ở cả hai cột (không so gần đúng).
var viTruth = []string{
	"Đăng nhập bằng tài khoản của bạn", "Mật khẩu không đúng, vui lòng thử lại", "Lưu trữ dữ liệu",
	"Tiếp tục thanh toán đơn hàng", "Xoá lịch sử chạy gần đây", "Người dùng đã được cập nhật thành công",
	"Chọn tệp ảnh để tải lên máy chủ", "Hủy bỏ thao tác hiện tại", "Kịch bản tự động đọc chữ tiếng Việt",
	"Bấm vào nút Xác nhận ở góc phải", "Lưu", "Đóng", "Thử lại", "Mở rộng", "Xác nhận",
}

func TestTesseractVietnamese(t *testing.T) {
	if _, err := tessBundleDir(); err != nil {
		t.Skip(err)
	}
	for _, c := range []struct {
		file string
		h    int // chiều cao phần trang thử (px)
	}{{"testdata/ocr-vi.png", 372}, {"testdata/ocr-vi-125.png", 465}} {
		f, err := os.Open(c.file)
		if err != nil {
			t.Fatal(err)
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		w := img.Bounds().Dx()
		for _, half := range []image.Rectangle{image.Rect(0, 0, w/2, c.h), image.Rect(w/2, 0, w, c.h)} {
			start := time.Now()
			lines, err := tessRead(img.(interface {
				SubImage(image.Rectangle) image.Image
			}).SubImage(half))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s %v: %d dòng, %v", c.file, half, len(lines), time.Since(start))
			for _, want := range viTruth {
				hits := vision.FindText(lines, vision.TextQuery{Text: want})
				if len(hits) == 0 {
					t.Errorf("%s %v: không thấy %q", c.file, half, want)
					continue
				}
				if b := hits[0].Box; b.X < 0 || b.X+b.W > half.Dx() || b.Y < 0 || b.Y+b.H > half.Dy() {
					t.Errorf("%q: khung %v ra ngoài ảnh %v", want, b, half.Size())
				}
			}
		}
	}
}
