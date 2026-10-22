// Package power – nguồn / phiên làm việc của máy cho node Nguồn (bước 43 S3): khoá màn, tắt màn hình, ngủ, đăng xuất,
// tắt, khởi động lại; và giữ máy thức trong lúc kịch bản chạy (KeepAwake). Lệnh chạy bằng công cụ có sẵn của hệ điều hành,
// không nhận chữ từ người dùng nên không có chỗ chèn lệnh.
package power

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Hành động của node Nguồn (Step.PowerOp). Khớp model.PowerOps.
const (
	Lock         = "lock"
	DisplaySleep = "displaySleep"
	Sleep        = "sleep"
	Logout       = "logout"
	Shutdown     = "shutdown"
	Restart      = "restart"
)

var ErrUnsupported = errors.New("nguồn máy chỉ điều khiển được trên macOS và Windows")

// Backend – làm một hành động nguồn. Test thay bản giả.
type Backend interface {
	Do(ctx context.Context, action string) error
}

// System – Backend thật.
type System struct{}

func (System) Do(ctx context.Context, action string) error {
	argv, err := command(runtime.GOOS, action)
	if err != nil {
		return err
	}
	out, err := exec.CommandContext(ctx, argv[0], argv[1:]...).CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if msg := strings.TrimSpace(string(out)); msg != "" {
			return fmt.Errorf("không thực hiện được %s: %s", action, msg)
		}
		return fmt.Errorf("không thực hiện được %s: %w", action, err)
	}
	return nil
}

const winDisplayOff = `(Add-Type '[DllImport("user32.dll")]public static extern int SendMessage(int h,int m,int w,int l);' -Name W -PassThru)::SendMessage(-1,0x0112,0xF170,2)`

// command – dòng lệnh làm action trên hệ goos.
func command(goos, action string) ([]string, error) {
	osa := func(script string) []string { return []string{"osascript", "-e", script} }
	switch goos {
	case "darwin":
		switch action {
		case Lock: // Ctrl+Cmd+Q – cần quyền Accessibility cho app
			return osa(`tell application "System Events" to keystroke "q" using {control down, command down}`), nil
		case DisplaySleep:
			return []string{"pmset", "displaysleepnow"}, nil
		case Sleep:
			return []string{"pmset", "sleepnow"}, nil
		case Logout:
			return osa(`tell application "loginwindow" to «event aevtrlgo»`), nil
		case Shutdown:
			return osa(`tell application "System Events" to shut down`), nil
		case Restart:
			return osa(`tell application "System Events" to restart`), nil
		}
	case "windows":
		switch action {
		case Lock:
			return []string{"rundll32.exe", "user32.dll,LockWorkStation"}, nil
		case DisplaySleep:
			return []string{"powershell", "-NoProfile", "-NonInteractive", "-Command", winDisplayOff}, nil
		case Sleep:
			return []string{"rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0"}, nil
		case Logout:
			return []string{"shutdown", "/l"}, nil
		case Shutdown:
			return []string{"shutdown", "/s", "/t", "0"}, nil
		case Restart:
			return []string{"shutdown", "/r", "/t", "0"}, nil
		}
	default:
		return nil, ErrUnsupported
	}
	return nil, fmt.Errorf("hành động nguồn không hợp lệ: %q", action)
}
