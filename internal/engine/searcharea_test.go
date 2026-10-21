package engine

import (
	"testing"

	"ghosthand/internal/device"
	"ghosthand/internal/model"
)

func TestCaptureAreaDisplays(t *testing.T) {
	ds := []model.Rect{{X: 0, Y: 0, W: 1920, H: 1080}, {X: 1920, Y: 406, W: 1512, H: 982}}
	if err := device.CaptureArea(model.Rect{X: 1920, Y: 406, W: 1512, H: 982}, ds); err != nil {
		t.Errorf("cả màn phụ: %v", err)
	}
	if err := device.CaptureArea(model.Rect{X: 1800, Y: 500, W: 300, H: 200}, ds); err == nil {
		t.Error("vắt qua 2 màn phải báo lỗi")
	}
	if err := device.CaptureArea(model.Rect{X: 100, Y: 2000, W: 50, H: 50}, ds); err == nil {
		t.Error("ngoài màn phải báo lỗi")
	}
	e := &Engine{sc: model.Script{Region: &model.Rect{X: 2000, Y: 500, W: 300, H: 200}}, screen: ds[0], displays: ds}
	if got := e.searchArea(model.Step{SearchArea: "screen"}); got != ds[1] {
		t.Errorf("cả màn với vùng ở màn phụ: %v", got)
	}
}
