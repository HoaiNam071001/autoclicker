// Package display – giao diện hệ điều hành cho node Màn hình / âm thanh (bước 43 S7): chế độ sáng / tối, độ sáng màn
// hình, chọn loa ra / micro vào mặc định, phát âm báo. Âm lượng giữ ở node Hệ thống (35).
package display

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrUnsupported = errors.New("tính năng này chưa có trên hệ điều hành này")
	// ErrNoBrightness – màn hình không cho chỉnh độ sáng bằng phần mềm (màn ngoài, máy bàn).
	ErrNoBrightness = errors.New("màn hình này không cho chỉnh độ sáng bằng phần mềm (thường chỉ màn hình laptop)")
)

// Audio – một thiết bị âm thanh.
type Audio struct {
	ID      string // mã bền của hệ điều hành (UID / endpoint)
	Name    string
	Input   bool // true = micro, false = loa
	Default bool // đang là mặc định (Windows: chưa đọc được → false)
}

// Backend – việc thật của hệ điều hành. Test thay bản giả.
type Backend interface {
	Dark() (bool, error)
	SetDark(dark bool) error
	Brightness() (percent int, err error)
	SetBrightness(percent int) error
	AudioDevices() ([]Audio, error)
	SetDefaultAudio(id string, input bool) error
	// PlaySound phát file âm thanh hoặc tên âm hệ thống (Glass, Ping, Beep…); wait: chờ phát xong.
	PlaySound(ctx context.Context, target string, wait bool) error
}

// FindAudio – chỉ số thiết bị khớp query (mã đúng, hoặc tên chứa query không phân biệt hoa thường) trong loại input;
// -1 = không có. Khớp tên đúng được ưu tiên hơn khớp chứa.
func FindAudio(list []Audio, query string, input bool) int {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return -1
	}
	contains := -1
	for i, a := range list {
		if a.Input != input {
			continue
		}
		name := strings.ToLower(a.Name)
		if a.ID == strings.TrimSpace(query) || name == q {
			return i
		}
		if contains < 0 && strings.Contains(name, q) {
			contains = i
		}
	}
	return contains
}
