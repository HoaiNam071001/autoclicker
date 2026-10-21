package store

import (
	"os"
	"path/filepath"
	"strings"

	"ghosthand/internal/model"
	"ghosthand/internal/script"
)

// ReadScriptFile đọc file kịch bản (kiểm tra qua script.DecodeScript); file chưa có tên thì lấy theo tên file.
func ReadScriptFile(path string) (model.Script, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return model.Script{}, err
	}
	sc, err := script.DecodeScript(b)
	if err != nil {
		return model.Script{}, err
	}
	if sc.Name == "" {
		sc.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	return sc, nil
}

// WriteScriptFile ghi kịch bản xuống path (ghi file tạm rồi đổi tên).
func WriteScriptFile(path string, sc model.Script) error {
	b, err := script.EncodeScript(sc)
	if err != nil {
		return err
	}
	return WriteFileAtomic(path, b)
}
