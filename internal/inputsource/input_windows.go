package inputsource

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// Windows: layout bàn phím tính theo từng luồng cửa sổ → đọc layout của cửa sổ đang ở trên cùng, đổi bằng
// WM_INPUTLANGCHANGEREQUEST gửi cho nó (phần Windows viết nhưng chưa chạy thử trên máy thật).
var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procForegroundWindow    = user32.NewProc("GetForegroundWindow")
	procWindowThreadProcess = user32.NewProc("GetWindowThreadProcessId")
	procGetKeyboardLayout   = user32.NewProc("GetKeyboardLayout")
	procGetLayoutList       = user32.NewProc("GetKeyboardLayoutList")
	procPostMessage         = user32.NewProc("PostMessageW")
	procLCIDToLocaleName    = kernel32.NewProc("LCIDToLocaleName")
)

const wmInputLangChangeRequest = 0x0050

func hklSource(hkl uintptr) Source {
	lang := uint16(hkl & 0xFFFF)
	buf := make([]uint16, 85)
	var name string
	if n, _, _ := procLCIDToLocaleName.Call(uintptr(lang), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0); n > 0 {
		name = syscall.UTF16ToString(buf)
	}
	primary, _, _ := strings.Cut(name, "-")
	return Source{ID: fmt.Sprintf("%08X", uint32(hkl)), Name: name, Lang: strings.ToLower(primary)}
}

func foregroundLayout() (hwnd, hkl uintptr) {
	hwnd, _, _ = procForegroundWindow.Call()
	tid, _, _ := procWindowThreadProcess.Call(hwnd, 0)
	hkl, _, _ = procGetKeyboardLayout.Call(tid)
	return hwnd, hkl
}

func (System) Current() (Source, error) {
	hwnd, hkl := foregroundLayout()
	if hwnd == 0 || hkl == 0 {
		return Source{}, errors.New("không đọc được bộ gõ của cửa sổ đang ở trên cùng")
	}
	return hklSource(hkl), nil
}

func (System) List() ([]Source, error) {
	n, _, _ := procGetLayoutList.Call(0, 0)
	if n == 0 {
		return nil, errors.New("không đọc được danh sách bộ gõ")
	}
	hkls := make([]uintptr, n)
	procGetLayoutList.Call(n, uintptr(unsafe.Pointer(&hkls[0])))
	out := make([]Source, 0, n)
	for _, h := range hkls {
		out = append(out, hklSource(h))
	}
	return out, nil
}

func (System) Select(id string) error {
	v, err := strconv.ParseUint(id, 16, 32)
	if err != nil {
		return fmt.Errorf("mã bộ gõ %q không hợp lệ", id)
	}
	hwnd, _ := foregroundLayout()
	if hwnd == 0 {
		return errors.New("không có cửa sổ nào ở trên cùng để đổi bộ gõ")
	}
	if r, _, e := procPostMessage.Call(hwnd, wmInputLangChangeRequest, 0, uintptr(uint32(v))); r == 0 {
		return e
	}
	return nil
}
