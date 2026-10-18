package web

import (
	"strings"
	"testing"
)

func TestPickTab(t *testing.T) {
	tabs := []tabInfo{
		{ID: "a", Title: "Trang chủ", URL: "https://example.com/"},
		{ID: "b", Title: "Lịch sử đơn", URL: "https://shop.vn/orders"},
	}
	cases := []struct{ filter, prefer, want string }{
		{"", "", "a"},             // tab dùng gần nhất đứng đầu
		{"", "b", "b"},            // tab vừa mở bằng node Mở
		{"", "đã đóng", "a"},      // tab ưu tiên không còn
		{"SHOP.vn", "a", "b"},     // lọc theo URL, không phân biệt hoa thường
		{"lịch sử", "", "b"},      // lọc theo tiêu đề
		{"không có", "a", "(no)"}, // lọc không khớp → không chọn bừa
	}
	for _, c := range cases {
		got, ok := pickTab(tabs, c.filter, c.prefer)
		if !ok {
			got.ID = "(no)"
		}
		if got.ID != c.want {
			t.Errorf("pickTab(%q, %q) = %s, cần %s", c.filter, c.prefer, got.ID, c.want)
		}
	}
	if _, ok := pickTab(nil, "", ""); ok {
		t.Error("không có tab mà vẫn chọn được")
	}
}

func TestBrowserPaths(t *testing.T) {
	env := map[string]string{"HOME": "/Users/a", "ProgramFiles": `C:\Program Files`, "LocalAppData": `C:\Users\a\AppData\Local`}
	get := func(k string) string { return env[k] }
	mac := browserPaths("darwin", get)
	if mac[0] != "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" {
		t.Errorf("macOS: Chrome phải thử đầu tiên, có %q", mac[0])
	}
	win := strings.Join(browserPaths("windows", get), "\n")
	for _, want := range []string{`C:\Program Files\Google\Chrome\Application\chrome.exe`, `C:\Users\a\AppData\Local\Microsoft\Edge\Application\msedge.exe`} {
		if !strings.Contains(win, want) {
			t.Errorf("Windows thiếu %s", want)
		}
	}
	if strings.Contains(win, `\\`) || strings.HasPrefix(win, `\`) {
		t.Errorf("biến môi trường trống tạo đường dẫn hỏng:\n%s", win)
	}
}
