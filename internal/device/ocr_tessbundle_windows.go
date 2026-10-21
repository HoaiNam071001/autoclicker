package device

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Bộ Tesseract kèm trong exe (bước 40 O2): tess/tess.zip = DLL đã strip + vie.traineddata, do scripts/fetch-tesseract.sh
// tạo (không commit). Giải nén một lần vào thư mục cache theo hash – đổi phiên bản là thư mục mới. Build không có zip
// (chưa chạy script) vẫn chạy được: tessBundleDir báo lỗi → OCR Windows.

//go:embed tess
var tessFS embed.FS

func tessBundleDir() (string, error) {
	data, err := tessFS.ReadFile("tess/tess.zip")
	if err != nil {
		return "", errors.New("bản build không kèm Tesseract (chạy scripts/fetch-tesseract.sh trước khi build)")
	}
	sum := sha256.Sum256(data)
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "autoclicker", "tesseract-"+hex.EncodeToString(sum[:6]))
	if _, err := os.Stat(filepath.Join(dir, ".ok")); err == nil {
		return dir, nil
	}
	// giải nén vào thư mục tạm rồi đổi tên – hai tiến trình (app + CLI) cùng giải nén không đè nhau dở dang
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp(filepath.Dir(dir), "tesseract-tmp-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	if err := unzipTo(data, tmp); err != nil {
		return "", fmt.Errorf("giải nén Tesseract: %w", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".ok"), nil, 0o644); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, dir); err != nil {
		if _, e := os.Stat(filepath.Join(dir, ".ok")); e != nil {
			return "", err
		}
	}
	return dir, nil
}

func unzipTo(data []byte, dir string) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		name := filepath.Base(f.Name) // gói phẳng – bỏ mọi đường dẫn con
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(filepath.Join(dir, name))
		if err == nil {
			_, err = io.Copy(out, rc)
			if e := out.Close(); err == nil {
				err = e
			}
		}
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
