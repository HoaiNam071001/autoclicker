package element

import (
	"strconv"
	"strings"
)

// uiaRoles – ControlType của Windows UI Automation (UIA_*ControlTypeId) → vai trò kiểu web, cùng tên với axRoles.
var uiaRoles = map[int]string{
	50000: "button", 50002: "checkbox", 50003: "combobox", 50004: "textbox", 50005: "link", 50006: "img",
	50007: "listitem", 50008: "list", 50009: "menu", 50010: "menubar", 50011: "menuitem", 50012: "progressbar",
	50013: "radio", 50014: "scrollbar", 50015: "slider", 50016: "spinbutton", 50017: "status", 50018: "tablist",
	50019: "tab", 50020: RoleText, 50021: "toolbar", 50022: "tooltip", 50023: "tree", 50024: "treeitem",
	50025: "custom", 50026: "group", 50028: "grid", 50029: "row", 50030: "document", 50031: "button",
	50032: "window", 50033: "pane", 50035: "columnheader", 50036: "table", 50037: "titlebar", 50038: "separator",
}

// RoleFromUIA – vai trò chuẩn hoá từ phần tử UIA: Chrome / Edge báo sẵn vai trò web (AriaRole – heading, searchbox…)
// thì dùng luôn, không thì theo ControlType; ngoài bảng → "uia<số>" (để gỡ rối). Khối chữ Chrome báo AriaRole
// "description" – bỏ, theo ControlType (Text → text).
func RoleFromUIA(controlType int, ariaRole string) string {
	if r := strings.ToLower(strings.TrimSpace(ariaRole)); r != "" && r != "description" {
		return r
	}
	if r, ok := uiaRoles[controlType]; ok {
		return r
	}
	return "uia" + strconv.Itoa(controlType)
}
