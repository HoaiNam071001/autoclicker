package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Nơi lưu kịch bản / ảnh chụp (bước 41): người dùng chọn trong Cài đặt → Lưu trữ; trống = thư mục mặc định
// Documents/Auto Clicker (kịch bản), Pictures/Auto Clicker (ảnh). Thư mục cấu hình (ConfigDir) không đổi.

const userDirName = "Auto Clicker"

// StorageSettings – thư mục đã chọn; "" = mặc định. AIScriptsDir trống = theo ScriptsDir.
type StorageSettings struct {
	ScriptsDir   string `json:"scriptsDir,omitempty"`
	AIScriptsDir string `json:"aiScriptsDir,omitempty"`
	ShotsDir     string `json:"shotsDir,omitempty"`
}

// StorageDefaults – thư mục mặc định (Cài đặt hiện khi ô để trống).
type StorageDefaults struct {
	ScriptsDir string `json:"scriptsDir"`
	ShotsDir   string `json:"shotsDir"`
}

// DefaultDirs – Documents / Pictures của máy (Windows qua SHGetKnownFolderPath) + "Auto Clicker".
func DefaultDirs() (StorageDefaults, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return StorageDefaults{}, fmt.Errorf("không tìm thấy thư mục người dùng: %w", err)
	}
	docs, pics := knownFolders(home)
	return StorageDefaults{ScriptsDir: filepath.Join(docs, userDirName), ShotsDir: filepath.Join(pics, userDirName)}, nil
}

// resolveDir – thư mục dùng thật: custom trống → def (tạo nếu chưa có); custom phải còn và là thư mục – không tự lặng lẽ
// đổi sang chỗ khác (người dùng tưởng file vào chỗ đã chọn).
func resolveDir(custom, def string) (string, error) {
	if custom == "" {
		return def, os.MkdirAll(def, 0o755)
	}
	if st, err := os.Stat(custom); err != nil || !st.IsDir() {
		return "", fmt.Errorf("không thấy thư mục %s (đã xoá hoặc ổ đã rút?) – chọn lại hoặc bấm Về mặc định ở Cài đặt → Lưu trữ", custom)
	}
	return custom, nil
}

func (s StorageSettings) scripts(d StorageDefaults) (string, error) {
	return resolveDir(s.ScriptsDir, d.ScriptsDir)
}

func (s StorageSettings) aiScripts(d StorageDefaults) (string, error) {
	if s.AIScriptsDir == "" {
		return s.scripts(d)
	}
	return resolveDir(s.AIScriptsDir, d.ScriptsDir)
}

func currentStorage() (StorageSettings, StorageDefaults, error) {
	d, err := DefaultDirs()
	return CurrentSettings().Storage, d, err
}

// ScriptsDir – thư mục mở sẵn khi Lưu lần đầu / Mở kịch bản.
func ScriptsDir() (string, error) {
	s, d, err := currentStorage()
	if err != nil {
		return "", err
	}
	return s.scripts(d)
}

// AIScriptsDir – thư mục cho kịch bản AI tạo qua MCP (create_script, bản sao khi AI sửa kịch bản đang dở).
func AIScriptsDir() (string, error) {
	s, d, err := currentStorage()
	if err != nil {
		return "", err
	}
	return s.aiScripts(d)
}

// ShotsDir – thư mục mặc định của node Chụp màn hình (saveDir trống). Không kiểm tra còn hay không – lúc ghi ảnh mới báo.
func ShotsDir() string {
	s, d, _ := currentStorage()
	if s.ShotsDir != "" {
		return s.ShotsDir
	}
	return d.ShotsDir
}

// SetStorage lưu thư mục đã chọn; thư mục phải là đường dẫn tuyệt đối và đang có.
func SetStorage(s StorageSettings) error {
	for _, p := range []string{s.ScriptsDir, s.AIScriptsDir, s.ShotsDir} {
		if p == "" {
			continue
		}
		if !filepath.IsAbs(p) {
			return fmt.Errorf("thư mục phải là đường dẫn đầy đủ: %s", p)
		}
		if _, err := resolveDir(p, ""); err != nil {
			return err
		}
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	cur := loadSettings(dir)
	cur.Storage = s
	return saveSettings(dir, cur)
}

// FreePath – dir/name.json; đã có file thì thêm " 2", " 3"… (không ghi đè kịch bản cũ).
func FreePath(dir, name string) string {
	p := filepath.Join(dir, name+".json")
	for i := 2; ; i++ {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return p
		}
		p = filepath.Join(dir, name+" "+strconv.Itoa(i)+".json")
	}
}
