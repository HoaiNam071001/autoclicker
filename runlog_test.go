package main

import (
	"context"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"autoclicker/internal/model"
)

func TestRunRecordAddKeepsLastAndEndMessage(t *testing.T) {
	var r RunRecord
	for i := 0; i < maxRunResults+3; i++ {
		r.add(model.StepResult{NodeID: "a", Attempts: i}, model.StepText)
	}
	if len(r.Results) != maxRunResults || r.Truncated != 3 || r.Results[0].Attempts != 3 || r.Nodes != maxRunResults+3 {
		t.Errorf("giữ %d kết quả, bỏ %d, đầu %d", len(r.Results), r.Truncated, r.Results[0].Attempts)
	}
	r.add(model.StepResult{Status: model.StatusSuccess, Output: map[string]any{"message": "Đã gửi"}}, model.StepEnd)
	r.add(model.StepResult{Status: model.StatusSuccess, Output: map[string]any{"message": "không phải end"}}, model.StepText)
	if r.Message != "Đã gửi" {
		t.Errorf("lời nhắn = %q", r.Message)
	}
}

func TestSaveRunPrunesAndLists(t *testing.T) {
	dir := runsDir(t.TempDir(), "Nhắn tin / Facebook", "/x/a.json")
	if !strings.HasPrefix(filepath.Base(dir), "nhắn-tin---facebook-") {
		t.Errorf("thư mục %q", dir)
	}
	if runsDir("/b", "A", "/x/a.json") == runsDir("/b", "A", "/x/b.json") {
		t.Error("2 file khác nhau cùng tên phải tách lịch sử")
	}
	start := time.Date(2026, 10, 17, 14, 0, 0, 0, time.Local)
	for i := 0; i < maxRuns+2; i++ {
		at := start.Add(time.Duration(i) * time.Second)
		rec := RunRecord{RunSummary: RunSummary{ID: newRunID(at), StartedAt: at, Status: RunSuccess}}
		shots := map[string]image.Image{}
		if i < 2 {
			shots["1.png"] = image.NewRGBA(image.Rect(0, 0, 4, 4))
		}
		if err := saveRun(dir, rec, shots); err != nil {
			t.Fatal(err)
		}
	}
	runs := listRuns(dir)
	if len(runs) != maxRuns || runs[0].ID != newRunID(start.Add((maxRuns+1)*time.Second)) {
		t.Fatalf("còn %d lần, mới nhất %v", len(runs), runs[0].ID)
	}
	if pngs, _ := filepath.Glob(filepath.Join(dir, "*.png")); len(pngs) != 0 {
		t.Errorf("ảnh của lần chạy đã xoá vẫn còn: %v", pngs)
	}
	if _, err := loadRun(dir, "../settings"); err == nil {
		t.Error("id có đường dẫn phải bị chặn")
	}
	if r, err := loadRun(dir, runs[0].ID); err != nil || r.Status != RunSuccess {
		t.Errorf("đọc lại: %+v %v", r, err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatal(err)
	}
}

func TestShootOnFail(t *testing.T) {
	e, _ := sendFlow(-1, nil)
	e.sc.Edges = e.sc.Edges[:1] // không thấy → dừng
	e.shots = map[string]image.Image{}
	var got []model.StepResult
	e.onNode = func(r model.StepResult) { got = append(got, r) }
	if err := e.run(context.Background()); err == nil {
		t.Fatal("phải lỗi")
	}
	if len(got) != 1 || got[0].Shot != "1.png" || e.shots["1.png"] == nil || got[0].Loop != 1 {
		t.Errorf("kết quả %+v, ảnh %v", got, e.shots)
	}
	if img := shrink(image.NewRGBA(image.Rect(0, 0, 3200, 2000)), shotMaxW); img.Bounds().Dx() != 1600 || img.Bounds().Dy() != 1000 {
		t.Errorf("thu nhỏ: %v", img.Bounds())
	}
}

// Debug: breakpoint chỉ dừng khi debug bật; chạy từng node dừng trước mọi node; Tiếp tắt chạy từng node.
func TestDebugBreakpointAndStep(t *testing.T) {
	sc := model.Script{Loops: 1, Start: "a", Steps: []model.Step{
		node("a", "n1", model.StepText, func(s *model.Step) { s.Text = "a" }),
		node("b", "n2", model.StepText, func(s *model.Step) { s.Text = "b"; s.Breakpoint = true }),
		node("c", "n3", model.StepText, func(s *model.Step) { s.Text = "c" }),
	}, Edges: []model.Edge{edge("a", model.PortSuccess, "b"), edge("b", model.PortSuccess, "c")}}

	// không debug: chạy thẳng
	if err := newEngine(&fakeInput{}, sc, nil).run(context.Background()); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var breaks []string
	in := &fakeInput{}
	e := newEngine(in, sc, func(p Progress) {
		if p.Break != "" {
			mu.Lock()
			breaks = append(breaks, p.Ref+":"+p.Break)
			mu.Unlock()
		}
	})
	e.debug = true
	done := make(chan error)
	go func() { done <- e.run(context.Background()) }()
	waitFor := func(n int) {
		t.Helper()
		for i := 0; i < 100; i++ {
			mu.Lock()
			k := len(breaks)
			mu.Unlock()
			if k >= n && e.gate.paused() {
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
		t.Fatalf("không dừng lần %d (%v)", n, breaks)
	}
	waitFor(1) // dừng trước n2 – n1 đã chạy
	if l := in.lines(); len(l) != 1 || l[0] != "type a" {
		t.Errorf("trước breakpoint: %q", l)
	}
	e.Step() // chạy n2 rồi dừng trước n3
	waitFor(2)
	if l := in.lines(); len(l) != 2 {
		t.Errorf("sau F10: %q", l)
	}
	e.Resume()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if w := "n2:breakpoint n3:step"; strings.Join(breaks, " ") != w {
		t.Errorf("dừng ở %v, muốn %s", breaks, w)
	}
}
