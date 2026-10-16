package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Cấu hình app (không thuộc kịch bản), lưu ở os.UserConfigDir()/autoclicker/settings.json.
// Bước 21 dùng chung thư mục này cho bản nháp và danh sách mở gần đây.

const (
	ThemeLight  = "light"
	ThemeDark   = "dark"
	ThemeSystem = "system"
)

type Settings struct {
	Theme string `json:"theme"` // light | dark | system
}

func defaultSettings() Settings { return Settings{Theme: ThemeSystem} }

const settingsFile = "settings.json"

// configDir: thư mục cấu hình của app, tạo nếu chưa có.
func configDir() (string, error) {
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
	// ghi file tạm rồi đổi tên – tắt app giữa chừng không làm hỏng file cũ
	tmp := filepath.Join(dir, settingsFile+".tmp")
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, settingsFile))
}

func validTheme(t string) bool { return t == ThemeLight || t == ThemeDark || t == ThemeSystem }

// GetSettings trả cấu hình đã lưu (mặc định nếu chưa có / không đọc được thư mục cấu hình).
func (a *App) GetSettings() Settings {
	dir, err := configDir()
	if err != nil {
		return defaultSettings()
	}
	return loadSettings(dir)
}

// SetTheme lưu theme: light | dark | system.
func (a *App) SetTheme(theme string) error {
	if !validTheme(theme) {
		return fmt.Errorf("theme không hợp lệ: %q", theme)
	}
	dir, err := configDir()
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
