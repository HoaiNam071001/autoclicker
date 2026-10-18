package device

import (
	"image"

	"autoclicker/internal/vision"
)

// Tìm theo chữ (bước Chờ ảnh / Click theo ảnh với FindBy = text): đọc chữ trên ảnh chụp bằng bộ nhận chữ của hệ điều hành
// (macOS: Vision, ocr_darwin.go) rồi tìm chuỗi cần chờ. Bền hơn ảnh mẫu khi đổi theme / zoom / hover / tỷ lệ màn.

// TextReader đọc chữ trên ảnh – tách interface để test engine bằng bộ đọc giả.
type TextReader interface {
	ReadText(img image.Image) ([]vision.TextLine, error)
}

type SystemOCR struct{}

func (SystemOCR) ReadText(img image.Image) ([]vision.TextLine, error) { return readText(img) }
