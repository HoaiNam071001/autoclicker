// Package toast hiện thông báo của hệ điều hành (Notification Center macOS / toast Windows) cho node Thông báo hệ
// thống (bước 43 S9). Chữ đi qua đối số / biến môi trường, không ghép vào chuỗi lệnh. Chưa có nút bấm (S9b).
package toast

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var ErrUnsupported = errors.New("thông báo hệ thống chỉ có trên macOS và Windows")

// Backend – hiện một thông báo. Test thay bản giả.
type Backend interface {
	Show(ctx context.Context, title, text string) error
}

// System – Backend thật.
type System struct{}

// macOS: osascript nhận chữ qua argv; AppleScript dùng item của argv nên không có chỗ chèn lệnh.
const appleScript = `on run argv
display notification (item 2 of argv) with title (item 1 of argv)
end run`

// Windows: toast hai dòng, mượn AppUserModelID của PowerShell (app chưa đăng ký ID riêng thì Windows không hiện toast).
const powershellScript = `
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
$x = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$t = $x.GetElementsByTagName('text')
$t.Item(0).AppendChild($x.CreateTextNode($env:GH_TOAST_TITLE)) | Out-Null
$t.Item(1).AppendChild($x.CreateTextNode($env:GH_TOAST_TEXT)) | Out-Null
$id = '{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\WindowsPowerShell\v1.0\powershell.exe'
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($id).Show([Windows.UI.Notifications.ToastNotification]::new($x))
`

func (System) Show(ctx context.Context, title, text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "osascript", "-e", appleScript, title, text)
	case "windows":
		cmd = exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", powershellScript)
		cmd.Env = append(os.Environ(), "GH_TOAST_TITLE="+title, "GH_TOAST_TEXT="+text)
	default:
		return ErrUnsupported
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if msg := strings.TrimSpace(string(out)); msg != "" {
			return fmt.Errorf("không hiện được thông báo: %s", msg)
		}
		return fmt.Errorf("không hiện được thông báo: %w", err)
	}
	return nil
}
