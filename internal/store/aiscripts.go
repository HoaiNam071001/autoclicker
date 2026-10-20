package store

import (
	"os"
	"path/filepath"
	"strconv"
)

// ScriptsDir – thư mục mặc định cho kịch bản AI tạo (bước 31): Documents/Auto Clicker, tạo nếu chưa có.
func ScriptsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Documents", "Auto Clicker")
	return dir, os.MkdirAll(dir, 0o755)
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
