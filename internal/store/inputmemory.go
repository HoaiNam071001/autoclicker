package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Giá trị lần trước + bộ giá trị (preset) của input – lưu trên máy theo từng kịch bản, không ghi vào file kịch bản.

const inputMemoryFile = "inputs.json"

type InputPreset struct {
	Name   string         `json:"name"`
	Values map[string]any `json:"values"`
}

type InputMemory struct {
	Last       map[string]any `json:"last,omitempty"`       // giá trị lần chạy trước
	LastPreset string         `json:"lastPreset,omitempty"` // bộ giá trị đã chọn lần trước ("" = tự điền)
	Presets    []InputPreset  `json:"presets,omitempty"`
	AutoRun    bool           `json:"autoRun,omitempty"` // "Lần sau chạy luôn với giá trị này" – bấm Chạy không hỏi
}

func loadInputMemory(dir string) map[string]InputMemory {
	all := map[string]InputMemory{}
	if b, err := os.ReadFile(filepath.Join(dir, inputMemoryFile)); err == nil {
		_ = json.Unmarshal(b, &all) // hỏng thì coi như chưa có
	}
	return all
}

// InputMemoryFor trả giá trị lần trước / preset của kịch bản id (đường dẫn file, hoặc tên khi chưa lưu).
func InputMemoryFor(id string) InputMemory {
	dir, err := ConfigDir()
	if err != nil {
		return InputMemory{}
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	return loadInputMemory(dir)[id]
}

// SaveInputMemory ghi giá trị lần trước / preset của kịch bản id.
func SaveInputMemory(id string, m InputMemory) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	all := loadInputMemory(dir)
	all[id] = m
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return WriteFileAtomic(filepath.Join(dir, inputMemoryFile), b)
}
