package web

import (
	"testing"

	"ghosthand/internal/model"
)

func TestOutlineText(t *testing.T) {
	o := Outline{URL: "https://vd.vn/dang-nhap", Title: "Đăng nhập", ViewW: 1280, ViewH: 720, ScrollY: 40, Total: 5, CrossFrames: 1,
		Items: []OutlineItem{
			{Role: "form", Name: "Đăng nhập", Selector: `role=form[name="Đăng nhập"]`, Rect: model.Rect{X: 10, Y: 20, W: 400, H: 300}, InView: true},
			{Role: "heading", Level: 1, Name: "Chào bạn", Selector: `role=heading[name="Chào bạn"]`, Depth: 1, InView: true},
			{Role: "textbox", Name: "Mật khẩu", Secret: true, Selector: "css=#pw", Depth: 1, InView: true, Focused: true},
			{Role: "checkbox", Name: "Nhớ tôi", Value: "", Checked: "false", Selector: "role=checkbox", Depth: 1, Disabled: true},
		}}
	want := `Đăng nhập – https://vd.vn/dang-nhap (khung nhìn 1280×720, đã cuộn 0,40)
form "Đăng nhập"  role=form[name="Đăng nhập"]  @10,20 400×300
  heading1 "Chào bạn"  role=heading[name="Chào bạn"]  @0,0 0×0
  textbox "Mật khẩu" = •••• [đang focus]  css=#pw  @0,0 0×0
  checkbox "Nhớ tôi" [chưa chọn, tắt, ngoài khung nhìn]  role=checkbox  @0,0 0×0
… còn 1 phần tử (cắt ở 4 dòng)
(1 iframe khác nguồn – không đọc được bên trong)
`
	if got := o.Text(); got != want {
		t.Errorf("có\n%s\nmuốn\n%s", got, want)
	}
}
