package script

import (
	"testing"

	"autoclicker/internal/model"
)

func TestDeleteOps(t *testing.T) {
	trash := model.Step{ID: "s1", Ref: "n1", Type: model.StepFile, Enabled: true, FileOp: model.FileTrash, FilePath: "~/Downloads", FilePattern: "*.png"}
	sc := model.Script{Steps: []model.Step{
		{ID: "s0", Ref: "n0", Type: model.StepFile, Enabled: true, FilePath: "~/a.txt"}, // ghi – không phải xoá
		trash,
		{ID: "s2", Ref: "n2", Type: model.StepFile, Enabled: false, FileOp: model.FileEmptyTrash}, // tắt – không chạy
	}}
	ops, fp := DeleteOps(sc)
	if len(ops) != 1 || ops[0].Folder != "~/Downloads" || ops[0].Pattern != "*.png" || fp == "" {
		t.Fatalf("ops %+v fp %q", ops, fp)
	}
	sc.Steps[1].FilePattern = "*"
	if _, fp2 := DeleteOps(sc); fp2 == fp {
		t.Error("đổi mẫu tên mà dấu vân tay không đổi – sẽ không hỏi lại")
	}
	sc.Steps[1].Name = "đổi tên"
	_, a := DeleteOps(sc)
	sc.Steps[1].Name = ""
	if _, b := DeleteOps(sc); a != b {
		t.Error("đổi tên node không được làm hỏi lại")
	}
	if ops, fp := DeleteOps(model.Script{Steps: sc.Steps[:1]}); ops != nil || fp != "" {
		t.Error("không có thao tác xoá mà vẫn có dấu vân tay")
	}
}
