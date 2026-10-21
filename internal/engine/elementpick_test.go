package engine

import (
	"reflect"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/element"
	"ghosthand/internal/model"
)

func TestPickElement(t *testing.T) {
	e := newEngine(&devicetest.FakeInput{}, model.Script{}, nil)
	e.elems = &devicetest.FakeElements{Frames: [][]element.Node{page(elemBox(100, 200, 80, 30))}}
	got, err := e.PickElement("", model.Point{X: 120, Y: 210})
	if err != nil {
		t.Fatal(err)
	}
	want := ElementPick{Selectors: []string{`text="Mua ngay"`, `role=button[name="Mua ngay"]`, "role=button"}, Occurrence: 1,
		Role: "button", Name: "Mua ngay", Box: elemBox(100, 200, 80, 30)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%+v, muốn %+v", got, want)
	}
	if _, err := e.PickElement("", model.Point{X: 2000, Y: 10}); err == nil {
		t.Error("ngoài cửa sổ phải báo lỗi")
	}
}
