package store

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// API local (bước 22): app đang mở ghi cổng + token của kênh cho CLI vào api.json (chỉ người dùng này đọc được);
// CLI đọc file để gửi lệnh sang app. Mỗi lần mở app một token mới, đóng app thì xoá file.

const apiInfoFile = "api.json"

type APIInfo struct {
	Port  int    `json:"port"`
	Token string `json:"token"`
	PID   int    `json:"pid"`
}

// NewToken – chuỗi ngẫu nhiên 32 byte (hex).
func NewToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("không tạo được token: %v", err))
	}
	return hex.EncodeToString(b)
}

func apiInfoPath() (string, error) {
	dir, err := ConfigDir()
	return filepath.Join(dir, apiInfoFile), err
}

func WriteAPIInfo(i APIInfo) error {
	p, err := apiInfoPath()
	if err != nil {
		return err
	}
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

// ReadAPIInfo – thông tin app đang mở; chưa mở (không có file) → lỗi os.ErrNotExist.
func ReadAPIInfo() (APIInfo, error) {
	var i APIInfo
	p, err := apiInfoPath()
	if err != nil {
		return i, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(b, &i)
}

// RemoveAPIInfo xoá api.json nếu còn là của tiến trình pid (app mở sau không bị xoá nhầm).
func RemoveAPIInfo(pid int) {
	if i, err := ReadAPIInfo(); err == nil && i.PID == pid {
		if p, err := apiInfoPath(); err == nil {
			os.Remove(p)
		}
	}
}
