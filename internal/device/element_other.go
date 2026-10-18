//go:build !darwin

package device

import (
	"errors"

	"autoclicker/internal/element"
)

// readElements: Windows (UI Automation) làm ở bước 37 E5 khi có máy Windows.
func readElements(string) ([]element.Node, error) {
	return nil, errors.New("tìm phần tử màn hình hiện chỉ có trên macOS")
}
