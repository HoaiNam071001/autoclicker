package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func sampleScript() Script {
	return Script{
		Name:   "Thử lưu",
		Region: &Rect{X: 10, Y: 20, W: 1000, H: 700},
		Loops:  3,
		Steps: []Step{
			{ID: "a", Type: StepClick, Points: []Point{{X: 300, Y: 300, Ax: AnchorEnd}}, Button: "right", Clicks: 2, Keys: []string{}, KeyMode: KeyChord, Modifiers: []string{}, DelayAfter: 500, Enabled: true},
			{ID: "b", Type: StepDrag, Points: []Point{{X: 1, Y: 2}, {X: 5000, Y: 40, Ax: AnchorRatio, Ay: AnchorCenter}}, Button: "left", DragMs: 400, DragRelative: true, Keys: []string{}, KeyMode: KeyChord, Modifiers: []string{}},
			{ID: "c", Type: StepKey, Points: []Point{}, Button: "left", Keys: []string{"c", "t"}, KeyMode: KeySequence, Modifiers: []string{"ctrl"}, JitterMs: 50},
			{ID: "d", Type: StepText, Points: []Point{}, Button: "left", Text: "xin chào", Keys: []string{}, KeyMode: KeyChord, Modifiers: []string{}},
		},
		AnchorMode: AnchorTopLeft,
	}
}

func TestScriptFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "k.json")
	sc := sampleScript()
	if err := writeScriptFile(path, sc); err != nil {
		t.Fatal(err)
	}
	got, err := readScriptFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sc.Version = ScriptVersion
	if !reflect.DeepEqual(got, sc) {
		t.Fatalf("đọc lại khác:\n%+v\nmuốn\n%+v", got, sc)
	}
}

func TestDecodeScriptErrors(t *testing.T) {
	cases := map[string]string{
		"không phải JSON":   `{"version":1,`,
		"thiếu version":     `{"name":"x","steps":[]}`,
		"version mới hơn":   `{"version":99,"steps":[]}`,
		"loại bước lạ":      `{"version":1,"steps":[{"id":"a","type":"jump"}]}`,
		"thiếu điểm":        `{"version":1,"steps":[{"id":"a","type":"drag","points":[{"x":1,"y":1}]}]}`,
		"neo lạ":            `{"version":1,"steps":[{"id":"a","type":"click","points":[{"x":1,"y":1,"ax":"left"}]}]}`,
		"vùng rỗng":         `{"version":1,"region":{"x":0,"y":0,"w":0,"h":10},"steps":[]}`,
		"số âm":             `{"version":1,"steps":[{"id":"a","type":"wait","delayAfter":-1}]}`,
		"phím kèm lạ":       `{"version":1,"steps":[{"id":"a","type":"key","keys":["a"],"modifiers":["super"]}]}`,
		"cách neo kịch bản": `{"version":1,"anchorMode":"nearest","steps":[]}`,
	}
	for name, in := range cases {
		if _, err := decodeScript([]byte(in)); err == nil {
			t.Errorf("%s: muốn lỗi, không có", name)
		}
	}
}

// File sửa tay / tối giản: thiếu field thì bổ sung giá trị trống, id trống / trùng được đặt lại.
func TestDecodeScriptFillsDefaults(t *testing.T) {
	sc, err := decodeScript([]byte(`{"version":1,"steps":[{"type":"wait"},{"id":"x","type":"click","points":[{"x":1,"y":1}]},{"id":"x","type":"type"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for i, s := range sc.Steps {
		if s.ID == "" || ids[s.ID] {
			t.Fatalf("bước %d: id %q trống / trùng", i+1, s.ID)
		}
		ids[s.ID] = true
		if s.Points == nil || s.Keys == nil || s.Modifiers == nil || s.Button != "left" || s.KeyMode != KeyChord {
			t.Fatalf("bước %d chưa bổ sung mặc định: %+v", i+1, s)
		}
	}
	if sc.Steps[1].ID != "x" {
		t.Fatalf("id hợp lệ bị đổi: %q", sc.Steps[1].ID)
	}
}

func TestDraftRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if loadDraft(dir) != nil {
		t.Fatal("chưa có nháp mà loadDraft khác nil")
	}
	d := Draft{Script: sampleScript(), Path: "/tmp/k.json", Dirty: true, SavedAt: 123}
	d.Script.Version = ScriptVersion
	b, _ := encodeDraft(d)
	if err := os.WriteFile(filepath.Join(dir, draftFile), b, 0o644); err != nil {
		t.Fatal(err)
	}
	got := loadDraft(dir)
	if got == nil || !reflect.DeepEqual(*got, d) {
		t.Fatalf("nháp đọc lại khác: %+v", got)
	}
	// nháp hỏng → bỏ qua, không lỗi
	os.WriteFile(filepath.Join(dir, draftFile), []byte(`{"script":{"version":1,"steps":[{"type":"?"}]}}`), 0o644)
	if loadDraft(dir) != nil {
		t.Fatal("nháp hỏng phải bỏ qua")
	}
}

func TestRecentList(t *testing.T) {
	var l []string
	for _, p := range []string{"a", "b", "c", "d", "e", "f", "c"} {
		l = pushRecent(l, p)
	}
	if want := []string{"c", "f", "e", "d", "b"}; !reflect.DeepEqual(l, want) {
		t.Fatalf("pushRecent = %v, muốn %v", l, want)
	}
	if got := removeRecent(l, "e"); strings.Join(got, "") != "cfdb" {
		t.Fatalf("removeRecent = %v", got)
	}
}

func TestSafeFileName(t *testing.T) {
	if got := safeFileName(` a/b:c? `); got != "a-b-c-" {
		t.Fatalf("safeFileName = %q", got)
	}
	if safeFileName("  ") != "kich-ban" {
		t.Fatal("tên trống phải có tên mặc định")
	}
}
