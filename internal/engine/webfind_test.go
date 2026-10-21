package engine

import (
	"context"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/model"
	"ghosthand/internal/web"
)

// Thử tìm: một lần, không cuộn / không đòi đứng yên, node chờ mất vẫn lấy khung; không đụng chuột.
func TestFindElement(t *testing.T) {
	b := &fakeBrowser{finds: []web.Found{found(1, model.Rect{X: 10, Y: 20, W: 30, H: 40})}}
	in := &devicetest.FakeInput{}
	e := New(in, model.Script{}, Options{Browser: b})
	s := webNode(model.StepWebWait, func(s *model.Step) { s.WebGone = true })
	r, err := e.FindElement(context.Background(), s)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Found || r.Selector != "css=#history" || r.X != 10 || r.H != 40 || r.Text != "Xem lịch sử" {
		t.Errorf("kết quả %+v", r)
	}
	if q := b.queries[0]; q.Scroll || q.CheckOnly || b.n != 1 || len(in.Lines()) != 0 {
		t.Errorf("query %+v, tìm %d lần, chuột %q", q, b.n, in.Lines())
	}

	b.finds = []web.Found{{Selector: -1, Hidden: 2}}
	if r, _ = e.FindElement(context.Background(), s); r.Found || r.Selector != "" || r.Hidden != 2 {
		t.Errorf("không thấy: %+v", r)
	}
	s.Selectors = []string{" "}
	if _, err := e.FindElement(context.Background(), s); err == nil {
		t.Error("selector trống phải báo lỗi")
	}
	if _, err := New(in, model.Script{}, Options{}).FindElement(context.Background(), s); err == nil {
		t.Error("không có trình duyệt phải báo lỗi")
	}
}
