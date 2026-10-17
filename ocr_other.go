//go:build !darwin

package main

import (
	"errors"
	"image"
)

// readText: Windows (Windows.Media.Ocr) làm khi có máy Windows (bước 17 / 18).
func readText(image.Image) ([]TextLine, error) {
	return nil, errors.New("tìm theo chữ hiện chỉ có trên macOS")
}
