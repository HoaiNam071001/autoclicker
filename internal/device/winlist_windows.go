package device

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

// Danh sách cửa sổ cấp cao đang hiện (Win32) – chọn app cho phần tử màn hình (37 E5); bước 26 W0 dùng lại.

var (
	procEnumWindows          = user32.NewProc("EnumWindows")
	procIsWindowVisible      = user32.NewProc("IsWindowVisible")
	procIsIconic             = user32.NewProc("IsIconic")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procGetWindowThreadPID   = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowLongPtrW    = user32.NewProc("GetWindowLongPtrW")
	procDwmGetWindowAttr     = syscall.NewLazyDLL("dwmapi.dll").NewProc("DwmGetWindowAttribute")
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess          = kernel32.NewProc("OpenProcess")
	procQueryFullProcessName = kernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandle          = kernel32.NewProc("CloseHandle")

	enumMu   sync.Mutex
	enumCB   uintptr // syscall.NewCallback có giới hạn số lần tạo – tạo một lần
	enumList []uintptr
)

// topWindow – cửa sổ cấp cao: hwnd, tiến trình, tiêu đề, đang thu nhỏ.
type topWindow struct {
	hwnd      uintptr
	pid       uint32
	title     string
	minimized bool
}

// topWindows – cửa sổ cấp cao theo thứ tự z (trên cùng trước): đang hiện, không phải cửa sổ công cụ, không bị DWM
// giấu (app UWP chạy nền, desktop ảo khác), có tiêu đề; kể cả cửa sổ thu nhỏ (minimized = true).
func topWindows() []topWindow {
	enumMu.Lock()
	defer enumMu.Unlock()
	if enumCB == 0 {
		enumCB = syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
			enumList = append(enumList, hwnd)
			return 1
		})
	}
	enumList = enumList[:0]
	procEnumWindows.Call(enumCB, 0)
	var out []topWindow
	for _, h := range enumList {
		if v, _, _ := procIsWindowVisible.Call(h); v == 0 {
			continue
		}
		const gwlExStyle, wsExToolWindow = ^uintptr(19), 0x80 // GWL_EXSTYLE = -20
		if ex, _, _ := procGetWindowLongPtrW.Call(h, gwlExStyle); ex&wsExToolWindow != 0 {
			continue
		}
		var cloaked uint32
		if r, _, _ := procDwmGetWindowAttr.Call(h, 14, uintptr(unsafe.Pointer(&cloaked)), 4); r == 0 && cloaked != 0 {
			continue // DWMWA_CLOAKED
		}
		title := windowTitle(h)
		if title == "" {
			continue
		}
		var pid uint32
		procGetWindowThreadPID.Call(h, uintptr(unsafe.Pointer(&pid)))
		icon, _, _ := procIsIconic.Call(h)
		out = append(out, topWindow{h, pid, title, icon != 0})
	}
	return out
}

func windowTitle(h uintptr) string {
	buf := make([]uint16, 512)
	n, _, _ := procGetWindowTextW.Call(h, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf[:n])
}

// processName – tên tệp chạy của tiến trình, bỏ .exe (chrome, msedge, notepad); không mở được → "".
func processName(pid uint32) string {
	const queryLimited = 0x1000 // PROCESS_QUERY_LIMITED_INFORMATION – đọc được cả tiến trình của người dùng khác quyền
	h, _, _ := procOpenProcess.Call(queryLimited, 0, uintptr(pid))
	if h == 0 {
		return ""
	}
	defer procCloseHandle.Call(h)
	buf := make([]uint16, 1024)
	n := uint32(len(buf))
	if r, _, _ := procQueryFullProcessName.Call(h, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n))); r == 0 {
		return ""
	}
	base := filepath.Base(syscall.UTF16ToString(buf[:n]))
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// windowMatches – app khớp tên tệp chạy (chrome / chrome.exe), cả tiêu đề, hoặc đuôi tiêu đề “… - Google Chrome”
// (kịch bản soạn trên macOS ghi tên app hiển thị). Không phân biệt hoa thường.
func windowMatches(app, exe, title string) bool {
	a, t := strings.ToLower(strings.TrimSpace(app)), strings.ToLower(title)
	e := strings.ToLower(exe)
	return a == e || a == e+".exe" || a == t || strings.HasSuffix(t, " - "+a) || strings.HasSuffix(t, " — "+a)
}

// appWindows – các cửa sổ đang hiện (không thu nhỏ) của app; "" = app có cửa sổ trên cùng, không tính app này.
// Trả kèm tên tệp chạy để báo lỗi.
func appWindows(app string) ([]uintptr, string, error) {
	self := uint32(os.Getpid())
	var pid uint32
	var minimized bool
	all := topWindows()
	for _, w := range all {
		if w.pid == self {
			continue
		}
		if app == "" && !w.minimized || app != "" && windowMatches(app, processName(w.pid), w.title) {
			if !w.minimized {
				pid = w.pid
				break
			}
			minimized = true
		}
	}
	if pid == 0 {
		switch {
		case app == "":
			return nil, "", ErrNoWindow
		case minimized:
			return nil, app, ErrNoWindow
		}
		return nil, app, ErrAppNotFound
	}
	var hwnds []uintptr
	for _, w := range all {
		if w.pid == pid && !w.minimized {
			hwnds = append(hwnds, w.hwnd)
		}
	}
	return hwnds, processName(pid), nil
}
