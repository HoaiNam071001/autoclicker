package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

// Xác nhận cho phép xoá (node Tệp / thư mục xoá vào / dọn Thùng rác, bước 35 N4): theo kịch bản (đường dẫn file, hoặc tên
// khi chưa lưu) giữ các dấu vân tay thao tác xoá đã đồng ý (script.DeleteOps) – đổi thư mục / mẫu tên thì phải hỏi lại.

const deleteConsentFile = "delete-consent.json"

// maxConsents – số dấu vân tay giữ mỗi kịch bản (mới nhất trước).
const maxConsents = 10

func loadDeleteConsent(dir string) map[string][]string {
	all := map[string][]string{}
	if b, err := os.ReadFile(filepath.Join(dir, deleteConsentFile)); err == nil {
		_ = json.Unmarshal(b, &all) // hỏng thì coi như chưa xác nhận gì – hỏi lại, an toàn
	}
	return all
}

// DeleteConfirmed – người dùng đã đồng ý thao tác xoá có dấu vân tay fp của kịch bản id chưa.
func DeleteConfirmed(id, fp string) bool {
	dir, err := ConfigDir()
	if err != nil {
		return false
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	return slices.Contains(loadDeleteConsent(dir)[id], fp)
}

// ConfirmDelete ghi nhớ đồng ý của kịch bản id cho dấu vân tay fp.
func ConfirmDelete(id, fp string) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	all := loadDeleteConsent(dir)
	list := append([]string{fp}, slices.DeleteFunc(all[id], func(x string) bool { return x == fp })...)
	all[id] = list[:min(len(list), maxConsents)]
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return WriteFileAtomic(filepath.Join(dir, deleteConsentFile), b)
}
