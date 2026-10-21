package web

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Thử trang trên Chrome thật (bước 44): trang thường / captcha giả / game canvas ra đúng kết luận. Chỉ khi
// GHOSTHAND_REAL=1.
func TestRealCheckPage(t *testing.T) {
	if os.Getenv("GHOSTHAND_REAL") != "1" {
		t.Skip("đặt GHOSTHAND_REAL=1 để thử với Chrome thật")
	}
	pages := map[string]struct{ html, want string }{
		"thuong.html": {`<!doctype html><meta charset="utf-8"><title>Cửa hàng</title><h1>Cửa hàng</h1><button>Mua</button><a href="#x">Giỏ hàng</a>`, PageOK},
		"captcha.html": {`<!doctype html><meta charset="utf-8"><title>Đăng ký</title><p>Điền form</p>
			<iframe title="reCAPTCHA" src="about:blank#recaptcha" style="width:304px;height:78px"></iframe>`, PageCaptcha},
		"canvas.html": {`<!doctype html><meta charset="utf-8"><title>Game A</title><style>body{margin:0}canvas{width:100vw;height:100vh;display:block}</style>
			<canvas></canvas>`, PageCanvas},
	}
	dir := t.TempDir()
	c := NewChrome(filepath.Join(os.TempDir(), "ghosthand-test-browser"), func(int, int) error { return nil })
	defer c.Close()
	for name, p := range pages {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(p.html), 0o644); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		pc, err := c.CheckPage(ctx, "file:///"+filepath.ToSlash(path))
		cancel()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		t.Logf("%s: %s, %d phần tử, %v", name, pc.Verdict, pc.Elements, pc.Signs)
		if pc.Verdict != p.want || !pc.Loaded {
			t.Errorf("%s: verdict = %s (loaded %v), muốn %s", name, pc.Verdict, pc.Loaded, p.want)
		}
	}
}
