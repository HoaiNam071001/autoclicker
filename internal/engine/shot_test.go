package engine

import (
	"context"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

// runShot chạy kịch bản, trả output của node cuối + lỗi.
func runShot(t *testing.T, sc model.Script, scr *fakeScreen, clip *devicetest.FakeClipboard, displays ...model.Rect) (map[string]any, error) {
	t.Helper()
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), nil)
	e.scr, e.clip, e.displays = scr, clip, displays
	var out map[string]any
	e.onNode = func(r model.StepResult) { out = r.Output }
	err := e.Run(context.Background())
	return out, err
}

func shotStep(f func(*model.Step)) model.Step {
	return modeltest.Step(model.StepScreenshot, f)
}

func TestScreenshotSavesPNG(t *testing.T) {
	dir := t.TempDir()
	frame := image.NewRGBA(image.Rect(0, 0, 800, 600)) // Retina: ảnh gấp đôi vùng
	scr := &fakeScreen{frames: []image.Image{frame}}
	clip := &devicetest.FakeClipboard{}
	sc := model.Script{Loops: 1, Region: &model.Rect{X: 10, Y: 20, W: 400, H: 300},
		Steps: []model.Step{shotStep(func(s *model.Step) { s.SaveDir, s.FileName, s.ToClipboard = dir, "đơn/hàng 1", true })}}
	out, err := runShot(t, sc, scr, clip)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "đơn-hàng 1.png")
	if out["path"] != want || out["x"] != 10 || out["y"] != 20 || out["w"] != 400 || out["h"] != 300 {
		t.Errorf("output %v, muốn path %s + vùng kịch bản", out, want)
	}
	if scr.shots[0] != *sc.Region {
		t.Errorf("chụp %v, muốn vùng kịch bản", scr.shots[0])
	}
	f, err := os.Open(want)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if cfg, err := png.DecodeConfig(f); err != nil || cfg.Width != 800 || cfg.Height != 600 {
		t.Errorf("PNG %+v %v, muốn 800×600 (giữ px thật)", cfg, err)
	}
	if clip.Image != want {
		t.Errorf("clipboard %q, muốn %q", clip.Image, want)
	}
	// chạy lại cùng tên → không ghi đè
	out, _ = runShot(t, sc, scr, clip)
	if out["path"] != filepath.Join(dir, "đơn-hàng 1 (2).png") {
		t.Errorf("lần 2: %v, muốn thêm (2)", out["path"])
	}
	if left, _ := filepath.Glob(filepath.Join(dir, ".shot-*")); len(left) > 0 {
		t.Errorf("còn file tạm %v", left)
	}
}

func TestScreenshotAreas(t *testing.T) {
	ds := []model.Rect{{X: 0, Y: 0, W: 1440, H: 900}, {X: 1440, Y: 0, W: 1920, H: 1080}}
	frame := image.NewRGBA(image.Rect(0, 0, 10, 10))
	dir := t.TempDir()
	for _, c := range []struct {
		name string
		f    func(*model.Step)
		want model.Rect
	}{
		{"màn chứa vùng", func(s *model.Step) { s.ShotArea = model.ShotScreen }, ds[0]},
		{"màn thứ 2", func(s *model.Step) { s.ShotArea, s.Display = model.ShotScreen, 2 }, ds[1]},
		{"vùng cố định", func(s *model.Step) {
			s.ShotArea, s.ShotRect = model.ShotRect, &model.Rect{X: 1500, Y: 40, W: 300, H: 200}
		}, model.Rect{X: 1500, Y: 40, W: 300, H: 200}},
	} {
		scr := &fakeScreen{frames: []image.Image{frame}}
		sc := model.Script{Loops: 1, Region: &model.Rect{X: 100, Y: 100, W: 400, H: 300},
			Steps: []model.Step{shotStep(func(s *model.Step) { c.f(s); s.SaveDir = dir })}}
		if _, err := runShot(t, sc, scr, &devicetest.FakeClipboard{}, ds...); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if scr.shots[0] != c.want {
			t.Errorf("%s: chụp %v, muốn %v", c.name, scr.shots[0], c.want)
		}
	}
	// màn không có → lỗi rõ, không chụp
	scr := &fakeScreen{frames: []image.Image{frame}}
	sc := model.Script{Loops: 1, Steps: []model.Step{shotStep(func(s *model.Step) { s.ShotArea, s.Display, s.SaveDir = model.ShotScreen, 3, dir })}}
	if _, err := runShot(t, sc, scr, &devicetest.FakeClipboard{}, ds...); err == nil || !strings.Contains(err.Error(), "màn hình thứ 3") {
		t.Errorf("muốn lỗi không có màn 3, có %v", err)
	}
	// vắt qua 2 màn → lỗi
	sc.Steps[0].ShotArea, sc.Steps[0].ShotRect = model.ShotRect, &model.Rect{X: 1400, Y: 0, W: 100, H: 100}
	if _, err := runShot(t, sc, scr, &devicetest.FakeClipboard{}, ds...); err == nil || len(scr.shots) > 0 {
		t.Errorf("vắt 2 màn: lỗi %v, đã chụp %v", err, scr.shots)
	}
}

// Chụp theo khung chỗ thấy của node chờ ảnh trước đó.
func TestScreenshotFromNode(t *testing.T) {
	with, tpl := screenWithButton(1, true)
	scr := &fakeScreen{frames: []image.Image{with}}
	sc := imageScript(tpl, 1)
	sc.Steps = append(sc.Steps, shotStep(func(s *model.Step) { s.ShotArea, s.ImageFrom, s.SaveDir = model.ShotNode, sc.Steps[0].ID, t.TempDir() }))
	out, err := runShot(t, sc, scr, &devicetest.FakeClipboard{})
	if err != nil {
		t.Fatal(err)
	}
	if want := (model.Rect{X: 250, Y: 120, W: 60, H: 20}); scr.shots[len(scr.shots)-1] != want || out["w"] != 60 {
		t.Errorf("chụp %v (output %v), muốn khung nút %v", scr.shots[len(scr.shots)-1], out, want)
	}
}

func TestShotPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	now := time.Date(2026, 10, 18, 9, 5, 7, 0, time.Local)
	set := filepath.Join(home, "Anh chup")
	for _, c := range []struct{ dir, def, name, want string }{
		{"", "", "", filepath.Join(home, "Pictures", "Ghosthand", "ghosthand-2026-10-18-090507.png")},
		{" ", set, "x", filepath.Join(set, "x.png")},
		{"~/Desktop", set, "a.png", filepath.Join(home, "Desktop", "a.png")},
		{"Documents/x", "", "b:c", filepath.Join(home, "Documents", "x", "b-c.png")},
	} {
		if got, err := shotPath(c.dir, c.def, c.name, now); err != nil || got != c.want {
			t.Errorf("shotPath(%q, %q) = %q %v, muốn %q", c.dir, c.name, got, err, c.want)
		}
	}
}
