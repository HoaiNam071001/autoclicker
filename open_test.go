package main

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/model"
)

func TestOpenCommand(t *testing.T) {
	cases := []struct {
		goos, target, app string
		name              string
		args              []string
	}{
		{"darwin", "youtube.com/watch?v=1", "", "open", []string{"https://youtube.com/watch?v=1"}},
		{"darwin", "https://a.vn", "Google Chrome", "open", []string{"-a", "Google Chrome", "https://a.vn"}},
		{"darwin", "Spotify", "", "open", []string{"-a", "Spotify"}},
		{"darwin", "spotify:search:tinh", "", "open", []string{"spotify:search:tinh"}},
		{"darwin", "/Applications/Zalo.app", "", "open", []string{"/Applications/Zalo.app"}},
		{"windows", "https://a.vn/?q=1&p=2", "", "powershell", []string{"-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath 'https://a.vn/?q=1&p=2'"}},
		{"windows", "a.vn", "Google Chrome", "powershell", []string{"-NoProfile", "-NonInteractive", "-Command", `Start-Process -FilePath 'chrome' -ArgumentList '"https://a.vn"'`}},
		{"windows", `C:\Program Files\O'Brien\x.exe`, "", "powershell", []string{"-NoProfile", "-NonInteractive", "-Command", `Start-Process -FilePath 'C:\Program Files\O''Brien\x.exe'`}},
		{"linux", "a.vn", "", "xdg-open", []string{"https://a.vn"}},
	}
	for _, c := range cases {
		name, args := openCommand(c.goos, c.target, c.app)
		if name != c.name || !reflect.DeepEqual(args, c.args) {
			t.Errorf("%s %q %q: có %s %q, muốn %s %q", c.goos, c.target, c.app, name, args, c.name, c.args)
		}
	}
}

func TestTargetKind(t *testing.T) {
	for target, want := range map[string]string{
		"facebook.com/abc": targetURL, "https://x": targetURL, "mailto:a@b.c": targetURL, "localhost:3000": targetURL,
		"/Users/a/x.pdf": targetPath, `C:\x.exe`: targetPath, "~/x.pdf": targetPath,
		"Spotify": targetApp, "Visual Studio Code": targetApp,
	} {
		if got := targetKind(target); got != want {
			t.Errorf("%q: %s, muốn %s", target, got, want)
		}
	}
}

// Biến trong URL được mã hoá, trong tên app / văn bản thì giữ nguyên.
func TestResolveOpenTarget(t *testing.T) {
	sc := model.Script{
		Steps: []model.Step{
			{ID: "o", Type: model.StepOpen, Target: "youtube.com/results?search_query={{bai_hat}}", Enabled: true},
			{ID: "t", Type: model.StepText, Text: "{{bai_hat}}", Enabled: true},
		},
		Inputs: []model.ScriptInput{{Key: "bai_hat", Type: model.InputText}},
	}
	out, err := resolveInputs(sc, map[string]any{"bai_hat": "Tinh Hà Say 2 & bạn"})
	if err != nil {
		t.Fatal(err)
	}
	if got := out.Steps[0].Target; got != "youtube.com/results?search_query=Tinh%20H%C3%A0%20Say%202%20%26%20b%E1%BA%A1n" {
		t.Errorf("URL: %q", got)
	}
	if out.Steps[1].Text != "Tinh Hà Say 2 & bạn" {
		t.Errorf("gõ chữ không mã hoá: %q", out.Steps[1].Text)
	}
}

func TestRunScrollAndOpen(t *testing.T) {
	in := &fakeInput{}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		step(model.StepOpen, func(s *model.Step) { s.Target = "a.vn"; s.App = "Firefox" }),
		step(model.StepScroll, func(s *model.Step) { s.Amount = 2 }),
		step(model.StepScroll, func(s *model.Step) { s.Direction, s.Amount, s.Points = "up", 1, []model.Point{{X: 10, Y: 20}} }),
	}}
	e := newEngine(in, linear(sc), nil)
	var opened []string
	e.open = func(target, app string) error {
		opened = append(opened, target+" | "+app)
		return nil
	}
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"scroll down", "scroll down", "move 10,20", "scroll up"}; !reflect.DeepEqual(in.lines(), want) {
		t.Errorf("có %q, muốn %q", in.lines(), want)
	}
	if !reflect.DeepEqual(opened, []string{"a.vn | Firefox"}) {
		t.Errorf("mở: %q", opened)
	}
}

func TestValidateScrollOpen(t *testing.T) {
	cases := map[string]model.Step{
		"hướng lạ":        step(model.StepScroll, func(s *model.Step) { s.Direction = "xeo" }),
		"0 nấc":           step(model.StepScroll, func(s *model.Step) { s.Amount = 0 }),
		"chưa nhập URL":   step(model.StepOpen, nil),
		"điểm ngoài vùng": step(model.StepScroll, func(s *model.Step) { s.Points = []model.Point{{X: 5000, Y: 5000}} }),
	}
	for name, s := range cases {
		sc := model.Script{Region: &model.Rect{W: 100, H: 100}, Steps: []model.Step{s}}
		if err := validateScript(linear(sc), model.Rect{}); err == nil || !strings.HasPrefix(err.Error(), "n1:") {
			t.Errorf("%s: muốn lỗi ở n1, có %v", name, err)
		}
	}
}
