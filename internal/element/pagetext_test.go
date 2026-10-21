package element

import (
	"reflect"
	"testing"

	"ghosthand/internal/model"
)

func TestTextLines(t *testing.T) {
	box := model.Rect{X: 0, Y: 0, W: 100, H: 20}
	nodes := []Node{
		{Role: "window", Title: "Trang A", Box: box, Depth: 0},
		{Role: "document", Box: box, Depth: 1},
		{Role: "heading", Title: "Giá vé", Box: box, Depth: 2},
		{Role: "button", Title: "Gửi", Box: box, Depth: 2}, // bọc khối chữ → chỉ lấy khối chữ
		{Role: RoleText, Value: "Gửi", Box: box, Depth: 3},
		{Role: RoleText, Value: "dòng 1\n  dòng   2", Box: box, Depth: 2},
		{Role: RoleText, Value: "dòng 2", Box: box, Depth: 2}, // trùng dòng liền trước
		{Role: RoleText, Value: "ngoài khung", Box: model.Rect{W: 100, H: 0}, Depth: 2},
		{Role: "textbox", Value: "abc", Depth: 2, Box: box},
		{Role: "button", Desc: "icon", Box: box, Depth: 2}, // chỉ có nhãn mô tả → không phải nội dung
	}
	want := []string{"Giá vé", "Gửi", "dòng 1", "dòng 2", "abc"}
	if got := TextLines(nodes, 0, len(nodes), false); !reflect.DeepEqual(got, want) {
		t.Errorf("cả cửa sổ: %q, muốn %q", got, want)
	}
	if got := TextLines(nodes, 0, len(nodes), true); len(got) != 6 || got[4] != "ngoài khung" {
		t.Errorf("gồm ngoài khung: %q", got)
	}
	if end := SubtreeEnd(nodes, 3); end != 5 {
		t.Errorf("SubtreeEnd = %d, muốn 5", end)
	}
	if got := TextLines(nodes, 3, SubtreeEnd(nodes, 3), false); !reflect.DeepEqual(got, []string{"Gửi"}) {
		t.Errorf("cây con nút: %q", got)
	}
}
