package store

import (
	"fmt"
	"sync"
)

// Danh sách file kịch bản mở / lưu gần đây (Settings.Recent).

const maxRecent = 5

// settingsMu: đọc – sửa – ghi settings.json (theme, mở gần đây) từ nhiều lời gọi Wails cùng lúc.
var settingsMu sync.Mutex

// PushRecent đưa path lên đầu danh sách (không trùng, tối đa maxRecent).
func PushRecent(list []string, path string) []string {
	out := []string{path}
	for _, p := range list {
		if p != path && len(out) < maxRecent {
			out = append(out, p)
		}
	}
	return out
}

// RemoveRecent bỏ path khỏi danh sách (file không còn mở được).
func RemoveRecent(list []string, path string) []string {
	out := []string{}
	for _, p := range list {
		if p != path {
			out = append(out, p)
		}
	}
	return out
}

// EditRecent sửa danh sách mở gần đây trong settings.json; lỗi chỉ ghi log – không chặn lưu / mở.
func EditRecent(f func([]string) []string) {
	settingsMu.Lock()
	defer settingsMu.Unlock()
	dir, err := ConfigDir()
	if err != nil {
		return
	}
	s := loadSettings(dir)
	s.Recent = f(s.Recent)
	if err := saveSettings(dir, s); err != nil {
		fmt.Println("không lưu được danh sách mở gần đây:", err)
	}
}

// RecentScripts trả danh sách file mở / lưu gần đây, mới nhất trước.
func RecentScripts() []string {
	dir, err := ConfigDir()
	if err != nil {
		return []string{}
	}
	if r := loadSettings(dir).Recent; r != nil {
		return r
	}
	return []string{}
}
