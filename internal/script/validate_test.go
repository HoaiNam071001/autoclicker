package script_test

// Kiểm tra trước khi chạy (ValidateScript) – test ngoài package để dùng scripttest.Linear (scripttest import script).

import (
	"strings"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script"
	"ghosthand/internal/script/scripttest"
	"ghosthand/internal/vision"
	"ghosthand/internal/vision/visiontest"
)

func TestValidateScript(t *testing.T) {
	cases := map[string]model.Script{
		"không có node nào được bật": {Steps: []model.Step{modeltest.Step(model.StepWait, func(s *model.Step) { s.Enabled = false })}},
		"cần 2 điểm":                 {Steps: []model.Step{modeltest.Step(model.StepDrag, func(s *model.Step) { s.Points = s.Points[:1] })}},
		"chưa nhập nội dung gõ":      {Steps: []model.Step{modeltest.Step(model.StepText, nil)}},
		"chưa chọn phím":             {Steps: []model.Step{modeltest.Step(model.StepKey, func(s *model.Step) { s.Keys = nil })}},
		"không hợp lệ":               {Steps: []model.Step{{Type: "scroll", Enabled: true}}},
	}
	for want, sc := range cases {
		err := script.ValidateScript(scripttest.Linear(sc), model.Rect{})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi chứa %q, nhận %v", want, err)
		}
	}
}

func TestValidatePointOutsideRegion(t *testing.T) {
	r := model.Rect{W: 1000, H: 700}
	sc := model.Script{Loops: 1, Region: &r, Steps: []model.Step{
		modeltest.Step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 1200, Y: 10}} }),
	}}
	if err := script.ValidateScript(scripttest.Linear(sc), model.Rect{}); err == nil || !strings.Contains(err.Error(), "ngoài vùng") {
		t.Errorf("err = %v", err)
	}
	// B lệch so với A không bị kiểm tra theo vùng
	sc.Steps = []model.Step{modeltest.Step(model.StepDrag, func(s *model.Step) {
		s.Points = []model.Point{{X: 900, Y: 10}, {X: 500, Y: 0}}
		s.DragRelative = true
	})}
	if err := script.ValidateScript(scripttest.Linear(sc), model.Rect{}); err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestValidateScrollOpen(t *testing.T) {
	cases := map[string]model.Step{
		"hướng lạ":        modeltest.Step(model.StepScroll, func(s *model.Step) { s.Direction = "xeo" }),
		"0 nấc":           modeltest.Step(model.StepScroll, func(s *model.Step) { s.Amount = 0 }),
		"chưa nhập URL":   modeltest.Step(model.StepOpen, nil),
		"điểm ngoài vùng": modeltest.Step(model.StepScroll, func(s *model.Step) { s.Points = []model.Point{{X: 5000, Y: 5000}} }),
	}
	for name, s := range cases {
		sc := model.Script{Region: &model.Rect{W: 100, H: 100}, Steps: []model.Step{s}}
		if err := script.ValidateScript(scripttest.Linear(sc), model.Rect{}); err == nil || !strings.HasPrefix(err.Error(), "n1:") {
			t.Errorf("%s: muốn lỗi ở n1, có %v", name, err)
		}
	}
}

func TestValidateImageSteps(t *testing.T) {
	tpl, _ := vision.EncodePNG(visiontest.Scene(100, 50, 1))
	cases := map[string][]model.Step{
		"click theo bước trước nhưng không có chờ ảnh": {modeltest.Step(model.StepClickImage, nil)},
		"chờ ảnh chưa có ảnh mẫu":                      {modeltest.Step(model.StepWaitImage, nil)},
		"ảnh mẫu hỏng":                                 {modeltest.Step(model.StepWaitImage, func(s *model.Step) { s.Template = "abc" })},
		"độ khớp 0":                                    {modeltest.Step(model.StepWaitImage, func(s *model.Step) { s.Template, s.Threshold = tpl, 0 })},
		"màu sai":                                      {modeltest.Step(model.StepWaitPixel, func(s *model.Step) { s.Color = "xanh" })},
	}
	for name, steps := range cases {
		if err := script.ValidateScript(scripttest.Linear(model.Script{Steps: steps}), model.Rect{}); err == nil {
			t.Errorf("%s: muốn lỗi", name)
		}
	}
	ok := []model.Step{modeltest.Step(model.StepWaitImage, func(s *model.Step) { s.Template = tpl }), modeltest.Step(model.StepClickImage, nil),
		modeltest.Step(model.StepDrag, func(s *model.Step) { s.DragFromPrev = true })}
	if err := script.ValidateScript(scripttest.Linear(model.Script{Steps: ok}), model.Rect{}); err != nil {
		t.Errorf("chờ ảnh → click theo ảnh → kéo từ điểm trước: %v", err)
	}
}

func TestValidateTextStep(t *testing.T) {
	sc := scripttest.TextScript("  ")
	if err := script.ValidateScript(scripttest.Linear(sc), model.Rect{W: 1920, H: 1080}); err == nil || !strings.Contains(err.Error(), "chưa nhập chữ") {
		t.Errorf("chữ trống: %v", err)
	}
	sc = scripttest.TextScript("Gửi cho {{ten}}")
	sc.Inputs = []model.ScriptInput{{Key: "ten", Label: "Tên", Type: model.InputText}}
	out, err := script.ResolveInputs(sc, map[string]any{"ten": "Lan"})
	if err != nil || out.Steps[0].FindText != "Gửi cho Lan" {
		t.Errorf("biến trong chữ cần tìm: %v %q", err, out.Steps[0].FindText)
	}
}

func TestValidateElementStep(t *testing.T) {
	el := func(sels ...string) model.Script {
		return model.Script{Loops: 1, Steps: []model.Step{modeltest.Step(model.StepWaitImage, func(s *model.Step) {
			s.FindBy, s.Selectors, s.TimeoutMs, s.IntervalMs = model.FindByElement, sels, 1000, 100
		})}}
	}
	cases := map[string]model.Script{
		"chưa nhập cách tìm phần tử": el(" "),
		"chỉ dùng được ở node web":   el("css=.buy"),
		"sai dạng":                   el("role=button[name"),
		"điều kiện phần tử": {Loops: 1, Steps: []model.Step{modeltest.Step(model.StepIf, func(s *model.Step) {
			s.Cond = &model.Cond{Kind: model.CondElement}
		})}},
	}
	for want, sc := range cases {
		err := script.ValidateScript(scripttest.Linear(sc), model.Rect{})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi chứa %q, nhận %v", want, err)
		}
	}
	// dòng có biến chỉ kiểm lúc chạy
	if err := script.ValidateScript(scripttest.Linear(el("{{loai}}=Mua", "text=Mua")), model.Rect{}); err != nil {
		t.Errorf("err = %v", err)
	}
}

// Đọc phần tử (37 E4): tìm nhãn bằng OCR phải đọc cạnh nhãn, mẫu lọc phải đọc được.
func TestValidateElementRead(t *testing.T) {
	rd := func(mod func(s *model.Step)) model.Script {
		return model.Script{Loops: 1, Steps: []model.Step{modeltest.Step(model.StepElementRead, func(s *model.Step) {
			s.Selectors, s.TimeoutMs, s.IntervalMs = []string{"text=Giá"}, 1000, 100
			mod(s)
		})}}
	}
	cases := map[string]model.Script{
		"chưa nhập chữ của nhãn": rd(func(s *model.Step) { s.FindBy, s.ReadNear = model.FindByText, model.ReadNearRight }),
		"phải chọn đọc bên phải": rd(func(s *model.Step) { s.FindBy, s.FindText = model.FindByText, "Giá" }),
		"mẫu lọc sai":            rd(func(s *model.Step) { s.ReadPattern = "(\\d+" }),
		"khoảng đọc phải trong":  rd(func(s *model.Step) { s.ReadNear, s.ReadSpan = model.ReadNearBelow, -1 }),
	}
	for want, sc := range cases {
		err := script.ValidateScript(scripttest.Linear(sc), model.Rect{})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi chứa %q, nhận %v", want, err)
		}
	}
	ok := rd(func(s *model.Step) {
		s.FindBy, s.FindText, s.ReadNear, s.ReadPattern = model.FindByText, "Giá", model.ReadNearRight, `(\d+)`
	})
	if err := script.ValidateScript(scripttest.Linear(ok), model.Rect{}); err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestValidateCoordSpaceAbsolute(t *testing.T) {
	r := model.Rect{X: 500, Y: 300, W: 400, H: 300}
	screen := model.Rect{W: 1920, H: 1080}
	mk := func(space string) model.Script {
		return scripttest.Linear(model.Script{Loops: 1, Region: &r, Steps: []model.Step{
			modeltest.Step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 100, Y: 100}}; s.CoordSpace = space }),
		}})
	}
	// region: 100,100 cộng gốc vùng = 600,400 → trong vùng; absolute: 100,100 là chỗ thật, ngoài vùng nhưng trong màn hình → hợp lệ
	for _, sp := range []string{"", "region", "absolute"} {
		if err := script.ValidateScript(mk(sp), screen); err != nil {
			t.Errorf("coordSpace %q: %v", sp, err)
		}
	}
	// ngoài màn hình: lỗi nêu gốc đang cộng
	sc := mk("")
	sc.Steps[0].Points = []model.Point{{X: 900, Y: 10}}
	if err := script.ValidateScript(sc, screen); err == nil || !strings.Contains(err.Error(), "gốc vùng 500,300") {
		t.Errorf("err = %v", err)
	}
	if p := script.ResolvePoint(model.Point{X: 100, Y: 100}, script.StepBounds(sc, screen, model.Step{CoordSpace: "absolute"})); p != (model.Point{X: 100, Y: 100}) {
		t.Errorf("absolute = %v", p)
	}
}

func TestValidateWaitScreen(t *testing.T) {
	ws := func(f func(*model.Step)) model.Script {
		return model.Script{Loops: 1, Steps: []model.Step{modeltest.Step(model.StepWaitScreen, f)}}
	}
	if err := script.ValidateScript(scripttest.Linear(ws(nil)), model.Rect{}); err != nil {
		t.Errorf("mặc định phải hợp lệ: %v", err)
	}
	cases := map[string]func(*model.Step){
		"chưa chọn vùng chụp":           func(s *model.Step) { s.ShotArea = model.ShotRect },
		"phải > 0":                      func(s *model.Step) { s.IntervalMs = 0 },
		"chưa chọn node lấy khung chụp": func(s *model.Step) { s.ShotArea = model.ShotNode },
	}
	for want, f := range cases {
		if err := script.ValidateScript(scripttest.Linear(ws(f)), model.Rect{}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi chứa %q, nhận %v", want, err)
		}
	}
}
