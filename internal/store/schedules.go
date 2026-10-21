package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"ghosthand/internal/model"
)

// Hẹn giờ (bước 22): schedules.json – mỗi file kịch bản một hẹn giờ, khoá theo đường dẫn.

const schedulesFile = "schedules.json"

func loadSchedules(dir string) map[string]model.Schedule {
	all := map[string]model.Schedule{}
	if b, err := os.ReadFile(filepath.Join(dir, schedulesFile)); err == nil {
		_ = json.Unmarshal(b, &all) // hỏng thì coi như chưa có
	}
	return all
}

// Schedules – mọi hẹn giờ, sắp theo đường dẫn.
func Schedules() []model.Schedule {
	dir, err := ConfigDir()
	if err != nil {
		return nil
	}
	settingsMu.Lock()
	all := loadSchedules(dir)
	settingsMu.Unlock()
	out := make([]model.Schedule, 0, len(all))
	for _, s := range all {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// EditSchedules sửa danh sách hẹn giờ dưới khoá (đặt / bỏ).
func EditSchedules(f func(m map[string]model.Schedule)) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	all := loadSchedules(dir)
	f(all)
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return WriteFileAtomic(filepath.Join(dir, schedulesFile), b)
}
