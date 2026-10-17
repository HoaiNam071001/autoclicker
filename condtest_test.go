package main

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestTestCondParts(t *testing.T) {
	_, tpl := screenWithButton(1, true)
	screen, _ := screenWithButton(1, true)
	c := Cond{Kind: CondAll, Conds: []Cond{
		{Kind: CondText, FindText: "Gửi"},
		{Kind: CondImage, Template: tpl, TemplateScale: 1, Threshold: 90},
		{Kind: CondPixel, Points: []Point{{X: 10, Y: 10}}, Color: "#FF0000", Tolerance: 5, Not: true},
		{Kind: CondInput, Left: "Lan", Op: "eq", Right: "lan"},
	}}
	e := newEngine(nil, Script{Region: &Rect{W: 400, H: 300}}, nil)
	e.scr = &fakeScreen{frames: []image.Image{screen}, colors: []color.RGBA{{0, 0, 255, 255}}}
	e.ocr = &fakeOCR{reads: [][]TextLine{sendBtn}}
	var parts []CondPart
	ok, known, err := e.testCond(c, "s", nil, &parts)
	if err != nil || !ok || !known || len(parts) != 4 {
		t.Fatalf("ok %v known %v err %v parts %d", ok, known, err, len(parts))
	}
	for i, want := range []string{"thấy “Gửi” tại X 120 Y 110", "thấy tại X 280 Y 130 · khớp", "màu #0000FF (cần #FF0000 ± 5", "“Lan” = “lan”"} {
		if !strings.Contains(parts[i].Detail, want) || !parts[i].Ok || len(parts[i].Path) != 1 || parts[i].Path[0] != i {
			t.Errorf("phần %d: %+v, muốn chứa %q", i, parts[i], want)
		}
	}
	if parts[1].Hit == nil || parts[0].Hit == nil {
		t.Error("ảnh / chữ phải có chỗ thấy để vẽ khung")
	}

	// bước trước: chỉ biết lúc chạy – "và" với phần sai thì vẫn biết là sai
	parts = nil
	ok, known, _ = e.testCond(Cond{Kind: CondAll, Conds: []Cond{{Kind: CondResult}, {Kind: CondInput, Left: "a", Op: "eq", Right: "b"}}}, "s", nil, &parts)
	if ok || !known || parts[0].Known {
		t.Errorf("result + sai: ok %v known %v", ok, known)
	}
	parts = nil
	_, known, _ = e.testCond(Cond{Kind: CondAll, Conds: []Cond{{Kind: CondResult}, {Kind: CondInput, Left: "a", Op: "eq", Right: "a"}}}, "s", nil, &parts)
	if known {
		t.Error("result + đúng: chưa biết được")
	}
}
