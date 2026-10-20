package script

import (
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

func TestRunFrom(t *testing.T) {
	sc := model.Script{Start: "a", Loops: 5, Edges: []model.Edge{
		modeltest.Edge("a", model.PortSuccess, "b"), modeltest.Edge("b", model.PortSuccess, "c"), modeltest.Edge("b", model.PortFail, "a"),
	}}
	got := RunFrom(sc, "b", false)
	if got.Start != "b" || got.Loops != 1 || len(got.Edges) != 3 {
		t.Errorf("chạy từ b = %+v", got)
	}
	got = RunFrom(sc, "b", true)
	if got.Start != "b" || got.Loops != 1 || len(got.Edges) != 1 || got.Edges[0].From != "a" {
		t.Errorf("chạy riêng b = %+v", got)
	}
	if len(sc.Edges) != 3 {
		t.Error("RunFrom sửa mất cạnh của kịch bản gốc")
	}
}

func TestRunPart(t *testing.T) {
	sc := model.Script{Start: "a", Loops: 5, Steps: []model.Step{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}}, Edges: []model.Edge{
		modeltest.Edge("a", model.PortSuccess, "b"), modeltest.Edge("b", model.PortSuccess, "c"), modeltest.Edge("c", model.PortSuccess, "d"),
		modeltest.Edge("c", model.PortFail, "a"),
	}}
	got := RunPart(sc, "", []string{"c", "b"})
	if got.Start != "b" || got.Loops != 1 || len(got.Edges) != 1 || got.Edges[0].From != "b" || got.Edges[0].To != "c" {
		t.Errorf("cụm b, c = %+v", got)
	}
	got = RunPart(sc, "c", []string{"d"}) // from luôn thuộc cụm
	if got.Start != "c" || len(got.Edges) != 1 || got.Edges[0].To != "d" {
		t.Errorf("cụm d từ c = %+v", got)
	}
}
