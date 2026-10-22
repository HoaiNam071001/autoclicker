package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ghosthand/internal/model"
)

// Node Chờ tải xong (bước 43 B6): quét thư mục tải về (FilePath, trống = Downloads) mỗi IntervalMs; file MỚI (tên chưa có
// lúc bắt đầu node, hoặc đã có nhưng bị ghi lại) khớp FilePattern, không còn đuôi tạm và cỡ đứng yên 2 lần quét liền → xong.
// Quét định kỳ thay vì FSEvents / ReadDirectoryChangesW: cùng một code hai hệ, thư mục tải nhỏ nên rẻ.

// downloadTemp – đuôi file trình duyệt đang tải dở (Chrome / Edge, Firefox, Safari, Opera, tải qua tmp).
var downloadTemp = []string{".crdownload", ".part", ".tmp", ".download", ".opdownload", ".partial"}

type fileStamp struct {
	size int64
	mod  time.Time
}

func (e *Engine) waitDownload(ctx context.Context, s model.Step) (map[string]any, error) {
	dir, err := downloadDir(s.FilePath)
	if err != nil {
		return nil, err
	}
	pattern := strings.ToLower(strings.TrimSpace(s.FilePattern))
	if pattern == "" {
		pattern = "*"
	}
	if _, err := filepath.Match(pattern, ""); err != nil {
		return nil, &invalidErr{fmt.Sprintf("mẫu tên file %q sai cú pháp: %v", s.FilePattern, err)}
	}
	before := map[string]fileStamp{}
	for n, st := range scanDir(dir) {
		before[n] = st
	}
	last := map[string]int64{} // cỡ lần quét trước của file ứng viên
	var found string
	var size int64
	_, err = e.waitLoop(ctx, s, "download", func() (bool, int, error) {
		cur, newest := scanDir(dir), time.Time{}
		next := map[string]int64{}
		for n, st := range cur {
			if old, ok := before[n]; ok && old == st {
				continue
			}
			low := strings.ToLower(n)
			if isDownloadTemp(low) {
				continue
			}
			if ok, _ := filepath.Match(pattern, low); !ok {
				continue
			}
			next[n] = st.size
			if prev, seen := last[n]; seen && prev == st.size && st.size > 0 && !st.mod.Before(newest) {
				found, size, newest = n, st.size, st.mod
			}
		}
		last = next
		return found != "", 0, nil
	})
	if err != nil {
		if err == errTimeout {
			return nil, &timeoutErr{msg: fmt.Sprintf("chưa thấy file tải xong trong %s sau %s", dir, fmtSeconds(s.TimeoutMs)), detail: map[string]any{"dir": dir, "pattern": s.FilePattern}}
		}
		return nil, err
	}
	return map[string]any{"path": filepath.Join(dir, found), "name": found, "size": size}, nil
}

// downloadDir – thư mục theo dõi: trống = Downloads của người dùng (cùng tên trên macOS / Windows).
func downloadDir(p string) (string, error) {
	if strings.TrimSpace(p) == "" {
		p = "~/Downloads"
	}
	dir, err := userPath(p)
	if err != nil {
		return "", err
	}
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return "", &invalidErr{fmt.Sprintf("%s không phải thư mục tải về", dir)}
	}
	return dir, nil
}

// scanDir – file thường (bỏ file ẩn) → cỡ + giờ sửa; thư mục đọc lỗi coi như trống, lần quét sau thử lại.
func scanDir(dir string) map[string]fileStamp {
	out := map[string]fileStamp{}
	entries, _ := os.ReadDir(dir)
	for _, en := range entries {
		if !en.Type().IsRegular() || strings.HasPrefix(en.Name(), ".") {
			continue
		}
		if info, err := en.Info(); err == nil {
			out[en.Name()] = fileStamp{info.Size(), info.ModTime()}
		}
	}
	return out
}

func isDownloadTemp(lowerName string) bool {
	for _, x := range downloadTemp {
		if strings.HasSuffix(lowerName, x) {
			return true
		}
	}
	return false
}
