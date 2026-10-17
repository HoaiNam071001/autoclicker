package main

import (
	"context"
	"reflect"
	"strings"
	"testing"
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
	sc := Script{
		Steps: []Step{
			{ID: "o", Type: StepOpen, Target: "youtube.com/results?search_query={{bai_hat}}", Enabled: true},
			{ID: "t", Type: StepText, Text: "{{bai_hat}}", Enabled: true},
		},
		Inputs: []ScriptInput{{Key: "bai_hat", Type: InputText}},
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
	sc := Script{Loops: 1, Steps: []Step{
		step(StepOpen, func(s *Step) { s.Target = "a.vn"; s.App = "Firefox" }),
		step(StepScroll, func(s *Step) { s.Amount = 2 }),
		step(StepScroll, func(s *Step) { s.Direction, s.Amount, s.Points = "up", 1, []Point{{X: 10, Y: 20}} }),
	}}
	e := newEngine(in, sc, nil)
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
	cases := map[string]Step{
		"hướng lạ":        step(StepScroll, func(s *Step) { s.Direction = "xeo" }),
		"0 nấc":           step(StepScroll, func(s *Step) { s.Amount = 0 }),
		"chưa nhập URL":   step(StepOpen, nil),
		"điểm ngoài vùng": step(StepScroll, func(s *Step) { s.Points = []Point{{X: 5000, Y: 5000}} }),
	}
	for name, s := range cases {
		sc := Script{Region: &Rect{W: 100, H: 100}, Steps: []Step{s}}
		if err := validateScript(sc, Rect{}); err == nil || !strings.HasPrefix(err.Error(), "bước 1:") {
			t.Errorf("%s: muốn lỗi bước 1, có %v", name, err)
		}
	}
}
