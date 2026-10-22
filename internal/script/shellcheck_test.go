package script

import (
	"errors"
	"testing"

	"ghosthand/internal/model"
)

func shellStep(mod func(*model.Step)) model.Step {
	s := model.Step{ID: "s1", Ref: "n1", Type: model.StepShell, Enabled: true, Script: "echo hi", TimeoutMs: 1000}
	if mod != nil {
		mod(&s)
	}
	return s
}

func TestDangerOpsShell(t *testing.T) {
	sc := model.Script{Steps: []model.Step{shellStep(nil)}}
	ops, fp := DangerOps(sc)
	if len(ops) != 1 || ops[0].Op != DangerShell || ops[0].Script != "echo hi" || fp == "" {
		t.Fatalf("ops %+v fp %q", ops, fp)
	}
	sc.Steps[0].Script = "echo hi "
	if _, fp2 := DangerOps(sc); fp2 == fp {
		t.Error("sửa lệnh mà dấu vân tay không đổi – sẽ chạy lệnh mới không hỏi")
	}
	sc.Steps[0].Script = "echo hi"
	sc.Steps[0].ShellEnv = []model.KV{{Key: "A", Value: "{{vars.x}}"}}
	if _, fp3 := DangerOps(sc); fp3 != fp {
		t.Error("đổi giá trị biến môi trường không được làm hỏi lại")
	}
	sc.Steps[0].Enabled = false
	if ops, _ := DangerOps(sc); ops != nil {
		t.Error("node tắt không cần xác nhận")
	}
}

func TestValidateShell(t *testing.T) {
	for name, tc := range map[string]struct {
		mod  func(*model.Step)
		want string
	}{
		"trống":   {func(s *model.Step) { s.Script = "  " }, "script"},
		"tên env": {func(s *model.Step) { s.ShellEnv = []model.KV{{Key: "1A", Value: "x"}} }, "shellEnv[0].key"},
		"ok":      {func(s *model.Step) { s.ShellEnv = []model.KV{{Key: "A_1", Value: "x"}} }, ""},
	} {
		err := validateShell(0, shellStep(tc.mod))
		var se *model.ScriptError
		if tc.want == "" && err != nil || tc.want != "" && (!errors.As(err, &se) || se.Field != tc.want) {
			t.Errorf("%s: err = %v, muốn trường %q", name, err, tc.want)
		}
	}
	if err := checkShell(shellStep(func(s *model.Step) { s.ShellLang = "ruby" })); err == nil {
		t.Error("ngôn ngữ lạ phải bị từ chối")
	}
}

func TestVarFieldsShellSkipsScript(t *testing.T) {
	s := shellStep(func(s *model.Step) { s.ShellEnv = []model.KV{{Key: "A", Value: "{{vars.x}}"}} })
	for _, f := range VarFields(&s) {
		if f.Ptr == &s.Script {
			t.Fatal("Script không được nhận biến")
		}
	}
}
