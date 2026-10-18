// Package store – dữ liệu app lưu trên máy (không thuộc file kịch bản): cấu hình, file mở gần đây, giá trị input lần
// trước, lịch sử chạy. Mọi thứ trong ConfigDir().
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Cấu hình app (không thuộc kịch bản), lưu ở os.UserConfigDir()/autoclicker/settings.json.
// Bản nháp tự lưu (draft.json, storage.go ở root) nằm cùng thư mục.

const (
	ThemeLight  = "light"
	ThemeDark   = "dark"
	ThemeSystem = "system"
)

type Settings struct {
	Theme  string   `json:"theme"`            // light | dark | system
	Recent []string `json:"recent,omitempty"` // file kịch bản mở / lưu gần đây, mới nhất trước (storage.go)
}

func defaultSettings() Settings { return Settings{Theme: ThemeSystem} }

const settingsFile = "settings.json"

// ConfigDir: thư mục cấu hình của app, tạo nếu chưa có.
func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "autoclicker")
	return dir, os.MkdirAll(dir, 0o755)
}

// loadSettings đọc cấu hình trong dir; chưa có file hoặc giá trị lạ → mặc định (không báo lỗi để app vẫn mở được).
func loadSettings(dir string) Settings {
	s := defaultSettings()
	b, err := os.ReadFile(filepath.Join(dir, settingsFile))
	if err != nil {
		return s
	}
	if json.Unmarshal(b, &s) != nil || !validTheme(s.Theme) {
		s.Theme = ThemeSystem
	}
	return s
}

func saveSettings(dir string, s Settings) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return WriteFileAtomic(filepath.Join(dir, settingsFile), b)
}

// WriteFileAtomic ghi file tạm cạnh file đích rồi đổi tên – tắt app giữa chừng không làm hỏng file cũ.
func WriteFileAtomic(path string, b []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func validTheme(t string) bool { return t == ThemeLight || t == ThemeDark || t == ThemeSystem }

// CurrentSettings trả cấu hình đã lưu (mặc định nếu chưa có / không đọc được thư mục cấu hình).
func CurrentSettings() Settings {
	dir, err := ConfigDir()
	if err != nil {
		return defaultSettings()
	}
	return loadSettings(dir)
}

// SetTheme lưu theme: light | dark | system.
func SetTheme(theme string) error {
	if !validTheme(theme) {
		return fmt.Errorf("theme không hợp lệ: %q", theme)
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	dir, err := ConfigDir()
	if err != nil {
		return fmt.Errorf("không mở được thư mục cấu hình: %w", err)
	}
	s := loadSettings(dir)
	s.Theme = theme
	if err := saveSettings(dir, s); err != nil {
		return fmt.Errorf("không lưu được cấu hình: %w", err)
	}
	return nil
}
