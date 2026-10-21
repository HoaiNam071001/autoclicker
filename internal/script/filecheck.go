package script

import (
	"path/filepath"
	"strings"

	"ghosthand/internal/model"
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
