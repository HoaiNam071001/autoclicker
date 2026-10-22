package engine

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

// runZip trả output + mã lỗi có cấu trúc ("" = không lỗi).
func runZip(t *testing.T, patch func(*model.Step)) (map[string]any, string, error) {
	t.Helper()
	s := modeltest.Step(model.StepZip, patch)
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	err := e.Run(t.Context())
	r := e.results[e.sc.Steps[0].ID]
	code := ""
	if r != nil && r.Error != nil {
		code = r.Error.Code
	}
	return r.Output, code, err
}

func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestZipRoundTripKeepsFolderAndFilter(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"bc/a.pdf": "AAA", "bc/sub/b.pdf": "BB", "bc/c.txt": "C"})
	out, code, err := runZip(t, func(s *model.Step) {
		s.FilePath, s.ZipTo, s.FilePattern = filepath.Join(dir, "bc"), filepath.Join(dir, "out", "bc.zip"), "*.pdf"
	})
	if err != nil || code != "" || out["count"] != 2 || out["size"] != 5 || out["path"] != filepath.Join(dir, "out", "bc.zip") {
		t.Fatalf("nén: out=%v err=%v", out, err)
	}
	dest := filepath.Join(dir, "giai")
	out, code, err = runZip(t, func(s *model.Step) {
		s.ZipOp, s.FilePath, s.ZipTo = model.ZipExtract, filepath.Join(dir, "out", "bc.zip"), dest
	})
	if err != nil || code != "" || out["count"] != 2 {
		t.Fatalf("giải nén: out=%v err=%v", out, err)
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "bc", "sub", "b.pdf")); string(b) != "BB" {
		t.Errorf("nội dung b.pdf = %q", b)
	}
	if _, err := os.Stat(filepath.Join(dest, "bc", "c.txt")); err == nil {
		t.Error("c.txt không khớp mẫu mà vẫn được nén")
	}
}

func TestZipDefaultDestinations(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"x.txt": "hi"})
	out, code, err := runZip(t, func(s *model.Step) { s.FilePath, s.ZipTo = filepath.Join(dir, "x.txt"), filepath.Join(dir, "pack.zip") })
	if err != nil || code != "" || out["path"] != filepath.Join(dir, "pack.zip") {
		t.Fatalf("out=%v err=%v", out, err)
	}
	out, code, err = runZip(t, func(s *model.Step) { s.ZipOp, s.FilePath = model.ZipExtract, filepath.Join(dir, "pack.zip") })
	if err != nil || code != "" || out["path"] != filepath.Join(dir, "pack") {
		t.Fatalf("out=%v err=%v", out, err)
	}
}

func TestZipDefaultNameIsSourcePlusZip(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"x.txt": "hi"})
	out, code, err := runZip(t, func(s *model.Step) { s.FilePath = filepath.Join(dir, "x.txt") })
	if err != nil || code != "" || out["path"] != filepath.Join(dir, "x.txt.zip") {
		t.Fatalf("out=%v err=%v", out, err)
	}
	// đích mặc định của giải nén trùng một file có sẵn → báo rõ, không ghi đè
	_, code, err = runZip(t, func(s *model.Step) { s.ZipOp, s.FilePath = model.ZipExtract, filepath.Join(dir, "x.txt.zip") })
	if code != model.ErrInvalidInput {
		t.Errorf("err=%v", err)
	}
}

func TestZipNothingToCompress(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"d/a.txt": "1"})
	_, code, err := runZip(t, func(s *model.Step) { s.FilePath, s.FilePattern = filepath.Join(dir, "d"), "*.png" })
	if code != model.ErrNotFound {
		t.Errorf("err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "d.zip")); err == nil {
		t.Error("để lại file zip rỗng")
	}
}

func TestUnzipRejectsEscapingEntries(t *testing.T) {
	for _, name := range []string{"../evil.txt", "a/../../evil.txt", "/abs/evil.txt", `..\evil.txt`} {
		dir := t.TempDir()
		zp := filepath.Join(dir, "bad.zip")
		f, _ := os.Create(zp)
		zw := zip.NewWriter(f)
		w, _ := zw.CreateHeader(&zip.FileHeader{Name: name})
		w.Write([]byte("x"))
		zw.Close()
		f.Close()
		_, code, err := runZip(t, func(s *model.Step) { s.ZipOp, s.FilePath, s.ZipTo = model.ZipExtract, zp, filepath.Join(dir, "out") })
		if code != model.ErrInvalidInput {
			t.Errorf("%q: err=%v", name, err)
		}
		if _, err := os.Stat(filepath.Join(dir, "evil.txt")); err == nil {
			t.Errorf("%q: file thoát ra ngoài", name)
		}
	}
}

func TestUnzipRejectsSymlinkEntry(t *testing.T) {
	dir := t.TempDir()
	zp := filepath.Join(dir, "link.zip")
	f, _ := os.Create(zp)
	zw := zip.NewWriter(f)
	h := &zip.FileHeader{Name: "l"}
	h.SetMode(os.ModeSymlink | 0o777)
	w, _ := zw.CreateHeader(h)
	w.Write([]byte("/etc"))
	zw.Close()
	f.Close()
	_, code, err := runZip(t, func(s *model.Step) { s.ZipOp, s.FilePath, s.ZipTo = model.ZipExtract, zp, filepath.Join(dir, "out") })
	if code != model.ErrInvalidInput {
		t.Errorf("err=%v", err)
	}
}

func TestUnzipMissingAndNotZip(t *testing.T) {
	dir := t.TempDir()
	_, code, err := runZip(t, func(s *model.Step) { s.ZipOp, s.FilePath = model.ZipExtract, filepath.Join(dir, "none.zip") })
	if code != model.ErrNotFound {
		t.Errorf("thiếu file: err=%v", err)
	}
	writeTree(t, dir, map[string]string{"fake.zip": "không phải zip"})
	_, code, err = runZip(t, func(s *model.Step) { s.ZipOp, s.FilePath = model.ZipExtract, filepath.Join(dir, "fake.zip") })
	if code != model.ErrInvalidInput {
		t.Errorf("không phải zip: err=%v", err)
	}
}
