package element

import (
	"reflect"
	"testing"

	"autoclicker/internal/model"
)

func TestAt(t *testing.T) {
	nodes := tree()
	cases := []struct {
		name string
		p    model.Point
		want int
	}{
		{"khối chữ trong nút – nhỏ nhất", model.Point{X: 150, Y: 215}, 3},
		{"mép nút ngoài khối chữ", model.Point{X: 102, Y: 202}, 2},
		{"nút icon", model.Point{X: 910, Y: 100}, 4},
		{"chỉ có cửa sổ", model.Point{X: 500, Y: 40}, 0},
		{"nút kẹp mép không chọn", model.Point{X: 150, Y: 790}, 1},
		{"ngoài cửa sổ", model.Point{X: 1200, Y: 40}, -1},
	}
	for _, c := range cases {
		if got := At(nodes, c.p); got != c.want {
			t.Errorf("%s: At = %d, muốn %d", c.name, got, c.want)
		}
	}
}

func TestQueryFor(t *testing.T) {
	nodes := tree()
	cases := []struct {
		name string
		i    int
		want Pick
	}{
		{"nút có chữ – nút đầu trong thứ tự đọc nên chỗ thứ 1", 2, Pick{Selectors: []string{`text="Thêm vào giỏ"`, `role=button[name="Thêm vào giỏ"]`, "role=button"}, Occurrence: 1}},
		{"nút icon – nhãn aria", 4, Pick{Selectors: []string{`aria="Đóng"`, `role=button[name="Đóng"]`}, Occurrence: 1}},
		{"ô nhập – chữ mờ", 6, Pick{Selectors: []string{`text="Tìm sản phẩm"`, `role=textbox[name="Tìm sản phẩm"]`, "role=textbox"}, Occurrence: 1}},
	}
	for _, c := range cases {
		if got := QueryFor(nodes, c.i); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: %+v, muốn %+v", c.name, got, c.want)
		}
	}
}

// Chữ trùng nhiều chỗ: chọn chỗ sau → chỗ thứ 2, dòng dự phòng khác chỗ thứ thì bỏ.
func TestQueryForOccurrence(t *testing.T) {
	nodes := []Node{
		{Role: "window", Box: box(0, 0, 500, 500)},
		{Role: "button", Title: "Mua", Box: box(10, 10, 80, 30), Depth: 1},
		{Role: "button", Title: "Mua", Box: box(10, 60, 80, 30), Depth: 1},
	}
	got := QueryFor(nodes, 2)
	want := Pick{Selectors: []string{`text="Mua"`, `role=button[name="Mua"]`, "role=button"}, Occurrence: 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%+v, muốn %+v", got, want)
	}
}
