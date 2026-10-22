package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

func downloadRun(t *testing.T, dir string, patch func(*model.Step), during func()) (map[string]any, error) {
	t.Helper()
	s := modeltest.Step(model.StepWaitDownload, func(s *model.Step) {
		s.FilePath, s.TimeoutMs, s.IntervalMs = dir, 1500, 10
		patch(s)
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	if during != nil {
		go func() { time.Sleep(60 * time.Millisecond); during() }()
	}
	err := e.Run(t.Context())
	if r := e.results[e.sc.Steps[0].ID]; r != nil {
		return r.Output, err
	}
	return nil, err
}

func write(t *testing.T, p, body string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Error(err)
	}
}

func TestWaitDownloadFindsNewFileAfterTempRename(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "cu.pdf"), "da co tu truoc")
	out, err := downloadRun(t, dir, func(s *model.Step) { s.FilePattern = "*.PDF" }, func() {
		write(t, filepath.Join(dir, "bao-cao.pdf.crdownload"), "do dang")
		time.Sleep(60 * time.Millisecond)
		os.Rename(filepath.Join(dir, "bao-cao.pdf.crdownload"), filepath.Join(dir, "bao-cao.pdf"))
	})
	if err != nil {
		t.Fatal(err)
	}
	if out["name"] != "bao-cao.pdf" || out["size"] != int64(len("do dang")) || out["path"] != filepath.Join(dir, "bao-cao.pdf") {
		t.Fatalf("output sai: %v", out)
	}
}

func TestWaitDownloadIgnoresOldAndNonMatching(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "cu.pdf"), "x")
	_, err := downloadRun(t, dir, func(s *model.Step) { s.FilePattern = "*.pdf"; s.TimeoutMs = 300 }, func() {
		write(t, filepath.Join(dir, "anh.png"), "y")
		write(t, filepath.Join(dir, "dang-tai.pdf.part"), "z")
	})
	if err == nil || !strings.Contains(err.Error(), "chưa thấy file tải xong") {
		t.Fatalf("muốn hết giờ, được %v", err)
	}
}

func TestWaitDownloadWaitsForSizeToSettle(t *testing.T) {
	dir := t.TempDir()
	out, err := downloadRun(t, dir, func(s *model.Step) { s.IntervalMs = 100 }, func() {
		p := filepath.Join(dir, "lon.zip")
		write(t, p, "a")
		time.Sleep(15 * time.Millisecond)
		write(t, p, "abcdef")
	})
	if err != nil || out["size"] != int64(6) {
		t.Fatalf("muốn size 6, được %v %v", out, err)
	}
}

func TestWaitDownloadMissingDir(t *testing.T) {
	_, err := downloadRun(t, filepath.Join(t.TempDir(), "khong-co"), func(*model.Step) {}, nil)
	if err == nil || !strings.Contains(err.Error(), "không phải thư mục tải về") {
		t.Fatalf("muốn lỗi thư mục, được %v", err)
	}
}
