package web

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ghosthand/internal/model"
)

const outlinePage = `<!doctype html><html><head><meta charset="utf-8"><title>Sân thử outline</title></head><body>
<nav aria-label="Menu chính"><a href="#a">Trang chủ</a> <a href="#b">Trang chủ cũ</a></nav>
<main><h1>Đăng nhập</h1>
<form aria-label="Đăng nhập"><label>Email <input id="email" value="a@b.vn"></label>
<label>Mật khẩu <input type="password" value="bimat"></label>
<label><input type="checkbox" checked> Nhớ tôi</label>
<button>Gửi</button><button>Gửi</button><button disabled data-testid="huy">Huỷ</button></form>
<button style="display:none">Ẩn</button>
<div style="height:3000px"></div><button>Cuối trang</button></main>
<iframe srcdoc='<button>Trong khung</button>'></iframe>
</body></html>`

// Cây rút gọn trên Chrome thật: đủ phần tử đang hiện, bỏ phần tử ẩn, không lộ mật khẩu, và mọi selector tìm lại
// đúng phần tử của dòng đó (kể cả hai nút trùng tên, nút trong iframe). Chỉ khi GHOSTHAND_REAL=1.
func TestRealPageOutline(t *testing.T) {
	if os.Getenv("GHOSTHAND_REAL") != "1" {
		t.Skip("đặt GHOSTHAND_REAL=1 để thử với Chrome thật")
	}
	dir := t.TempDir()
	page := filepath.Join(dir, "outline.html")
	if err := os.WriteFile(page, []byte(outlinePage), 0o644); err != nil {
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
	o, err := c.PageOutline(ctx, "outline.html")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n" + o.Text())
	byName := map[string]OutlineItem{}
	for _, it := range o.Items {
		byName[it.Role+" "+it.Name] = it
	}
	if _, ok := byName["button Ẩn"]; ok {
		t.Error("có nút đang ẩn")
	}
	if pw := byName["textbox Mật khẩu"]; !pw.Secret || pw.Value != "" {
		t.Errorf("ô mật khẩu: %+v", pw)
	}
	if em := byName["textbox Email"]; em.Value != "a@b.vn" || em.Selector != `role=textbox[name="Email"]` || em.Depth != 2 {
		t.Errorf("ô email: %+v", em)
	}
	if cb := byName["checkbox Nhớ tôi"]; cb.Checked != "true" {
		t.Errorf("ô chọn: %+v", cb)
	}
	if end := byName["button Cuối trang"]; end.InView {
		t.Errorf("nút cuối trang phải ngoài khung nhìn: %+v", end)
	}
	if _, ok := byName["button Trong khung"]; !ok {
		t.Error("thiếu nút trong iframe")
	}
	tab, _, err := c.tabFor(ctx, "outline.html")
	if err != nil {
		t.Fatal(err)
	}
	for i, it := range o.Items {
		sel, err := model.ParseSelector(it.Selector)
		if err != nil {
			t.Errorf("dòng %d %q: %v", i, it.Selector, err)
			continue
		}
		r, err := runFind(ctx, tab.ctx, Query{Selectors: []model.Selector{sel}, CheckOnly: true})
		if err != nil || r.Sel < 0 {
			t.Errorf("dòng %d %q không tìm thấy: %v", i, it.Selector, err)
			continue
		}
		// khung find.js trả phải trùng khung của dòng này → đúng phần tử
		got := model.Rect{X: int(r.X + .5), Y: int(r.Y + .5), W: int(r.W + .5), H: int(r.H + .5)}
		if abs(got.X-it.Rect.X) > 1 || abs(got.Y-it.Rect.Y) > 1 {
			t.Errorf("dòng %d %s %q ra phần tử khác: %+v ≠ %+v", i, it.Role, it.Selector, got, it.Rect)
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
