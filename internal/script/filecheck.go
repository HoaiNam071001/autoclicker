package script

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"strings"

	"autoclicker/internal/model"
)

// checkFileStep – node Tệp / thư mục đủ dữ liệu để chạy (bước 35 N4). Ô còn biến thì phần còn lại kiểm tra lúc chạy.
func checkFileStep(i int, s model.Step) error {
	p := strings.TrimSpace(s.FilePath)
	if s.FileOp != model.FileEmptyTrash && p == "" {
		what := "file"
		if s.FileOp == model.FileMkdir || s.FileOp == model.FileTrash {
			what = "thư mục"
		}
		return stepErr(i, s, "filePath", "chưa nhập %s", what)
	}
	switch s.FileOp {
	case model.FileRow:
		if ext := strings.ToLower(filepath.Ext(p)); !strings.Contains(p, "{{") && ext != ".csv" && ext != ".xlsx" {
			return stepErr(i, s, "filePath", "thêm dòng chỉ cho file .csv hoặc .xlsx (file này: %q)", filepath.Base(p))
		}
		if len(s.Cells) == 0 {
			return stepErr(i, s, "cells", "chưa có ô nào cho dòng mới")
		}
	case model.FileTrash:
		pat := strings.TrimSpace(s.FilePattern)
		if pat == "" {
			return stepErr(i, s, "filePattern", "chưa nhập mẫu tên file cần xoá, vd *.png")
		}
		if strings.ContainsAny(pat, `/\`) {
			return stepErr(i, s, "filePattern", "mẫu tên chỉ là tên file (vd *.png) – thư mục đặt ở ô Thư mục")
		}
		if _, err := filepath.Match(pat, ""); err != nil {
			return stepErr(i, s, "filePattern", "mẫu tên “%s” sai dạng", pat)
		}
	}
	return nil
}

// DeleteOp – một thao tác xoá trong kịch bản (node Tệp / thư mục trash / emptyTrash) – hiện cho người dùng xác nhận.
type DeleteOp struct {
	Ref     string `json:"ref"`
	Op      string `json:"op"`      // trash | emptyTrash
	Folder  string `json:"folder"`  // trash: thư mục (có thể còn biến node)
	Pattern string `json:"pattern"` // trash: mẫu tên file
}

// DeleteOps – thao tác xoá của các node đang bật, theo thứ tự trong mảng (input đã thay). Fingerprint đổi khi thư mục /
// mẫu tên / thao tác đổi → phải xác nhận lại.
func DeleteOps(sc model.Script) (ops []DeleteOp, fingerprint string) {
	for _, s := range sc.Steps {
		if !s.Enabled || s.Type != model.StepFile || (s.FileOp != model.FileTrash && s.FileOp != model.FileEmptyTrash) {
			continue
		}
		op := DeleteOp{Ref: s.Ref, Op: s.FileOp}
		if s.FileOp == model.FileTrash {
			op.Folder, op.Pattern = strings.TrimSpace(s.FilePath), strings.TrimSpace(s.FilePattern)
		}
		ops = append(ops, op)
	}
	if len(ops) == 0 {
		return nil, ""
	}
	b, _ := json.Marshal(ops)
	h := sha256.Sum256(b)
	return ops, hex.EncodeToString(h[:8])
}
