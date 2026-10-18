package web

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"

	"autoclicker/internal/model"
)

const bgPage = `<!doctype html><html><head><meta charset="utf-8"></head><body style="height:3000px">
<form id="f"><input id="name" value="chữ cũ"></form>
<button id="b">Bấm</button>
<div id="box" style="height:100px;width:200px;overflow:auto"><div style="height:2000px">dài</div></div>
<script>
window.log = {clicks: 0, dbl: 0, trusted: true, submit: 0};
const b = document.getElementById('b');
b.addEventListener('click', e => { log.clicks++; log.trusted = log.trusted && e.isTrusted; });
b.addEventListener('dblclick', () => log.dbl++);
document.getElementById('f').addEventListener('submit', e => { e.preventDefault(); log.submit++; });
</script></body></html>`

// Chạy nền với Chrome thật đang thu nhỏ: gõ, bấm, double-click, lăn qua CDP – chuột thật không được động tới.
// Chỉ khi AUTOCLICKER_REAL=1.
func TestRealBackground(t *testing.T) {
	if os.Getenv("AUTOCLICKER_REAL") != "1" {
		t.Skip("đặt AUTOCLICKER_REAL=1 để thử với Chrome thật")
	}
	page := filepath.Join(t.TempDir(), "bg.html")
	if err := os.WriteFile(page, []byte(bgPage), 0o644); err != nil {
		t.Fatal(err)
	}
	moved := 0
	c := NewChrome(filepath.Join(os.TempDir(), "autoclicker-test-browser"), func(int, int) error { moved++; return nil })
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.Open(ctx, "file://"+page); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	tab, _, err := c.tabFor(ctx, "bg.html")
	if err != nil {
		t.Fatal(err)
	}
	setWindow(t, ctx, tab.ctx, browser.WindowStateMinimized)
	defer setWindow(t, ctx, tab.ctx, browser.WindowStateNormal)
	time.Sleep(500 * time.Millisecond)

	find := func(css string) Found {
		f, err := c.Find(ctx, Query{Selectors: []model.Selector{{By: model.SelCSS, Value: css}}, Scroll: true, Background: true})
		if err != nil || !f.OK {
			t.Fatalf("%s: %+v %v", css, f, err)
		}
		return f
	}
	center := func(f Found) (int, int) { return f.Rect.X + f.Rect.W/2, f.Rect.Y + f.Rect.H/2 }
	eval := func(js string) any {
		v, err := run(ctx, tab.ctx, "đọc", chromedp.Evaluate[any](js))
		if err != nil {
			t.Fatal(err)
		}
		return v
	}

	x, y := center(find("#name"))
	if err := c.BgClick(ctx, x, y, "left", 1); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if ok, at, err := c.Focused(ctx); !ok || err != nil {
		t.Fatalf("ô chưa nhận focus: %s %v", at, err)
	}
	if err := c.BgType(ctx, "Xin chào {1}", true, true); err != nil {
		t.Fatal(err)
	}
	if v := eval(`document.getElementById('name').value`); v != "Xin chào {1}" {
		t.Errorf("giá trị ô %q", v)
	}

	x, y = center(find("#b"))
	for i := 1; i <= 2; i++ {
		if err := c.BgClick(ctx, x, y, "left", i); err != nil {
			t.Fatal(err)
		}
	}
	x, y = center(find("#box"))
	for range 3 {
		if err := c.BgWheel(ctx, x, y, "down"); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(300 * time.Millisecond)
	got := eval(`JSON.stringify({...log, top: document.getElementById('box').scrollTop})`)
	if want := `{"clicks":2,"dbl":1,"trusted":true,"submit":1,"top":300}`; got != want {
		t.Errorf("trang ghi %v, muốn %s", got, want)
	}
	if moved != 0 {
		t.Errorf("chạy nền mà chuột thật di %d lần", moved)
	}
}

func setWindow(t *testing.T, ctx, tctx context.Context, state browser.WindowState) {
	_, err := run(ctx, tctx, "đổi cửa sổ", func(ctx context.Context, tg *chromedp.Target) (cdp.Empty, error) {
		w, err := cdp.Call(ctx, tg, browser.GetWindowForTarget, browser.GetWindowForTargetParams{})
		if err != nil {
			return cdp.Empty{}, err
		}
		return cdp.Call(ctx, tg, browser.SetWindowBounds, browser.SetWindowBoundsParams{WindowID: w.WindowID, Bounds: &browser.Bounds{WindowState: state}})
	})
	if err != nil {
		t.Fatal(err)
	}
}
