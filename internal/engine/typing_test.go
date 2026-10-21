package engine

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/model"
	"ghosthand/internal/script/scripttest"
)

func TestUsePaste(t *testing.T) {
	cases := []struct {
		mode             string
		pasteAuto, human bool
		want             bool
	}{
		{model.TypeAuto, true, false, true}, // Windows, không giống người → dán
		{model.TypeAuto, true, true, false}, // giống người cần nhịp gõ
		{model.TypeAuto, false, false, false},
		{model.TypePaste, false, true, true},
		{model.TypeKeys, true, false, false},
	}
	for _, c := range cases {
		if got := usePaste(c.mode, c.pasteAuto, c.human); got != c.want {
			t.Errorf("usePaste(%q, %v, %v) = %v", c.mode, c.pasteAuto, c.human, got)
		}
	}
}

// clipLog – clipboard giả ghi lại mọi lần ghi chữ.
type clipLog struct {
	devicetest.FakeClipboard
	writes []string
}

func (c *clipLog) WriteText(t string) error {
	c.writes = append(c.writes, t)
	return c.FakeClipboard.WriteText(t)
}

func runTyping(t *testing.T, s model.Step, setup func(e *Engine)) (*devicetest.FakeInput, model.StepResult) {
	t.Helper()
	in := &devicetest.FakeInput{}
	e := newEngine(in, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	var res model.StepResult
	e.onNode = func(r model.StepResult) { res = r }
	setup(e)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	return in, res
}

// Dán: ghi chữ vào clipboard, Ctrl/Cmd+V, rồi trả chữ cũ – không gõ phím nào (bộ gõ tiếng Việt không xen vào được).
func TestTypePasteRestoresClipboard(t *testing.T) {
	clip := &clipLog{FakeClipboard: devicetest.FakeClipboard{Text: "chữ cũ"}}
	s := model.Step{ID: "s1", Type: model.StepText, Text: "mình thử", TypeMode: model.TypePaste, Enabled: true}
	in, _ := runTyping(t, s, func(e *Engine) { e.clip = clip })
	mod := shortcutMod()
	want := []string{"down " + mod + " []", "down v [" + mod + "]", "up v", "up " + mod}
	if got := in.Lines(); !reflect.DeepEqual(got, want) {
		t.Errorf("phím: có %q, muốn %q", got, want)
	}
	if !reflect.DeepEqual(clip.writes, []string{"mình thử", "chữ cũ"}) {
		t.Errorf("clipboard: %q", clip.writes)
	}
}

// Tự chọn trên Windows (pasteAuto) mà bật giống người → vẫn gõ phím; thấy UniKey thì nhắc trong nhật ký.
func TestTypeKeysWarnsVietnameseIME(t *testing.T) {
	s := model.Step{ID: "s1", Type: model.StepText, Text: "ab", Enabled: true}
	in, res := runTyping(t, s, func(e *Engine) {
		e.imeCheck = true
		e.apps = &devicetest.FakeApps{Apps: []model.AppInfo{{Name: "explorer"}, {Name: "UniKeyNT"}}}
	})
	if got := in.Lines(); !reflect.DeepEqual(got, []string{"type a", "type b"}) {
		t.Errorf("phím: %q", got)
	}
	if len(res.Logs) != 1 || res.Logs[0].Level != "warn" || !strings.Contains(res.Logs[0].Msg, "UniKeyNT") {
		t.Errorf("nhật ký: %+v", res.Logs)
	}
}
