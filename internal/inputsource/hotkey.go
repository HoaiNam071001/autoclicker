package inputsource

import (
	"fmt"
	"slices"
	"strings"
)

// Phím tắt bật / tắt tiếng Việt của Unikey / EVKey / OpenKey (bước 43 S6b), viết "ctrl+shift", "alt+z", "ctrl+space".

var modifierAliases = map[string]string{
	"ctrl": "ctrl", "control": "ctrl", "shift": "shift", "alt": "alt", "option": "alt", "opt": "alt",
	"cmd": "cmd", "command": "cmd", "win": "cmd", "windows": "cmd", "super": "cmd",
}

// ParseHotkey tách phím tắt thành phím bổ trợ (ctrl, shift, alt, cmd) và phím chính (rỗng nếu chỉ có phím bổ trợ như
// Ctrl+Shift). Tên theo robotgo; không có phím nào, hoặc hai phím chính → lỗi.
func ParseHotkey(s string) (mods []string, key string, err error) {
	for _, part := range strings.Split(strings.ToLower(s), "+") {
		part = strings.TrimSpace(part)
		switch m, isMod := modifierAliases[part]; {
		case part == "":
			continue
		case isMod:
			if !slices.Contains(mods, m) {
				mods = append(mods, m)
			}
		case key != "":
			return nil, "", fmt.Errorf("phím tắt %q có hai phím chính (%s, %s) – chỉ một phím chính, còn lại là ctrl / shift / alt / cmd", s, key, part)
		default:
			key = part
		}
	}
	if len(mods) == 0 && key == "" {
		return nil, "", fmt.Errorf("phím tắt %q trống – vd ctrl+shift, alt+z", s)
	}
	return mods, key, nil
}
