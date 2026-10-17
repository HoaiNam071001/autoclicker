package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

// Chụp màn hình / đọc màu pixel cho bước chờ ảnh, click theo ảnh, chờ pixel (bước 25). Tách interface để test engine
// bằng màn hình giả. Toạ độ theo hệ robotgo (macOS: point, Windows: px thật); ảnh chụp theo px thật nên có thể lớn hơn
// vùng chụp (Retina 2x) – tỷ lệ = cỡ ảnh / cỡ vùng.
type Screen interface {
	Capture(r Rect) (image.Image, error)
	Pixel(x, y int) (color.RGBA, error)
}

type robotgoScreen struct{}

// Capture chụp vùng r trên màn hình chính (màn phụ – bước 17). Ưu tiên cách chụp không dính cửa sổ app (thanh mini
// đang chạy, overlay) và con trỏ chuột – không thì thanh mini che mất chỗ cần tìm.
func (robotgoScreen) Capture(r Rect) (image.Image, error) {
	if err := screenCaptureCheck(); err != nil {
		return nil, err
	}
	if img, ok, err := captureNoSelf(r); ok {
		return img, err
	}
	img, err := robotgo.CaptureImg(r.X, r.Y, r.W, r.H)
	if err != nil {
		return nil, fmt.Errorf("không chụp được màn hình: %w", err)
	}
	return img, nil
}

// Pixel đọc màu tại (x, y) – chụp 1×1 cùng cách như Capture (không dính thanh mini / con trỏ chuột).
func (s robotgoScreen) Pixel(x, y int) (color.RGBA, error) {
	if err := screenCaptureCheck(); err != nil {
		return color.RGBA{}, err
	}
	if img, ok, err := captureNoSelf(Rect{X: x, Y: y, W: 1, H: 1}); ok {
		if err != nil {
			return color.RGBA{}, err
		}
		r, g, b, _ := img.At(img.Bounds().Min.X, img.Bounds().Min.Y).RGBA()
		return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}, nil
	}
	return parseColor("#" + robotgo.GetPixelColor(x, y))
}

// captureArea kiểm tra vùng chụp nằm trong màn hình chính (robotgo chưa chụp được màn phụ).
func captureArea(r, main Rect) error {
	if r.W <= 0 || r.H <= 0 {
		return errors.New("vùng tìm rỗng")
	}
	if main.W > 0 && (r.X < main.X || r.Y < main.Y || r.X+r.W > main.X+main.W || r.Y+r.H > main.Y+main.H) {
		return errors.New("vùng tìm nằm ngoài màn hình chính – chưa hỗ trợ tìm ảnh trên màn hình phụ")
	}
	return nil
}

// imageScale: px ảnh chụp trên 1 đơn vị toạ độ.
func imageScale(img image.Image, r Rect) float64 {
	if r.W == 0 {
		return 1
	}
	return float64(img.Bounds().Dx()) / float64(r.W)
}

func encodePNG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func decodePNG(s string) (image.Image, error) {
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

// parseColor đọc "#RRGGBB".
func parseColor(s string) (color.RGBA, error) {
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

func hexColor(c color.RGBA) string { return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B) }

// colorDiff: chênh lệch lớn nhất giữa các kênh màu.
func colorDiff(a, b color.RGBA) int {
	d := func(x, y uint8) int { return abs(int(x) - int(y)) }
	return max(d(a.R, b.R), d(a.G, b.G), d(a.B, b.B))
}
