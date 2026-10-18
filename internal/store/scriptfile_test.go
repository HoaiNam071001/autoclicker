package store

import (
	"path/filepath"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

func TestScriptFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Nhắn tin.json")
	sc := model.NewScript()
	sc.Name = ""
	sc.Steps = []model.Step{modeltest.Node("a", "n1", model.StepClick, nil)}
	sc.Start = "a"
	if err := WriteScriptFile(path, sc); err != nil {
		t.Fatal(err)
	}
	got, err := ReadScriptFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Nhắn tin" || len(got.Steps) != 1 || got.Start != "a" {
		t.Fatalf("đọc lại khác: %+v", got)
	}
}
