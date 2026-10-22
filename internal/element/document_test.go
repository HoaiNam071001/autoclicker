package element

import (
	"testing"

	"ghosthand/internal/model"
)

func TestDocumentIndex(t *testing.T) {
	nodes := []Node{
		{Role: "window"},
		{Role: "searchbox", Value: "tìm", Box: model.Rect{W: 100, H: 20}},
		{Role: "textbox", Value: "thân bài\ndòng hai", Box: model.Rect{W: 800, H: 600}, Selected: "dòng hai"},
		{Role: "textbox", Value: "", Box: model.Rect{W: 900, H: 900}},
	}
	if got := DocumentIndex(nodes, false); got != 2 {
		t.Errorf("lớn nhất có chữ: %d", got)
	}
	if got := DocumentIndex(nodes, true); got != 2 {
		t.Errorf("phần chọn: %d", got)
	}
	nodes[2].Selected = ""
	if got := DocumentIndex(nodes, true); got != -1 {
		t.Errorf("không chọn gì: %d", got)
	}
	if got := DocumentIndex(nil, false); got != -1 {
		t.Errorf("rỗng: %d", got)
	}
}
