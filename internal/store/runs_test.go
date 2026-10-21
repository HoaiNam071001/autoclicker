package store

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/model"
)

func TestSaveRunPrunesAndLists(t *testing.T) {
	dir := RunsDir(t.TempDir(), "Nhắn tin / Facebook", "/x/a.json")
	if !strings.HasPrefix(filepath.Base(dir), "nhắn-tin---facebook-") {
		t.Errorf("thư mục %q", dir)
	}
	if RunsDir("/b", "A", "/x/a.json") == RunsDir("/b", "A", "/x/b.json") {
		t.Error("2 file khác nhau cùng tên phải tách lịch sử")
	}
	start := time.Date(2026, 10, 17, 14, 0, 0, 0, time.Local)
	for i := 0; i < maxRuns+2; i++ {
		at := start.Add(time.Duration(i) * time.Second)
		rec := model.RunRecord{RunSummary: model.RunSummary{ID: NewRunID(at), StartedAt: at, Status: model.RunSuccess}}
		shots := map[string]image.Image{}
		if i < 2 {
			shots["1.png"] = image.NewRGBA(image.Rect(0, 0, 4, 4))
		}
		if err := SaveRun(dir, rec, shots); err != nil {
			t.Fatal(err)
		}
	}
	runs := listRuns(dir)
	if len(runs) != maxRuns || runs[0].ID != NewRunID(start.Add((maxRuns+1)*time.Second)) {
		t.Fatalf("còn %d lần, mới nhất %v", len(runs), runs[0].ID)
	}
	if pngs, _ := filepath.Glob(filepath.Join(dir, "*.png")); len(pngs) != 0 {
		t.Errorf("ảnh của lần chạy đã xoá vẫn còn: %v", pngs)
	}
	if _, err := loadRun(dir, "../settings"); err == nil {
		t.Error("id có đường dẫn phải bị chặn")
	}
	if r, err := loadRun(dir, runs[0].ID); err != nil || r.Status != model.RunSuccess {
		t.Errorf("đọc lại: %+v %v", r, err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteRun(t *testing.T) {
	dir := RunsDir(t.TempDir(), "A", "/x/a.json")
	at := time.Date(2026, 10, 20, 9, 0, 0, 0, time.Local)
	var ids []string
	for i := range 2 {
		rec := model.RunRecord{RunSummary: model.RunSummary{ID: NewRunID(at.Add(time.Duration(i) * time.Second)), Status: model.RunFail}}
		if err := SaveRun(dir, rec, map[string]image.Image{"1.png": image.NewRGBA(image.Rect(0, 0, 2, 2))}); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, rec.ID)
	}
	for _, bad := range []string{"*", "../x", ".hidden", "a/b", `a\b`} {
		if err := deleteRun(dir, bad); err == nil {
			t.Errorf("id %q phải bị chặn", bad)
		}
	}
	if err := deleteRun(dir, ids[0]); err != nil {
		t.Fatal(err)
	}
	if runs := listRuns(dir); len(runs) != 1 || runs[0].ID != ids[1] {
		t.Errorf("còn %v, muốn chỉ %s", runs, ids[1])
	}
	if pngs, _ := filepath.Glob(filepath.Join(dir, "*.png")); len(pngs) != 1 || !strings.HasPrefix(filepath.Base(pngs[0]), ids[1]) {
		t.Errorf("ảnh còn lại %v", pngs)
	}
	if err := deleteRun(dir, ids[0]); err != nil {
		t.Errorf("xoá lần đã xoá không được lỗi: %v", err)
	}
}
