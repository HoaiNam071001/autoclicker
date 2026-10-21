package engine

import (
	"reflect"
	"strings"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/web"
)

func TestWebTypeClicksChecksFocusThenTypes(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{found(0, model.Rect{X: 100, Y: 200, W: 40, H: 20})}}
	in, res, err := runWeb(t, b, webNode(model.StepWebType, func(s *model.Step) { s.Text, s.WebClear, s.WebEnter = "ab", true, true }))
	if err != nil {
		t.Fatal(err)
	}
	mod := shortcutMod()
	want := []string{"move 120,210", "click left false", "down " + mod + " []", "down a [" + mod + "]", "up a", "up " + mod,
		"down backspace []", "up backspace", "type a", "type b", "down enter []", "up enter"}
	if got := in.Lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("có %q\nmuốn %q", got, want)
	}
	if b.nFocus != 1 || res[0].Output["text"] != "ab" || res[0].Output["clickX"] != 120 {
		t.Errorf("focus %d lần, output %v", b.nFocus, res[0].Output)
	}
}

// Click rồi mà ô không nhận focus (cách tìm trỏ nhầm nút) → không gõ, cổng fail not_found.
func TestWebTypeNoFocus(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{found(0, model.Rect{W: 10, H: 10})}, focus: []bool{false}}
	in, res, err := runWeb(t, b, webNode(model.StepWebType, func(s *model.Step) { s.Text = "x" }))
	if err == nil || !strings.Contains(err.Error(), "không nhận focus (đang focus: button “Gửi”)") || res[0].Error.Code != model.ErrNotFound {
		t.Errorf("lỗi %v", err)
	}
	for _, l := range in.Lines() {
		if strings.HasPrefix(l, "type") {
			t.Errorf("không được gõ khi chưa focus: %q", in.Lines())
		}
	}
}

func TestWebRead(t *testing.T) {
	f := found(1, model.Rect{})
	f.Value, f.Count = "https://shop.vn/don/7", 3
	b := &fakeBrowser{finds: []web.Found{missing, f}}
	in, res, err := runWeb(t, b, webNode(model.StepWebRead, func(s *model.Step) { s.ReadFrom, s.Attr = model.ReadAttr, " href " }))
	if err != nil {
		t.Fatal(err)
	}
	q := b.queries[0]
	if !q.CheckOnly || q.Scroll || q.Read != "attr" || q.Attr != "href" || len(in.Lines()) != 0 {
		t.Errorf("query %+v, chuột %q", q, in.Lines())
	}
	if want := map[string]any{"text": "https://shop.vn/don/7", "count": 3, "selector": "css=#history"}; !reflect.DeepEqual(ownOutput(res[0]), want) {
		t.Errorf("output %v", res[0].Output)
	}
	// mặc định đọc chữ hiển thị; không thấy → fail not_found như webWait
	b = &fakeBrowser{finds: []web.Found{missing}}
	_, res, err = runWeb(t, b, webNode(model.StepWebRead, nil))
	if b.queries[0].Read != "text" || err == nil || res[0].Error.Code != model.ErrNotFound {
		t.Errorf("read %q, lỗi %v", b.queries[0].Read, err)
	}
}

func TestWebScrollTo(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{found(0, model.Rect{X: 10, Y: 20, W: 30, H: 40})}}
	in, res, err := runWeb(t, b, webNode(model.StepWebScroll, nil))
	if err != nil {
		t.Fatal(err)
	}
	if !b.queries[0].Scroll || len(b.hovered) != 0 || len(in.Lines()) != 0 || res[0].Output["y"] != 40 {
		t.Errorf("query %+v, hover %q, chuột %q, output %v", b.queries[0], b.hovered, in.Lines(), res[0].Output)
	}
}

func TestWebScrollIn(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{found(0, model.Rect{X: 100, Y: 200, W: 40, H: 20})}}
	in, _, err := runWeb(t, b, webNode(model.StepWebScroll, func(s *model.Step) { s.ScrollMode, s.Direction, s.Amount = model.ScrollModeIn, "down", 2 }))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"move 120,210", "scroll down", "scroll down"}; !reflect.DeepEqual(in.Lines(), want) || len(b.hovered) != 1 {
		t.Errorf("có %q, hover %q", in.Lines(), b.hovered)
	}
}

// Tới khi thấy: lăn giữa trang (khung nhìn) tới khi đích hiện trong khung nhìn; thấy mà còn ngoài khung thì lăn tiếp.
func TestWebScrollUntil(t *testing.T) {
	view := web.Found{OK: true, Selector: 0, Rect: model.Rect{X: 0, Y: 100, W: 800, H: 600}}
	off := found(0, model.Rect{X: 10, Y: 900, W: 30, H: 40})
	off.OffView = true
	b := &fakeBrowser{finds: []web.Found{missing, view, off, view, found(0, model.Rect{X: 10, Y: 300, W: 30, H: 40})}}
	in, res, err := runWeb(t, b, webNode(model.StepWebScroll, func(s *model.Step) {
		s.ScrollMode, s.Direction, s.Amount, s.TimeoutMs = model.ScrollModeUntil, "down", 1, 2000
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !b.queries[1].Viewport || !b.queries[1].Scroll || b.queries[0].Scroll || b.queries[0].Viewport {
		t.Errorf("query đích %+v, khung %+v", b.queries[0], b.queries[1])
	}
	if want := []string{"move 400,400", "scroll down", "move 400,400", "scroll down"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("có %q", in.Lines())
	}
	if b.hovered[0] != "400,400 false" || res[0].Output["scrolls"] != 2 || res[0].Output["y"] != 320 {
		t.Errorf("hover %q, output %v", b.hovered, res[0].Output)
	}
	// khung cuộn riêng: tìm theo scrollIn; lăn mãi không thấy → báo số lần đã lăn
	b = &fakeBrowser{}
	for range 50 {
		b.finds = append(b.finds, missing, view)
	}
	_, _, err = runWeb(t, b, webNode(model.StepWebScroll, func(s *model.Step) {
		s.ScrollMode, s.Direction, s.Amount, s.ScrollIn = model.ScrollModeUntil, "down", 1, []string{"css=#feed"}
	}))
	if err == nil || !strings.Contains(err.Error(), "đã lăn") || b.queries[1].Viewport || b.queries[1].Selectors[0].Value != "#feed" {
		t.Errorf("lỗi %v, khung %+v", err, b.queries[1])
	}
}

// Ô nhận focus chậm (React state, hiệu ứng): hỏi lại vài lần thay vì báo lỗi ngay ở lần đầu.
func TestWebTypeLateFocus(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{found(0, model.Rect{W: 10, H: 10})}, focus: []bool{false, false, true}}
	_, res, err := runWeb(t, b, webNode(model.StepWebType, func(s *model.Step) { s.Text = "x" }))
	if err != nil || res[0].Output["text"] != "x" || b.nFocus != 3 {
		t.Errorf("lỗi %v, focus hỏi %d lần, output %v", err, b.nFocus, res[0].Output)
	}
}

// skipFocusCheck (46 A3): ô tuỳ biến không nhận focus kiểu thường vẫn gõ, không hỏi focus; scope / chỗ thứ N đi vào query.
func TestWebTypeSkipFocusCheck(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{found(0, model.Rect{W: 10, H: 10})}, focus: []bool{false}}
	_, res, err := runWeb(t, b, webNode(model.StepWebType, func(s *model.Step) {
		s.Text, s.SkipFocusCheck, s.WebScope, s.Occurrence = "x", true, model.WebScopeDialog, 2
	}))
	if err != nil || res[0].Output["text"] != "x" || b.nFocus != 0 {
		t.Errorf("lỗi %v, focus hỏi %d lần, output %v", err, b.nFocus, res[0].Output)
	}
	if q := b.queries[0]; q.Scope != model.WebScopeDialog || q.Occurrence != 2 {
		t.Errorf("query %+v", q)
	}
}

// Hết giờ khi chỉ tìm trong hộp thoại mà trang không có hộp thoại → lời báo nói rõ.
func TestWebNoDialogMessage(t *testing.T) {
	f := missing
	f.NoDialog = true
	b := &fakeBrowser{finds: []web.Found{f}}
	_, _, err := runWeb(t, b, webNode(model.StepWebClick, func(s *model.Step) { s.WebScope, s.TimeoutMs = model.WebScopeDialog, 50 }))
	if err == nil || !strings.Contains(err.Error(), "không có hộp thoại nào đang mở") {
		t.Errorf("lỗi %v", err)
	}
}

// webKey (46 B3): phím đi qua trình duyệt (CDP), không đụng bàn phím thật; output keys + tab.
func TestWebKey(t *testing.T) {
	b := &fakeBrowser{}
	in, res, err := runWeb(t, b, webNode(model.StepWebKey, func(s *model.Step) {
		s.Keys, s.Modifiers, s.Tab, s.KeyMode = []string{"tab", "enter"}, []string{"shift"}, "shop", model.KeySequence
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(in.Lines()) != 0 || len(b.bg) != 1 || b.bg[0] != "keys shop [shift] [tab enter] true" {
		t.Errorf("bàn phím thật %q, trình duyệt %q", in.Lines(), b.bg)
	}
	if res[0].Output["keys"] != "shift+tab+enter" || res[0].Output["tab"] != "https://shop.vn/" {
		t.Errorf("output %v", res[0].Output)
	}
}
