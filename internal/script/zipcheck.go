package script

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"ghosthand/internal/model"
)

// checkZip – miền giá trị node Nén / giải nén khi mở file.
func checkZip(s model.Step) error {
	if !slices.Contains(model.ZipOps, s.ZipOp) {
		return fmt.Errorf("thao tác nén không hợp lệ: %q (trống = nén | extract)", s.ZipOp)
	}
	return nil
}

// validateZip – đủ dữ liệu để chạy. Ô còn biến thì phần còn lại kiểm tra lúc chạy.
func validateZip(i int, s model.Step) error {
	p := strings.TrimSpace(s.FilePath)
	if p == "" {
		if s.ZipOp == model.ZipExtract {
			return stepErr(i, s, "filePath", "chưa nhập file .zip cần giải nén")
		}
		return stepErr(i, s, "filePath", "chưa nhập file / thư mục cần nén")
	}
	if s.ZipOp == model.ZipExtract {
		if ext := strings.ToLower(filepath.Ext(p)); !strings.Contains(p, "{{") && ext != ".zip" {
			return stepErr(i, s, "filePath", "giải nén chỉ cho file .zip (file này: %q)", filepath.Base(p))
		}
		return nil
	}
	if to := strings.TrimSpace(s.ZipTo); to != "" && !strings.Contains(to, "{{") && strings.ToLower(filepath.Ext(to)) != ".zip" {
		return stepErr(i, s, "zipTo", "file nén phải có đuôi .zip (đang là %q)", filepath.Base(to))
	}
	if pat := strings.TrimSpace(s.FilePattern); pat != "" && !strings.Contains(pat, "{{") {
		if strings.ContainsAny(pat, `/\`) {
			return stepErr(i, s, "filePattern", "mẫu tên chỉ là tên file (vd *.png) – thư mục đặt ở ô nguồn")
		}
		if _, err := filepath.Match(pat, ""); err != nil {
			return stepErr(i, s, "filePattern", "mẫu tên “%s” sai dạng", pat)
		}
	}
	return nil
}
