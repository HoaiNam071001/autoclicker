package store

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

func TestDraftRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if loadDraft(dir) != nil {
		t.Fatal("chưa có nháp mà loadDraft khác nil")
	}
	sc := model.NewScript()
	sc.Steps = []model.Step{modeltest.Node("a", "n1", model.StepClick, nil), modeltest.Node("b", "n2", model.StepText, nil)}
	sc.Start, sc.Edges = "a", []model.Edge{modeltest.Edge("a", model.PortSuccess, "b")}
	d := Draft{Script: sc, Path: "/tmp/k.json", Dirty: true, SavedAt: 123}
	d.Script.Version = model.ScriptVersion
	b, _ := encodeDraft(d)
	if err := os.WriteFile(filepath.Join(dir, draftFile), b, 0o644); err != nil {
		t.Fatal(err)
	}
	got := loadDraft(dir)
	if got == nil || !reflect.DeepEqual(*got, d) {
		t.Fatalf("nháp đọc lại khác: %+v", got)
	}
	// nháp hỏng → bỏ qua, không lỗi
	os.WriteFile(filepath.Join(dir, draftFile), []byte(`{"script":{"version":4,"steps":[{"type":"?"}]}}`), 0o644)
	if loadDraft(dir) != nil {
		t.Fatal("nháp hỏng phải bỏ qua")
	}
}
