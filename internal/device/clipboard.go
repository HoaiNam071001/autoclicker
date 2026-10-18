package device

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Clipboard – đọc / ghi clipboard của hệ điều hành (node Chụp màn hình chép ảnh – bước 35 N1, node Clipboard – N2).
// Qua interface để test engine bằng bản giả (devicetest.FakeClipboard).
type Clipboard interface {
	ReadText() (string, error)
	WriteText(text string) error
	WriteImage(pngPath string) error // ảnh PNG từ file
}

// SystemClipboard – clipboard thật: macOS pbcopy / pbpaste / osascript, Windows PowerShell (clipboard_<os>.go).
type SystemClipboard struct{}

// run chạy lệnh, stdin = in; lỗi kèm stderr để người dùng biết vì sao.
func run(cmd *exec.Cmd, in string) (string, error) {
	var out, errb bytes.Buffer
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(in), &out, &errb
	hideWindow(cmd)
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(errb.String()); msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}
	return out.String(), nil
}
