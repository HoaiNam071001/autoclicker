//go:build !darwin && !windows

package device

import (
	"errors"
	"image"

	"ghosthand/internal/vision"
)

// readText: Windows (Windows.Media.Ocr) làm khi có máy Windows (bước 17 / 18).
func readText(image.Image) ([]vision.TextLine, error) {
	return nil, errors.New("tìm theo chữ hiện chỉ có trên macOS")
}
