//go:build windows

// Package winexclude: cho cửa sổ của app không lọt vào ảnh chụp màn hình (thanh mini, overlay) – Windows 10 2004+.
package winexclude

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procEnumWindows         = user32.NewProc("EnumWindows")
	procGetWindowThreadPid  = user32.NewProc("GetWindowThreadProcessId")
	procSetWindowDisplayAff = user32.NewProc("SetWindowDisplayAffinity")
	wdaExcludeFromCapture   = uintptr(0x11)
	callback                uintptr // syscall.NewCallback có giới hạn số lần tạo – tạo một lần
	matched                 int
)

// ExcludeOwnWindows đặt WDA_EXCLUDEFROMCAPTURE cho mọi cửa sổ cấp cao của tiến trình này. Trả số cửa sổ đặt được.
func ExcludeOwnWindows() int {
	if procSetWindowDisplayAff.Find() != nil {
		return 0
	}
	pid := uint32(os.Getpid())
	matched = 0
	if callback == 0 {
		callback = syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
			var p uint32
			procGetWindowThreadPid.Call(hwnd, uintptr(unsafe.Pointer(&p)))
			if p == pid {
				if r, _, _ := procSetWindowDisplayAff.Call(hwnd, wdaExcludeFromCapture); r != 0 {
					matched++
				}
			}
			return 1 // tiếp tục duyệt
		})
	}
	procEnumWindows.Call(callback, 0)
	return matched
}
