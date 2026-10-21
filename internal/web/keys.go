package web

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/chromedp/cdproto/input"
)

// Node webKey (46 B3): bấm phím vào tab của Chrome của app qua CDP – không cần Chrome ở trước, không đụng bàn phím thật.
// Tên phím như node Phím tắt (robotgo: enter, tab, esc, f5, a…); sự kiện CDP đi qua đường nhập của trình duyệt nên
// trang nhận như phím thật (Tab chuyển ô, Enter gửi form).

type keyDef struct {
	key, code string
	vk        int64
	text      string // ký tự chèn khi không giữ ctrl / alt / cmd
}

var namedKeys = map[string]keyDef{
	"enter": {"Enter", "Enter", 13, "\r"}, "tab": {"Tab", "Tab", 9, ""}, "esc": {"Escape", "Escape", 27, ""},
	"escape": {"Escape", "Escape", 27, ""}, "backspace": {"Backspace", "Backspace", 8, ""}, "delete": {"Delete", "Delete", 46, ""},
	"space": {" ", "Space", 32, " "}, "up": {"ArrowUp", "ArrowUp", 38, ""}, "down": {"ArrowDown", "ArrowDown", 40, ""},
	"left": {"ArrowLeft", "ArrowLeft", 37, ""}, "right": {"ArrowRight", "ArrowRight", 39, ""}, "home": {"Home", "Home", 36, ""},
	"end": {"End", "End", 35, ""}, "pageup": {"PageUp", "PageUp", 33, ""}, "pagedown": {"PageDown", "PageDown", 34, ""},
}

// modifier → bit của CDP (Alt 1, Ctrl 2, Meta 4, Shift 8).
var modBits = map[string]int64{"alt": 1, "ctrl": 2, "control": 2, "cmd": 4, "command": 4, "meta": 4, "shift": 8}

// editCommands – lệnh soạn thảo cho tổ hợp sửa chữ (CDP không tự chạy lệnh của Ctrl / Cmd + phím trong ô nhập).
var editCommands = map[string]string{"a": "selectAll", "c": "copy", "x": "cut", "v": "paste", "z": "undo"}

// keyOf – phím theo tên robotgo: phím đặt tên, f1–f12, một chữ cái / chữ số / dấu.
func keyOf(name string) (keyDef, bool) {
	n := strings.ToLower(strings.TrimSpace(name))
	if d, ok := namedKeys[n]; ok {
		return d, true
	}
	var f int
	if _, err := fmt.Sscanf(n, "f%d", &f); err == nil && f >= 1 && f <= 12 && n == fmt.Sprintf("f%d", f) {
		return keyDef{key: strings.ToUpper(n), code: strings.ToUpper(n), vk: int64(111 + f)}, true
	}
	r := []rune(n)
	if len(r) != 1 || r[0] > unicode.MaxASCII || !unicode.IsPrint(r[0]) {
		return keyDef{}, false
	}
	c := r[0]
	switch {
	case c >= 'a' && c <= 'z':
		return keyDef{key: n, code: "Key" + strings.ToUpper(n), vk: int64(unicode.ToUpper(c)), text: n}, true
	case c >= '0' && c <= '9':
		return keyDef{key: n, code: "Digit" + n, vk: int64(c), text: n}, true
	}
	return keyDef{key: n, text: n}, true // dấu: chỉ cần key + chữ
}

// CheckKey – tên phím / phím kèm node webKey dùng được không (ValidateScript).
func CheckKey(name string, modifier bool) error {
	if modifier {
		if _, ok := modBits[strings.ToLower(strings.TrimSpace(name))]; !ok {
			return fmt.Errorf("phím kèm không hợp lệ: %q (cmd | ctrl | alt | shift)", name)
		}
		return nil
	}
	if _, ok := keyOf(name); !ok {
		return fmt.Errorf("phím %q không bấm được qua Chrome của app (enter, tab, esc, backspace, delete, space, mũi tên, home, end, pageup, pagedown, f1–f12, chữ / số)", name)
	}
	return nil
}

// keyEvents – sự kiện xuống / lên của phím name khi giữ các phím kèm mods.
func keyEvents(name string, mods []string) (down, up input.DispatchKeyEventParams) {
	d, _ := keyOf(name)
	var bits int64
	for _, m := range mods {
		bits |= modBits[strings.ToLower(strings.TrimSpace(m))]
	}
	down = input.DispatchKeyEventParams{Type: input.DispatchKeyEventTypeRawKeyDown, Modifiers: bits, Key: d.key, Code: d.code,
		WindowsVirtualKeyCode: d.vk}
	if bits&(1|2|4) == 0 && d.text != "" { // không giữ alt / ctrl / cmd → phím có chữ thì chèn chữ
		text := d.text
		if bits&8 != 0 {
			text = strings.ToUpper(text)
		}
		down.Type, down.Text, down.UnmodifiedText = input.DispatchKeyEventTypeKeyDown, text, d.text
	} else if cmd, ok := editCommands[d.key]; ok && bits&(2|4) != 0 {
		down.Commands = []string{cmd}
	}
	up = input.DispatchKeyEventParams{Type: input.DispatchKeyEventTypeKeyUp, Modifiers: bits, Key: d.key, Code: d.code,
		WindowsVirtualKeyCode: d.vk}
	return down, up
}

// Keys bấm phím vào tab (tab: id / URL / tiêu đề chứa; "" = tab node web đang dùng) – sequence: lần lượt từng phím (giữ
// phím kèm), không thì cả tổ hợp. Trả URL tab đã bấm.
func (c *Chrome) Keys(ctx context.Context, tab string, keys, mods []string, sequence bool) (string, error) {
	t, url, err := c.tabFor(ctx, tab)
	if err != nil {
		return "", err
	}
	if err := t.emulateFocus(ctx); err != nil { // trang tưởng đang được focus – ô nhập nhận phím dù Chrome ở sau
		return "", err
	}
	send := func(p input.DispatchKeyEventParams) error {
		return dispatch(ctx, t, "bấm phím "+p.Key, input.DispatchKeyEvent, p)
	}
	if sequence {
		for _, k := range keys {
			down, up := keyEvents(k, mods)
			if err := send(down); err != nil {
				return "", err
			}
			if err := send(up); err != nil {
				return "", err
			}
		}
		return url, nil
	}
	ups := make([]input.DispatchKeyEventParams, 0, len(keys))
	for _, k := range keys {
		down, up := keyEvents(k, mods)
		if err := send(down); err != nil {
			return "", err
		}
		ups = append(ups, up)
	}
	for i := len(ups) - 1; i >= 0; i-- {
		if err := send(ups[i]); err != nil {
			return "", err
		}
	}
	return url, nil
}
