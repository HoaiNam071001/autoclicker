package engine

import (
	"context"
	"image"
	"reflect"
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/element"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/vision"
	"ghosthand/internal/vision/visiontest"
)

func findAllText(text string, mod func(*model.Step)) model.Step {
	return modeltest.Node("a", "n1", model.StepFindAll, func(s *model.Step) {
		s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = model.FindByText, text, 2000, 1
		if mod != nil {
			mod(s)
		}
	})
}

// Chữ: chưa đủ minCount thì chờ tiếp; items theo thứ tự đọc, text = cả dòng, toạ độ màn hình (vùng 100, 50).
func TestFindAllTextWaitsForMin(t *testing.T) {
	ocr := &fakeOCR{reads: [][]vision.TextLine{
		{line("Thêm", 10, 10, 40, 20)},
		{line("Thêm 2", 10, 200, 60, 20), line("Thêm 1", 300, 10, 60, 20), line("Thêm 0", 10, 10, 60, 20)},
	}}
	e, _ := ctlEngine(model.Script{Start: "a", Steps: []model.Step{findAllText("thêm", func(s *model.Step) { s.MinCount = 3 })}}, ocr)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	out := e.results["a"].Output
	items, _ := out["items"].([]any)
	if out["count"] != 3 || len(items) != 3 {
		t.Fatalf("output %v", out)
	}
	var texts []string
	for _, it := range items {
		texts = append(texts, it.(map[string]any)["text"].(string))
	}
	if want := []string{"Thêm 0", "Thêm 1", "Thêm 2"}; !reflect.DeepEqual(texts, want) {
		t.Errorf("thứ tự %q, muốn %q", texts, want)
	}
	if items[0].(map[string]any)["y"] != 50+10+10 || items[0].(map[string]any)["score"] != 100 {
		t.Errorf("chỗ đầu %v", items[0])
	}
}

// Hết giờ: nhiều hơn maxCount → lỗi not_found, detail có count / minCount / maxCount.
func TestFindAllTooMany(t *testing.T) {
	ocr := &fakeOCR{reads: [][]vision.TextLine{{line("Gửi", 0, 0, 40, 20), line("Gửi", 0, 100, 40, 20)}}}
	e, _ := ctlEngine(model.Script{Start: "a", Steps: []model.Step{findAllText("gửi", func(s *model.Step) { s.MaxCount, s.TimeoutMs = 1, 5 })}}, ocr)
	err := e.Run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "nhiều hơn tối đa 1") {
		t.Fatalf("lỗi: %v", err)
	}
	r := e.results["a"]
	if r.Error == nil || r.Error.Code != model.ErrNotFound || r.Error.Detail["count"] != 2 || r.Error.Detail["maxCount"] != 1 || r.Error.Detail["minCount"] != 1 {
		t.Errorf("lỗi node %+v", r.Error)
	}
}

// Ảnh mẫu: hai nút giống nhau → hai chỗ, toạ độ màn hình theo vùng kịch bản.
func TestFindAllImage(t *testing.T) {
	img := visiontest.Scene(400, 300, 11)
	visiontest.Button(img, 40, 50, 60, 20)
	visiontest.Button(img, 250, 200, 60, 20)
	tpl, _ := vision.EncodePNG(visiontest.Crop(img, image.Rect(40, 50, 100, 70)))
	s := modeltest.Node("a", "n1", model.StepFindAll, func(s *model.Step) { s.Template, s.TemplateScale, s.TimeoutMs, s.IntervalMs = tpl, 1, 2000, 1 })
	e, _ := ctlEngine(model.Script{Start: "a", Steps: []model.Step{s}}, nil)
	e.scr = &fakeScreen{frames: []image.Image{img}}
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	items, _ := e.results["a"].Output["items"].([]any)
	if len(items) != 2 || items[0].(map[string]any)["x"] != 100+70 || items[0].(map[string]any)["y"] != 50+60 || items[1].(map[string]any)["x"] != 100+280 || items[1].(map[string]any)["y"] != 50+210 {
		t.Fatalf("items %v", items)
	}
	res, err := e.FindOnce(s)
	if err != nil || !res.Found || res.Count != 2 || len(res.All) != 2 || res.X != 140 {
		t.Errorf("Thử tìm: %+v %v", res, err)
	}
}

// Phần tử: chỉ lấy phần tử đang hiện (khung cao 1 px = ngoài khung nhìn).
func TestFindAllElement(t *testing.T) {
	nodes := append(page(elemBox(100, 200, 80, 30)), element.Node{Role: "button", Raw: "AXButton", Title: "Mua ngay", Box: elemBox(100, 790, 80, 1), Depth: 1},
		element.Node{Role: "button", Raw: "AXButton", Title: "Mua ngay", Box: elemBox(300, 200, 80, 30), Depth: 1})
	_, e, err := elementRun(t, &devicetest.FakeElements{Frames: [][]element.Node{nodes}}, elementStep(model.StepFindAll, nil))
	if err != nil {
		t.Fatal(err)
	}
	out := e.results[e.sc.Steps[0].ID].Output
	items, _ := out["items"].([]any)
	if out["count"] != 2 || items[1].(map[string]any)["x"] != 340 || items[0].(map[string]any)["role"] != "button" {
		t.Errorf("output %v", out)
	}
}

// Vòng findAll → forEach → click theo item.x / item.y: click từng chỗ một lần theo thứ tự đọc.
func TestFindAllForEachClicks(t *testing.T) {
	ocr := &fakeOCR{reads: [][]vision.TextLine{{line("Thêm", 300, 10, 40, 20), line("Thêm", 10, 10, 40, 20)}}}
	sc := model.Script{Start: "a", Steps: []model.Step{
		findAllText("thêm", nil),
		modeltest.Node("f", "n2", model.StepForEach, func(s *model.Step) { s.Value = "{{n1.output.items}}" }),
		modeltest.Node("c", "n3", model.StepClick, func(s *model.Step) {
			s.Points = []model.Point{{XVar: "{{n2.output.item.x}}", YVar: "{{n2.output.item.y}}"}}
		}),
	}, Edges: []model.Edge{modeltest.Edge("a", model.PortSuccess, "f"), modeltest.Edge("f", model.PortItem, "c"), modeltest.Edge("c", model.PortSuccess, "f")}}
	e, in := ctlEngine(sc, ocr)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"move 130,70", "click left false", "move 420,70", "click left false"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("có %q, muốn %q", in.Lines(), want)
	}
}
