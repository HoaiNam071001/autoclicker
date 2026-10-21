package web

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/chromedp/chromedp"

	"ghosthand/internal/model"
)

// PageOutline (D4c, cho MCP 31): cây accessibility rút gọn của trang – AI đọc để biết trang có gì và viết selector
// cho node web, rẻ token hơn ảnh chụp. Không đụng chuột / phím, không đưa tab lên trước.

//go:embed outline.js
var outlineJS string

// maxOutline – số dòng tối đa (trang lớn: phần sau bị cắt, Outline.Total cho biết còn bao nhiêu).
const maxOutline = 300

// Outline – phần tử đang hiện của một tab theo thứ tự trang.
type Outline struct {
	URL         string        `json:"url"`
	Title       string        `json:"title"`
	ViewW       int           `json:"viewW"` // khung nhìn (CSS px)
	ViewH       int           `json:"viewH"`
	ScrollX     int           `json:"scrollX"`
	ScrollY     int           `json:"scrollY"`
	Items       []OutlineItem `json:"items"`
	Total       int           `json:"total"`       // số phần tử khớp trước khi cắt ở maxOutline
	CrossFrames int           `json:"crossFrames"` // iframe khác nguồn – không đọc được bên trong
}

// OutlineItem – một phần tử: vai trò, tên truy cập, selector đã kiểm (find tìm ra đúng nó) + trạng thái.
type OutlineItem struct {
	Role     string     `json:"role"`
	Name     string     `json:"name"`
	Selector string     `json:"sel"`
	Rect     model.Rect `json:"rect"`   // toạ độ client của trang chính (CSS px) – chưa đổi sang màn hình
	InView   bool       `json:"inView"` // tâm trong vùng nhìn thấy (không cần cuộn)
	Depth    int        `json:"depth"`  // số vùng (nav / main / form / dialog…) bao ngoài
	Level    int        `json:"level,omitempty"`
	Value    string     `json:"value,omitempty"`
	Secret   bool       `json:"secret,omitempty"`   // ô mật khẩu – không đọc giá trị
	Checked  string     `json:"checked,omitempty"`  // true | false | mixed
	Expanded string     `json:"expanded,omitempty"` // true | false
	Selected bool       `json:"selected,omitempty"`
	Disabled bool       `json:"disabled,omitempty"`
	Focused  bool       `json:"focused,omitempty"`
}

// outlineItem – dạng script trả (khung tách 4 số).
type outlineItem struct {
	OutlineItem
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type outlineResult struct {
	Outline
	Items []outlineItem `json:"items"`
}

// PageOutline đọc tab lọc theo tab (như Query.Tab, "" = tab đang dùng).
func (c *Chrome) PageOutline(ctx context.Context, tab string) (Outline, error) {
	t, _, err := c.tabFor(ctx, tab)
	if err != nil {
		return Outline{}, err
	}
	r, err := run(ctx, t.ctx, "đọc cấu trúc trang", chromedp.Evaluate[outlineResult](pageJS(outlineJS, strconv.Itoa(maxOutline))))
	if err != nil {
		return Outline{}, err
	}
	o := r.Outline
	o.Items = make([]OutlineItem, len(r.Items))
	for i, it := range r.Items {
		o.Items[i] = it.OutlineItem
		o.Items[i].Rect = model.Rect{X: it.X, Y: it.Y, W: it.W, H: it.H}
	}
	return o, nil
}

// Text – bản chữ thụt lề theo vùng, mỗi dòng: vai trò "tên" = giá trị [trạng thái] selector @ khung.
func (o Outline) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s – %s (khung nhìn %d×%d, đã cuộn %d,%d)\n", o.Title, o.URL, o.ViewW, o.ViewH, o.ScrollX, o.ScrollY)
	for _, it := range o.Items {
		b.WriteString(strings.Repeat("  ", it.Depth))
		b.WriteString(it.line())
		b.WriteByte('\n')
	}
	if n := o.Total - len(o.Items); n > 0 {
		fmt.Fprintf(&b, "… còn %d phần tử (cắt ở %d dòng)\n", n, len(o.Items))
	}
	if o.CrossFrames > 0 {
		fmt.Fprintf(&b, "(%d iframe khác nguồn – không đọc được bên trong)\n", o.CrossFrames)
	}
	return b.String()
}

func (it OutlineItem) line() string {
	s := it.Role
	if it.Level > 0 {
		s += strconv.Itoa(it.Level)
	}
	if it.Name != "" {
		s += " " + strconv.Quote(it.Name)
	}
	if it.Secret {
		s += " = ••••"
	} else if it.Value != "" {
		s += " = " + strconv.Quote(it.Value)
	}
	var st []string
	add := func(on bool, w string) {
		if on {
			st = append(st, w)
		}
	}
	add(it.Checked == "true", "đã chọn")
	add(it.Checked == "false", "chưa chọn")
	add(it.Checked == "mixed", "chọn một phần")
	add(it.Expanded == "true", "đang mở")
	add(it.Expanded == "false", "đang đóng")
	add(it.Selected, "đang chọn")
	add(it.Disabled, "tắt")
	add(it.Focused, "đang focus")
	add(!it.InView, "ngoài khung nhìn")
	if len(st) > 0 {
		s += " [" + strings.Join(st, ", ") + "]"
	}
	r := it.Rect
	return fmt.Sprintf("%s  %s  @%d,%d %d×%d", s, it.Selector, r.X, r.Y, r.W, r.H)
}
