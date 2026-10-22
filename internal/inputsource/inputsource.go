// Package inputsource đọc / đổi bộ gõ – bàn phím đang dùng cho node Bộ gõ (bước 43 S6). macOS: Text Input Source
// (HIToolbox); Windows: layout bàn phím của cửa sổ đang ở trên cùng. Test thay Backend bản giả.
package inputsource

import (
	"slices"
	"strings"
)

// Source – một bộ gõ / layout.
type Source struct {
	ID   string // macOS com.apple.keylayout.ABC · Windows 04090409 (HKL)
	Name string // tên hiện cho người dùng: "ABC", "Telex" · Windows "en-US"
	Lang string // mã ngôn ngữ chính: en, vi…
}

// Backend – thao tác với bộ gõ của máy.
type Backend interface {
	Current() (Source, error)
	List() ([]Source, error)
	Select(id string) error
}

// System – Backend thật (input_darwin.go / input_windows.go; hệ khác báo không hỗ trợ).
type System struct{}

// Find – bộ gõ đầu tiên khớp want: trùng đúng ID / tên (không phân biệt hoa thường) trước, rồi mới chứa. -1 = không có.
func Find(list []Source, want string) int {
	want = strings.ToLower(strings.TrimSpace(want))
	if want == "" {
		return -1
	}
	for i, s := range list {
		if strings.ToLower(s.ID) == want || strings.ToLower(s.Name) == want {
			return i
		}
	}
	for i, s := range list {
		if strings.Contains(strings.ToLower(s.ID), want) || strings.Contains(strings.ToLower(s.Name), want) {
			return i
		}
	}
	return -1
}

// IsVietnamese – bộ gõ hệ điều hành là tiếng Việt (Telex / VNI của macOS, layout vi-VN của Windows).
func IsVietnamese(s Source) bool {
	return s.Lang == "vi" || strings.Contains(strings.ToLower(s.ID+" "+s.Name), "vietnam")
}

// imeApps – phần mềm gõ tiếng Việt bên thứ ba (đứng ngoài hệ điều hành, tự bật / tắt bằng phím tắt riêng).
var imeApps = []string{"unikey", "evkey", "openkey", "gotiengviet", "gõ tiếng việt"}

// FindIME – tên tiến trình đầu tiên trong names là phần mềm gõ tiếng Việt; "" = không có.
func FindIME(names []string) string {
	for _, n := range names {
		low := strings.ToLower(n)
		if slices.ContainsFunc(imeApps, func(a string) bool { return strings.Contains(low, a) }) {
			return n
		}
	}
	return ""
}
