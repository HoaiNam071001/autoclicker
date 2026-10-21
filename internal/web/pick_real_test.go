package web

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

const pickPage = `<!doctype html><html><head><meta charset="utf-8"></head><body>
<h1>Bấm Đồng ý bên dưới</h1><button id="ok-2">Đồng ý</button>
<button id="send" aria-label="Gửi đơn hàng"><span>Gửi đơn</span></button>
<p>Gửi đơn ở đây</p><div id=":r7:"><a href="#x" class="x1n2onr6">Xem lịch sử</a></div>
<p>Bấm Xem lịch sử ở trên</p>
<iframe srcdoc='<div role="checkbox" id="agree">Tôi đồng ý</div>'></iframe>
</body></html>`

// Chọn phần tử trên trang bằng Chrome thật (click giả lập bằng JS – không cần chuột): selector sinh ra phải tìm lại
// đúng phần tử, ưu tiên chữ / aria / role, bỏ id tự sinh. Chỉ khi GHOSTHAND_REAL=1.
func TestRealPick(t *testing.T) {
	if os.Getenv("GHOSTHAND_REAL") != "1" {
		t.Skip("đặt GHOSTHAND_REAL=1 để thử với Chrome thật")
	}
	dir := t.TempDir()
	page := filepath.Join(dir, "pick.html")
	if err := os.WriteFile(page, []byte(pickPage), 0o644); err != nil {
		t.Fatal(err)
	}
	c := NewChrome(filepath.Join(os.TempDir(), "ghosthand-test-browser"), func(int, int) error { return nil })
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.Open(ctx, "file://"+page); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second) // trang tải xong
	cases := []struct {
		click string // JS trả phần tử để click
		want  []string
	}{
		{`document.querySelector('#send span')`, []string{"text=Gửi đơn", `aria=Gửi đơn hàng`, `role=button[name="Gửi đơn hàng"]`}},
		{`document.querySelector('h1 + button')`, []string{`text="Đồng ý"`, `role=button[name="Đồng ý"]`, "css=#ok-2"}},
		{`document.querySelector('a')`, []string{"text=Xem lịch sử", `role=link[name="Xem lịch sử"]`, "css=body > div > a"}},
		{`document.querySelector('iframe').contentDocument.getElementById('agree')`, []string{"text=Tôi đồng ý", `role=checkbox[name="Tôi đồng ý"]`, "css=#agree"}},
	}
	for _, tc := range cases {
		t.Run(tc.want[0], func(t *testing.T) {
			tab, _, err := c.tabFor(ctx, "pick.html")
			if err != nil {
				t.Fatal(err)
			}
			go func() {
				time.Sleep(500 * time.Millisecond)
				js := `(() => { (` + tc.click + `).dispatchEvent(new MouseEvent('click', {bubbles: true, cancelable: true})); return true; })()`
				_, _ = run(ctx, tab.ctx, "click", chromedp.Evaluate[bool](js))
			}()
			p, err := c.Pick(ctx, "pick.html")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(p.Selectors, tc.want) {
				t.Errorf("có %q, muốn %q", p.Selectors, tc.want)
			}
		})
	}
}
