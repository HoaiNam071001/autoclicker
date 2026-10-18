package samples

import (
	"embed"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Trang sân thử: HTML + CSS + JS thuần, offline, tự hiện kết quả (đã click mấy lần, kéo từ đâu tới đâu…).
// Ghi ra đĩa và mở bằng file:// thay vì server HTTP cục bộ: không chiếm cổng, bản sao mẫu đã lưu vẫn mở được lần sau.

//go:embed pages
var pageFS embed.FS

// WriteSandbox ghi (đè) mọi trang sân thử vào dir – luôn khớp bản app đang chạy.
func WriteSandbox(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return fs.WalkDir(pageFS, "pages", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := pageFS.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dir, filepath.Base(p)), b, 0o644)
	})
}

// FileURL – URL file:// của thư mục (không có / cuối); dấu cách → %20, ổ Windows C:\ → file:///C:/.
func FileURL(dir string) string {
	p := filepath.ToSlash(dir)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return (&url.URL{Scheme: "file", Path: p}).String()
}
