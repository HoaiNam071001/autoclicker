package web

import (
	"strings"
	"testing"
)

func TestClassifyPage(t *testing.T) {
	article := strings.Repeat("Bài viết dài nói về captcha và just a moment. ", 60)
	cases := []struct {
		name     string
		s        pageSigns
		elements int
		loaded   bool
		want     string
	}{
		{"trang thường", pageSigns{Title: "Wikipedia", Text: "Trang chính"}, 120, true, PageOK},
		{"reCAPTCHA hiện", pageSigns{Frames: []pageFrame{{Src: "https://www.google.com/recaptcha/api2/anchor?k=1", W: 304, H: 78, Shown: true}}}, 30, true, PageCaptcha},
		{"reCAPTCHA v3 ẩn", pageSigns{Frames: []pageFrame{{Src: "https://www.google.com/recaptcha/api2/anchor?k=1&size=invisible", W: 256, H: 60, Shown: true}}}, 30, true, PageOK},
		{"iframe captcha ẩn", pageSigns{Frames: []pageFrame{{Src: "https://hcaptcha.com/x", W: 0, H: 0}}}, 30, true, PageOK},
		{"Cloudflare", pageSigns{Title: "Just a moment...", Text: "Checking if the site connection is secure"}, 2, true, PageCaptcha},
		{"khối turnstile", pageSigns{Marks: []string{".cf-turnstile"}}, 10, true, PageCaptcha},
		{"bị chặn", pageSigns{Title: "Access Denied", Text: "You don't have permission"}, 1, true, PageBlocked},
		{"game canvas", pageSigns{Title: "Game A", Canvas: 0.85}, 3, true, PageCanvas},
		{"canvas nhỏ trong trang lớn", pageSigns{Title: "Báo", Canvas: 0.6}, 80, true, PageOK},
		{"trang trống", pageSigns{Text: "  "}, 0, false, PageEmpty},
		{"bài dài nhắc chữ captcha", pageSigns{Title: "Tin", Text: article}, 50, true, PageOK},
	}
	for _, c := range cases {
		got := classifyPage(c.s, c.elements, c.loaded)
		if got.Verdict != c.want || got.UseWeb != (c.want == PageOK) {
			t.Errorf("%s: verdict = %s (useWeb %v), muốn %s; dấu hiệu %v", c.name, got.Verdict, got.UseWeb, c.want, got.Signs)
		}
	}
}

func TestClassifyPageSigns(t *testing.T) {
	got := classifyPage(pageSigns{Frames: []pageFrame{{Src: "https://challenges.cloudflare.com/cdn-cgi/x?y", W: 300, H: 65, Shown: true}}}, 5, false)
	joined := strings.Join(got.Signs, " | ")
	for _, want := range []string{"challenges.cloudflare.com", "300×65", "chưa tải xong"} {
		if !strings.Contains(joined, want) {
			t.Errorf("dấu hiệu thiếu %q: %s", want, joined)
		}
	}
}
