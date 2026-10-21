package vision

import "image"

// So hai lần chụp một vùng (node Chờ màn hình, bước 42 G1): ảnh xám thu nhỏ rồi đếm % điểm ảnh đổi rõ.
// Thu nhỏ lấy trung bình theo ô nên nhiễu nén / khử răng cưa mất đi; đếm theo % điểm ảnh (không lấy độ lệch trung bình)
// để vùng nhấp nháy nhỏ (con trỏ chữ, đồng hồ giây) chỉ góp phần rất nhỏ, còn đổi lớn (trang tải, hộp thoại) ra % cao.

const (
	thumbMaxW = 320 // rộng tối đa ảnh thu nhỏ – vùng cả màn 1920 px: ô 6 × 6
	pixelTol  = 20  // điểm ảnh lệch quá mức này (thang 0..255) mới tính là đổi
)

// Thumb – ảnh xám thu nhỏ (rộng ≤ 320 px, giữ tỷ lệ) để so bằng DiffPercent.
func Thumb(img image.Image) Gray {
	g := ToGray(img)
	if g.w <= thumbMaxW {
		return g
	}
	return g.resize(thumbMaxW, g.h*thumbMaxW/g.w)
}

// DiffPercent – % điểm ảnh lệch > pixelTol giữa hai ảnh thu nhỏ (0..100). Khác cỡ (vùng đổi cỡ giữa chừng) = 100.
func DiffPercent(a, b Gray) float64 {
	if a.w != b.w || a.h != b.h {
		return 100
	}
	if a.empty() {
		return 0
	}
	n := 0
	for i, v := range a.p {
		if d := v - b.p[i]; d > pixelTol || d < -pixelTol {
			n++
		}
	}
	return float64(n) * 100 / float64(len(a.p))
}
