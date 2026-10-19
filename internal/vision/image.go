package vision

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"image/png"
	"strconv"
	"strings"
)

// Mã hoá PNG / base64, màu "#RRGGBB", thu nhỏ ảnh – dùng chung cho chụp màn hình, engine, lịch sử chạy.

func EncodePNG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodePNG(s string) (image.Image, error) {
	b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, "data:image/png;base64,"))
	if err != nil {
		return nil, fmt.Errorf("ảnh mẫu hỏng: %w", err)
	}
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("ảnh mẫu hỏng: %w", err)
	}
	return img, nil
}

// ParseColor đọc "#RRGGBB".
func ParseColor(s string) (color.RGBA, error) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(h) != 6 {
		return color.RGBA{}, fmt.Errorf("màu không hợp lệ %q (cần #RRGGBB)", s)
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("màu không hợp lệ %q (cần #RRGGBB)", s)
	}
	return color.RGBA{uint8(v >> 16), uint8(v >> 8), uint8(v), 255}, nil
}

func HexColor(c color.RGBA) string { return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B) }

// ColorDiff: chênh lệch lớn nhất giữa các kênh màu.
func ColorDiff(a, b color.RGBA) int {
	d := func(x, y uint8) int { return int(max(x, y) - min(x, y)) }
	return max(d(a.R, b.R), d(a.G, b.G), d(a.B, b.B))
}

// Shrink thu nhỏ ảnh về rộng tối đa w (lấy trung bình từng ô – đủ để xem lại, không cần đẹp).
func Shrink(img image.Image, w int) image.Image {
	b := img.Bounds()
	k := (b.Dx() + w - 1) / w
	if k <= 1 {
		return img
	}
	out := image.NewRGBA(image.Rect(0, 0, b.Dx()/k, b.Dy()/k))
	for y := 0; y < out.Rect.Dy(); y++ {
		for x := 0; x < out.Rect.Dx(); x++ {
			var r, g, bl uint32
			for dy := 0; dy < k; dy++ {
				for dx := 0; dx < k; dx++ {
					cr, cg, cb, _ := img.At(b.Min.X+x*k+dx, b.Min.Y+y*k+dy).RGBA()
					r, g, bl = r+cr, g+cg, bl+cb
				}
			}
			n := uint32(k * k)
			i := out.PixOffset(x, y)
			out.Pix[i], out.Pix[i+1], out.Pix[i+2], out.Pix[i+3] = uint8(r/n>>8), uint8(g/n>>8), uint8(bl/n>>8), 255
		}
	}
	return out
}

// PixelHash – dấu vân tay điểm ảnh (FNV-1a): hai ảnh cùng cỡ cùng hash thì coi như y hệt – đọc chữ lại ảnh y hệt
// là phí (OCR ~100 ms mỗi lần), chờ giá trị đổi chỉ OCR khi vùng đổi.
func PixelHash(img image.Image) uint64 {
	h := fnv.New64a()
	b := img.Bounds()
	fmt.Fprintf(h, "%d×%d|", b.Dx(), b.Dy())
	if m, ok := img.(*image.RGBA); ok {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			o := m.PixOffset(b.Min.X, y)
			h.Write(m.Pix[o : o+4*b.Dx()])
		}
		return h.Sum64()
	}
	px := make([]byte, 4)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			px[0], px[1], px[2], px[3] = c.R, c.G, c.B, c.A
			h.Write(px)
		}
	}
	return h.Sum64()
}
