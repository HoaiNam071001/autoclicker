package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

// Xác nhận thao tác nguy hiểm (xoá file, buộc thoát app, đóng cửa sổ – script.DangerOps): theo kịch bản (đường dẫn
// file, hoặc tên khi chưa lưu) giữ các dấu vân tay đã đồng ý – đổi đích (thư mục, mẫu tên, app…) thì phải hỏi lại.

const dangerConsentFile = "danger-consent.json"

// maxConsents – số dấu vân tay giữ mỗi kịch bản (mới nhất trước).
const maxConsents = 10

func loadDangerConsent(dir string) map[string][]string {
	all := map[string][]string{}
	if b, err := os.ReadFile(filepath.Join(dir, dangerConsentFile)); err == nil {
		_ = json.Unmarshal(b, &all) // hỏng thì coi như chưa xác nhận gì – hỏi lại, an toàn
	}
	return all
}

// DangerConfirmed – người dùng đã đồng ý thao tác nguy hiểm có dấu vân tay fp của kịch bản id chưa.
func DangerConfirmed(id, fp string) bool {
	dir, err := ConfigDir()
	if err != nil {
		return false
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	return slices.Contains(loadDangerConsent(dir)[id], fp)
}

// ConfirmDanger ghi nhớ đồng ý của kịch bản id cho dấu vân tay fp.
func ConfirmDanger(id, fp string) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	settingsMu.Lock()
	defer settingsMu.Unlock()
	all := loadDangerConsent(dir)
	list := append([]string{fp}, slices.DeleteFunc(all[id], func(x string) bool { return x == fp })...)
	all[id] = list[:min(len(list), maxConsents)]
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return WriteFileAtomic(filepath.Join(dir, dangerConsentFile), b)
}
