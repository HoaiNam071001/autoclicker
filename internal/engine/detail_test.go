package engine

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/element"
	"autoclicker/internal/model"
	"autoclicker/internal/web"
)

// detailJSON – detail của lỗi dạng JSON (như AI nhận qua wait_run).
func detailJSON(t *testing.T, err error) string {
	t.Helper()
	b, jerr := json.Marshal(errDetail(err))
	if jerr != nil {
		t.Fatal(jerr)
	}
	return string(b)
}

func wantAll(t *testing.T, name, got string, parts ...string) {
	t.Helper()
	for _, p := range parts {
		if !strings.Contains(got, p) {
			t.Errorf("%s: detail %s thiếu %s", name, got, p)
		}
	}
}

// Không thấy phần tử → ứng viên gần đúng (vai trò, tên, khung, % giống).
func TestElementMissDetail(t *testing.T) {
	e := newEngine(&devicetest.FakeInput{}, model.Script{}, nil)
	e.elems = &devicetest.FakeElements{Frames: [][]element.Node{page(elemBox(10, 20, 80, 30))}}
	_, err := e.findImage(context.Background(), elementStep(model.StepWaitImage, func(s *model.Step) {
		s.Selectors, s.TimeoutMs = []string{"text=Mua ngey"}, 5
	}))
	wantAll(t, "phần tử", detailJSON(t, err), `"kind":"element"`, `"read":2`, `"count":0`,
		`"near":[{"role":"button","name":"Mua ngay","box":{"x":10,"y":20,"w":80,"h":30`)
}

// Không thấy cửa sổ → cửa sổ gần giống; app không lên cửa sổ → tiến trình đang chạy + app tên gần giống.
func TestWindowAppMissDetail(t *testing.T) {
	e, _ := winEngine()
	_, err := e.windowStep(t.Context(), winStep(model.WinWaitOpen, func(s *model.Step) { s.WinApp, s.WinTitle, s.TimeoutMs = "chromium", "San thu", 0 }))
	wantAll(t, "cửa sổ", detailJSON(t, err), `"kind":"window"`, `"title":"Sân thử – Google Chrome"`, `"score":`)
	if d := detailJSON(t, err); strings.Contains(d, `"Zalo"`) {
		t.Errorf("Zalo không giống mà vẫn là ứng viên: %s", d)
	}

	a, _ := appEngine(false)
	_, err = a.appStep(t.Context(), appStep(model.AppWaitStart, "notepad", func(s *model.Step) { s.TimeoutMs = 0 }))
	wantAll(t, "app chưa có cửa sổ", detailJSON(t, err), `"kind":"app"`, `"running":[{"name":"Notepad","id":"","pid":20,"windows":0}]`)
	_, err = a.appStep(t.Context(), appStep(model.AppWaitStart, "chromium", func(s *model.Step) { s.TimeoutMs = 0 }))
	wantAll(t, "app sai tên", detailJSON(t, err), `"running":[]`, `"near":[{"name":"chrome","id":"C:/Chrome/chrome.exe"`)
}

// Web hết giờ → trạng thái lần tìm cuối có cấu trúc.
func TestWebMissDetail(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{{Selector: 0, Covered: "div#popup", Count: 1, Tab: "file:///a.html"}}}
	_, res, _ := runWeb(t, b, webNode(model.StepWebClick, nil))
	d, _ := json.Marshal(res[0].Error.Detail)
	wantAll(t, "web", string(d), `"kind":"web"`, `"matched":"text=Xem lịch sử"`, `"covered":"div#popup"`, `"count":1`, `"tab":"file:///a.html"`)
}

// Đọc phần tử không thấy → cùng detail element (tính lúc hết giờ).
func TestElementReadMissDetail(t *testing.T) {
	e := newEngine(&devicetest.FakeInput{}, model.Script{}, nil)
	e.elems = &devicetest.FakeElements{Frames: [][]element.Node{page(elemBox(10, 20, 80, 30))}}
	_, err := e.elementRead(context.Background(), elementStep(model.StepElementRead, func(s *model.Step) {
		s.Selectors, s.TimeoutMs = []string{"text=Mua ngey"}, 5
	}))
	wantAll(t, "đọc phần tử", detailJSON(t, err), `"kind":"element"`, `"name":"Mua ngay"`)
}
