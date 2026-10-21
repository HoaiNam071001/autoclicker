package script

import (
	"reflect"
	"testing"

	"ghosthand/internal/model"
)

func sampleScript() model.Script {
	return model.Script{
		Name:   "Thử lưu",
		Region: &model.Rect{X: 10, Y: 20, W: 1000, H: 700},
		Loops:  3,
		Steps: []model.Step{
			{ID: "a", Ref: "n1", Type: model.StepClick, Points: []model.Point{{X: 300, Y: 300, Ax: AnchorEnd}}, Button: "right", Clicks: 2, Keys: []string{}, KeyMode: model.KeyChord, Modifiers: []string{}, DelayAfter: 500, Enabled: true},
			{ID: "b", Ref: "n2", Type: model.StepDrag, Points: []model.Point{{X: 1, Y: 2}, {X: 5000, Y: 40, Ax: AnchorRatio, Ay: AnchorCenter}}, Button: "left", DragMs: 400, DragRelative: true, Keys: []string{}, KeyMode: model.KeyChord, Modifiers: []string{}},
			{ID: "c", Ref: "n3", Type: model.StepKey, Points: []model.Point{}, Button: "left", Keys: []string{"c", "t"}, KeyMode: model.KeySequence, Modifiers: []string{"ctrl"}, JitterMs: 50},
			{ID: "d", Ref: "n4", Type: model.StepText, Points: []model.Point{}, Button: "left", Text: "xin chào", Keys: []string{}, KeyMode: model.KeyChord, Modifiers: []string{}},
		},
		Start: "a",
		Edges: []model.Edge{
			{ID: "e1", From: "a", Port: model.PortSuccess, To: "b"},
			{ID: "e2", From: "a", Port: model.PortFail, To: "d"},
			{ID: "e3", From: "b", Port: model.PortSuccess, To: "c"},
		},
		AnchorMode: model.AnchorTopLeft,
	}
}

func TestScriptRoundTrip(t *testing.T) {
	sc := sampleScript()
	b, err := EncodeScript(sc)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeScript(b)
	if err != nil {
		t.Fatal(err)
	}
	sc.Version = model.ScriptVersion
	if !reflect.DeepEqual(got, sc) {
		t.Fatalf("đọc lại khác:\n%+v\nmuốn\n%+v", got, sc)
	}
}

func TestDecodeScriptErrors(t *testing.T) {
	cases := map[string]string{
		"không phải JSON":   `{"version":4,`,
		"thiếu version":     `{"name":"x","steps":[]}`,
		"version mới hơn":   `{"version":99,"steps":[]}`,
		"loại bước lạ":      `{"version":4,"steps":[{"id":"a","type":"jump"}]}`,
		"thiếu điểm":        `{"version":4,"steps":[{"id":"a","type":"drag","points":[{"x":1,"y":1}]}]}`,
		"neo lạ":            `{"version":4,"steps":[{"id":"a","type":"click","points":[{"x":1,"y":1,"ax":"left"}]}]}`,
		"vùng rỗng":         `{"version":4,"region":{"x":0,"y":0,"w":0,"h":10},"steps":[]}`,
		"số âm":             `{"version":4,"steps":[{"id":"a","type":"wait","delayAfter":-1}]}`,
		"phím kèm lạ":       `{"version":4,"steps":[{"id":"a","type":"key","keys":["a"],"modifiers":["super"]}]}`,
		"cách neo kịch bản": `{"version":4,"anchorMode":"nearest","steps":[]}`,
	}
	for name, in := range cases {
		if _, err := DecodeScript([]byte(in)); err == nil {
			t.Errorf("%s: muốn lỗi, không có", name)
		}
	}
}

// Sửa tay file: cạnh trỏ node không còn / thiếu cổng / trùng cổng bị bỏ, id cạnh trống được đặt, start không còn → trống.
func TestDecodeScriptRepairsEdges(t *testing.T) {
	sc, err := DecodeScript([]byte(`{"version":4,"start":"zz","steps":[{"id":"a","type":"wait"},{"id":"b","type":"wait"}],
		"edges":[{"from":"a","port":"success","to":"b"},{"from":"a","port":"success","to":"a"},{"from":"a","port":"fail","to":"x"},{"from":"b","to":"a"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(sc.Edges) != 1 || sc.Edges[0].To != "b" || sc.Edges[0].ID == "" || sc.Start != "" {
		t.Fatalf("cạnh / start chưa sửa: %+v start %q", sc.Edges, sc.Start)
	}
	if sc.Steps[0].Ref != "n1" || sc.Steps[1].Ref != "n2" {
		t.Fatalf("chưa đặt ref: %q %q", sc.Steps[0].Ref, sc.Steps[1].Ref)
	}
}

// File sửa tay / tối giản: thiếu field thì bổ sung giá trị trống, id trống / trùng được đặt lại.
func TestDecodeScriptFillsDefaults(t *testing.T) {
	sc, err := DecodeScript([]byte(`{"version":4,"steps":[{"type":"wait"},{"id":"x","type":"click","points":[{"x":1,"y":1}]},{"id":"x","type":"type"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for i, s := range sc.Steps {
		if s.ID == "" || ids[s.ID] {
			t.Fatalf("bước %d: id %q trống / trùng", i+1, s.ID)
		}
		ids[s.ID] = true
		if s.Points == nil || s.Keys == nil || s.Modifiers == nil || s.Button != "left" || s.KeyMode != model.KeyChord {
			t.Fatalf("bước %d chưa bổ sung mặc định: %+v", i+1, s)
		}
	}
	if sc.Steps[1].ID != "x" {
		t.Fatalf("id hợp lệ bị đổi: %q", sc.Steps[1].ID)
	}
}

func TestSafeFileName(t *testing.T) {
	if got := SafeFileName(` a/b:c? `); got != "a-b-c-" {
		t.Fatalf("safeFileName = %q", got)
	}
	if SafeFileName("  ") != "kich-ban" {
		t.Fatal("tên trống phải có tên mặc định")
	}
}
