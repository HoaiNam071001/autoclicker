package device

import "testing"

func TestWindowMatches(t *testing.T) {
	const title = "Phần tử màn hình – Sân thử - Google Chrome"
	cases := []struct {
		app  string
		want bool
	}{
		{"chrome", true}, {"Chrome.exe", true}, {"Google Chrome", true}, // tên tệp chạy / tên hiển thị (kịch bản soạn trên macOS)
		{title, true}, {"Microsoft Edge", false}, {"chrom", false}, {"Sân thử", false},
	}
	for _, c := range cases {
		if got := windowMatches(c.app, "chrome", title); got != c.want {
			t.Errorf("windowMatches(%q) = %v, cần %v", c.app, got, c.want)
		}
	}
}
