package device

import (
	"fmt"
	"os"
	"os/exec"
)

// pbcopy / pbpaste đọc ghi theo bảng mã của locale – đặt UTF-8 để giữ tiếng Việt.
func utf8Env(cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = append(os.Environ(), "LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8")
	return cmd
}

func hideWindow(*exec.Cmd) {}

func (SystemClipboard) ReadText() (string, error) {
	out, err := run(utf8Env(exec.Command("pbpaste")), "")
	if err != nil {
		return "", fmt.Errorf("không đọc được clipboard: %w", err)
	}
	return out, nil
}

func (SystemClipboard) WriteText(text string) error {
	if _, err := run(utf8Env(exec.Command("pbcopy")), text); err != nil {
		return fmt.Errorf("không ghi được clipboard: %w", err)
	}
	return nil
}

// WriteImage – đường dẫn qua argv (không ghép vào chuỗi AppleScript, khỏi lo dấu nháy trong tên file).
func (SystemClipboard) WriteImage(pngPath string) error {
	cmd := exec.Command("osascript", "-e", "on run argv", "-e",
		"set the clipboard to (read (POSIX file (item 1 of argv)) as «class PNGf»)", "-e", "end run", pngPath)
	if _, err := run(cmd, ""); err != nil {
		return fmt.Errorf("không chép được ảnh vào clipboard: %w", err)
	}
	return nil
}
