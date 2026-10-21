// Package script – kịch bản như dữ liệu: mã hoá / giải mã file, kiểm tra lúc mở (checkScript) và trước khi chạy
// (ValidateScript), input, ô chữ nhận biến, neo điểm. Không đụng thiết bị hay đĩa.
package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"ghosthand/internal/model"
)

// EncodeScript ghi kịch bản ra JSON với version hiện tại.
func EncodeScript(sc model.Script) ([]byte, error) {
	sc.Version = model.ScriptVersion
	return json.MarshalIndent(sc, "", "  ")
}

// DecodeScript đọc file kịch bản: kiểm tra version, cấu trúc (checkScript) và bổ sung giá trị trống.
// Không chạy validateScript – bản đang soạn dở (vd bước gõ chữ chưa nhập) vẫn mở được, lúc Chạy mới kiểm tra.
func DecodeScript(b []byte) (model.Script, error) {
	var sc model.Script
	if err := json.Unmarshal(b, &sc); err != nil {
		return model.Script{}, fmt.Errorf("file hỏng hoặc không phải JSON: %w", err)
	}
	switch {
	case sc.Version == 0:
		return model.Script{}, errors.New("không phải file kịch bản (thiếu version)")
	case sc.Version > model.ScriptVersion:
		return model.Script{}, fmt.Errorf("file tạo bởi phiên bản app mới hơn (version %d, app đọc được tới %d) – hãy cập nhật app", sc.Version, model.ScriptVersion)
	case sc.Version < model.ScriptVersion:
		return model.Script{}, fmt.Errorf("kịch bản tạo bằng bản cũ (version %d, dạng danh sách bước), không mở được – bản này dùng dạng flow (version %d)", sc.Version, model.ScriptVersion)
	}
	if err := checkScript(&sc); err != nil {
		return model.Script{}, err
	}
	return sc, nil
}

// SafeFileName bỏ ký tự không dùng được trong tên file (Windows khắt khe nhất).
func SafeFileName(name string) string {
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) || r < 32 {
			return '-'
		}
		return r
	}, strings.TrimSpace(name))
	if name == "" {
		return "kich-ban"
	}
	return name
}
