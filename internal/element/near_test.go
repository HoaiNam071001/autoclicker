package element

import (
	"testing"

	"ghosthand/internal/model"
)

func TestNearMatches(t *testing.T) {
	q := Query{Selectors: []model.Selector{sel(t, "text=Them vao gio hang"), sel(t, `role=button[name="Dong lai"]`)}}
	got := NearMatches(tree(), q, 3)
	if len(got) != 3 || got[0].Role != "button" || got[0].Name != "Thêm vào giỏ" || got[0].Score < 60 {
		t.Fatalf("ứng viên = %+v", got)
	}
	for _, n := range got {
		if n.Name == "Shop" || n.Name == "Tìm sản phẩm" {
			t.Errorf("không liên quan mà vẫn có: %+v", n)
		}
	}
	if got := NearMatches(tree(), Query{Selectors: []model.Selector{sel(t, "role=checkbox")}}, 5); len(got) != 0 {
		t.Errorf("sai vai trò, không tên → không ứng viên: %+v", got)
	}
}
