package session

import (
	"context"
	"testing"

	"ghosthand/internal/model"
)

// Kịch bản 2 node Chờ s1 → s2.
func twoWaits(t *testing.T) model.Script {
	t.Helper()
	sc := waitScript(t, 1)
	s2 := sc.Steps[0]
	s2.ID, s2.Ref = "s2", "n2"
	sc.Steps = append(sc.Steps, s2)
	sc.Edges = []model.Edge{{ID: "e1", From: "s1", Port: model.PortSuccess, To: "s2"}}
	return sc
}

func TestStopAtTraceSeed(t *testing.T) {
	m, _, _ := testManager(t)
	id, err := m.Start(twoWaits(t), nil, Options{Source: SourceAI, StopAt: "s1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.Wait(context.Background(), id); err != nil {
		t.Fatal(err)
	}
	tr, err := m.Trace(id)
	if err != nil {
		t.Fatal(err)
	}
	if !tr.Reached || tr.State != StateSuccess || len(tr.Results) != 1 || tr.Results[0].NodeID != "s1" {
		t.Fatalf("trace = %+v", tr)
	}

	id2, err := m.Start(twoWaits(t), nil, Options{Source: SourceAI, From: "s2", SeedRun: id})
	if err != nil {
		t.Fatal(err)
	}
	m.Wait(context.Background(), id2)
	tr, _ = m.Trace(id2)
	if len(tr.Results) != 1 || tr.Results[0].NodeID != "s2" || tr.Reached {
		t.Errorf("chạy tiếp = %+v", tr)
	}

	if _, err := m.Start(twoWaits(t), nil, Options{StopAt: "n9"}); err == nil {
		t.Error("StopAt node không có mà không báo lỗi")
	}
	if _, err := m.Start(twoWaits(t), nil, Options{SeedRun: "khong-co"}); err == nil {
		t.Error("SeedRun không có mà không báo lỗi")
	}
}
