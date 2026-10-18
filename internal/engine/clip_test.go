package engine

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

func clipStep(f func(*model.Step)) model.Step { return modeltest.Step(model.StepClipboard, f) }

func TestClipboardReadWrite(t *testing.T) {
	clip := &devicetest.FakeClipboard{Text: "cũ"}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		clipStep(nil), // đọc
		clipStep(func(s *model.Step) { s.ClipOp, s.Text = model.ClipWrite, "mã {{n1.output.text}}" }),
	}}
	out, err := runShot(t, sc, &fakeScreen{}, clip)
	if err != nil {
		t.Fatal(err)
	}
	if clip.Text != "mã cũ" || out["text"] != "mã cũ" {
		t.Errorf("clipboard %q, output %v – muốn ghi “mã cũ” (biến từ lần đọc)", clip.Text, out)
	}
}

func TestClipboardImage(t *testing.T) {
	dir := t.TempDir()
	png1 := filepath.Join(dir, "a.png")
	f, _ := os.Create(png1)
	png.Encode(f, image.NewRGBA(image.Rect(0, 0, 2, 2)))
	f.Close()
	os.WriteFile(filepath.Join(dir, "b.jpg"), []byte("x"), 0o644)

	clip := &devicetest.FakeClipboard{}
	one := func(p string) (map[string]any, error) {
		return runShot(t, model.Script{Loops: 1, Steps: []model.Step{clipStep(func(s *model.Step) { s.ClipOp, s.FilePath = model.ClipImage, p })}}, &fakeScreen{}, clip)
	}
	if out, err := one(png1); err != nil || clip.Image != png1 || out["path"] != png1 {
		t.Errorf("ảnh PNG: %v, clipboard %q, output %v", err, clip.Image, out)
	}
	for _, c := range []struct{ path, want string }{
		{filepath.Join(dir, "b.jpg"), "PNG"},
		{filepath.Join(dir, "khong-co.png"), "không thấy file"},
	} {
		if _, err := one(c.path); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: lỗi %v, muốn chứa %q", c.path, err, c.want)
		}
	}
	// clipboard của máy lỗi → node lỗi
	clip.Err = errors.New("không ghi được clipboard")
	if _, err := one(png1); err == nil {
		t.Error("clipboard lỗi mà node vẫn thành công")
	}
}
