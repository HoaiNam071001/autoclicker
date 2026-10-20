package mcpsrv

import (
	"image"
	"maps"
	"slices"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/session"
)

// Ảnh gửi kèm: mặc định ảnh lỗi + ảnh sau node mới nhất; AI xin theo tên thì đúng ảnh đó; ảnh đã bỏ không báo tên.
func TestShotImages(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 40, 20))
	tr := session.Trace{Shots: map[string]image.Image{}}
	for _, r := range []model.StepResult{
		{NodeID: "a", Ref: "n1", Status: model.StatusSuccess, Shot: "each-1-n1.png"}, // ảnh đã bị bỏ (giữ 12 ảnh mới nhất)
		{NodeID: "b", Ref: "n2", Status: model.StatusSuccess, Shot: "each-2-n2.png"},
		{NodeID: "c", Ref: "n3", Status: model.StatusSuccess, Shot: "each-3-n3.png"},
		{NodeID: "d", Ref: "n4", Status: model.StatusSuccess, Shot: "each-4-n4.png"},
		{NodeID: "e", Ref: "n5", Status: model.StatusSuccess, Shot: "each-5-n5.png"},
		{NodeID: "f", Ref: "n6", Status: model.StatusFail, Shot: "1.png"},
	} {
		tr.Results = append(tr.Results, r)
		if r.Ref != "n1" {
			tr.Shots[r.Shot] = img
		}
	}
	sent := func(w map[string]int) []string { return slices.Sorted(maps.Keys(w)) }

	_, w := shotImages(tr, nil)
	if got := sent(w); !slices.Equal(got, []string{"1.png", "each-3-n3.png", "each-4-n4.png", "each-5-n5.png"}) {
		t.Errorf("mặc định gửi %v", got)
	}
	_, w = shotImages(tr, []string{"each-2-n2.png", "each-1-n1.png", "khong-co.png"})
	if got := sent(w); !slices.Equal(got, []string{"each-2-n2.png"}) {
		t.Errorf("xin theo tên gửi %v", got)
	}
	v := traceView(tr, model.Script{}, nil)["nodes"].([]nodeView)
	if v[0].Shot != "" || v[1].Shot != "each-2-n2.png" {
		t.Errorf("tên ảnh trong vết: %q %q", v[0].Shot, v[1].Shot)
	}
	if _, err := shotEach("all"); err == nil {
		t.Error("shots sai phải lỗi")
	}
}
