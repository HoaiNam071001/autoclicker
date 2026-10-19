package script

import (
	"testing"

	"autoclicker/internal/model"
)

func TestDangerOps(t *testing.T) {
	trash := model.Step{ID: "s1", Ref: "n1", Type: model.StepFile, Enabled: true, FileOp: model.FileTrash, FilePath: "~/Downloads", FilePattern: "*.png"}
	sc := model.Script{Steps: []model.Step{
		{ID: "s0", Ref: "n0", Type: model.StepFile, Enabled: true, FilePath: "~/a.txt"}, // ghi – không phải xoá
		trash,
		{ID: "s2", Ref: "n2", Type: model.StepFile, Enabled: false, FileOp: model.FileEmptyTrash}, // tắt – không chạy
	}}
	ops, fp := DangerOps(sc)
	if len(ops) != 1 || ops[0].Folder != "~/Downloads" || ops[0].Pattern != "*.png" || fp == "" {
		t.Fatalf("ops %+v fp %q", ops, fp)
	}
	sc.Steps[1].FilePattern = "*"
	if _, fp2 := DangerOps(sc); fp2 == fp {
		t.Error("đổi mẫu tên mà dấu vân tay không đổi – sẽ không hỏi lại")
	}
	sc.Steps[1].Name = "đổi tên"
	_, a := DangerOps(sc)
	sc.Steps[1].Name = ""
	if _, b := DangerOps(sc); a != b {
		t.Error("đổi tên node không được làm hỏi lại")
	}
	if ops, fp := DangerOps(model.Script{Steps: sc.Steps[:1]}); ops != nil || fp != "" {
		t.Error("không có thao tác xoá mà vẫn có dấu vân tay")
	}
}

func TestDangerOpsAppWindow(t *testing.T) {
	quit := model.Step{Ref: "n1", Type: model.StepApp, Enabled: true, AppOp: model.AppQuit, AppName: "notepad"}
	closeWin := model.Step{Ref: "n2", Type: model.StepWindow, Enabled: true, WinOp: model.WinClose, WinApp: "notepad", WinTitle: "Không tên"}
	sc := model.Script{Steps: []model.Step{quit, closeWin, {Ref: "n3", Type: model.StepApp, Enabled: true, AppOp: model.AppForceQuit, AppName: "zalo"}}}
	ops, fp := DangerOps(sc)
	if len(ops) != 2 || ops[0].Op != DangerClose || ops[0].Target != "notepad “Không tên”" || ops[1].Op != DangerForceQuit || ops[1].Target != "zalo" {
		t.Fatalf("thoát không buộc không phải nguy hiểm: %+v", ops)
	}
	sc.Steps[0].AppForce = true
	if ops, fp2 := DangerOps(sc); len(ops) != 3 || ops[0].Op != DangerQuit || fp2 == fp {
		t.Errorf("thoát + buộc thoát khi quá giờ: %+v", ops)
	}
}
