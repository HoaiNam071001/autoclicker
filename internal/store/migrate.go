package store

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Chuyển dữ liệu từ tên cũ Auto Clicker sang Ghosthand (bước 45). Cầu nối tạm – xoá ở bước phát hành đầu tiên cho
// người khác (18b) hoặc sau 2026-12-21, cùng lúc bỏ lời gọi trong main.go.

const (
	oldConfigName  = "autoclicker"
	oldUserDirName = "Auto Clicker"
	newConfigName  = "ghosthand"
)

type dirMove struct{ from, to string }

// MigrateOldName – gọi một lần lúc mở app, trước mọi đọc / ghi cấu hình. Lỗi chỉ ghi log, không chặn mở app.
func MigrateOldName() {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return
	}
	moves := []dirMove{{filepath.Join(cfg, oldConfigName), filepath.Join(cfg, newConfigName)}}
	var scripts dirMove
	if home, err := os.UserHomeDir(); err == nil {
		docs, pics := knownFolders(home)
		scripts = dirMove{filepath.Join(docs, oldUserDirName), filepath.Join(docs, userDirName)}
		moves = append(moves, scripts, dirMove{filepath.Join(pics, oldUserDirName), filepath.Join(pics, userDirName)})
	}
	moved := moveOldDirs(moves)
	if len(moved) == 0 {
		return
	}
	newCfg := moves[0].to
	if err := rewritePaths(newCfg, moved); err != nil {
		log.Printf("đổi đường dẫn cũ trong cấu hình: %v", err)
	}
	for _, m := range moved {
		if m == scripts {
			moveRunHistory(newCfg, m)
		}
	}
	if cache, err := os.UserCacheDir(); err == nil {
		os.RemoveAll(filepath.Join(cache, oldConfigName)) // bộ Tesseract giải nén theo tên cũ – lần đọc chữ sau giải nén lại
	}
}

// moveOldDirs đổi tên thư mục cũ sang mới khi cũ còn mà mới chưa có; trả các cặp đã chuyển. Lỗi (đang mở, khác ổ) → để nguyên.
func moveOldDirs(moves []dirMove) []dirMove {
	var done []dirMove
	for _, m := range moves {
		if st, err := os.Stat(m.from); err != nil || !st.IsDir() {
			continue
		}
		if _, err := os.Stat(m.to); err == nil {
			continue
		}
		if err := os.Rename(m.from, m.to); err != nil {
			log.Printf("chuyển %s → %s: %v", m.from, m.to, err)
			continue
		}
		done = append(done, m)
	}
	return done
}

// rewritePaths thay đường dẫn cũ → mới trong các file .json ngay trong dir (cài đặt, Mở gần đây, nháp, lịch hẹn).
// So trên dạng đã mã hoá JSON (Windows: "\\") và chỉ khớp cả tên thư mục (theo sau là dấu tách hoặc hết chuỗi).
func rewritePaths(dir string, moved []dirMove) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return err
	}
	var pairs []string
	for _, m := range moved {
		from, to := jsonInner(m.from), jsonInner(m.to)
		sep := jsonInner(string(filepath.Separator))
		pairs = append(pairs, from+sep, to+sep, from+`"`, to+`"`)
	}
	r := strings.NewReplacer(pairs...)
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if out := r.Replace(string(b)); out != string(b) {
			if err := os.WriteFile(f, []byte(out), 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

// jsonInner – chuỗi như khi nằm trong JSON (không có ngoặc kép hai đầu).
func jsonInner(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}

// moveRunHistory – lịch sử chạy đặt theo hash đường dẫn kịch bản (RunsDir): kịch bản trong thư mục đã chuyển có đường dẫn
// mới → đổi tên thư mục lịch sử tương ứng để không mất.
func moveRunHistory(configDir string, m dirMove) {
	files, _ := filepath.Glob(filepath.Join(m.to, "*.json"))
	for _, f := range files {
		var sc struct {
			Name string `json:"name"`
		}
		if b, err := os.ReadFile(f); err != nil || json.Unmarshal(b, &sc) != nil {
			continue
		}
		oldPath := filepath.Join(m.from, filepath.Base(f))
		moveOldDirs([]dirMove{{RunsDir(configDir, sc.Name, oldPath), RunsDir(configDir, sc.Name, f)}})
	}
}
