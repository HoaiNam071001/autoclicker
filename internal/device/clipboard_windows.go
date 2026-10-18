package device

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// hideWindow – không bật cửa sổ PowerShell đen lên khi chạy.
func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
}

// ps chạy script PowerShell (-STA: clipboard WinForms cần luồng STA); dữ liệu qua stdin / biến môi trường, không ghép vào script.
func ps(script string, env ...string) *exec.Cmd {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-STA", "-Command",
		"[Console]::InputEncoding=[Text.Encoding]::UTF8;[Console]::OutputEncoding=[Text.Encoding]::UTF8;"+script)
	cmd.Env = append(os.Environ(), env...)
	return cmd
}

func (SystemClipboard) ReadText() (string, error) {
	out, err := run(ps("Get-Clipboard -Raw"), "")
	if err != nil {
		return "", fmt.Errorf("không đọc được clipboard: %w", err)
	}
	return strings.TrimSuffix(out, "\r\n"), nil
}

func (SystemClipboard) WriteText(text string) error {
	if _, err := run(ps("Set-Clipboard -Value ([Console]::In.ReadToEnd())"), text); err != nil {
		return fmt.Errorf("không ghi được clipboard: %w", err)
	}
	return nil
}

func (SystemClipboard) WriteImage(pngPath string) error {
	script := "Add-Type -AssemblyName System.Windows.Forms,System.Drawing;" +
		"$img=[System.Drawing.Image]::FromFile($env:AC_IMAGE);[System.Windows.Forms.Clipboard]::SetImage($img);$img.Dispose()"
	if _, err := run(ps(script, "AC_IMAGE="+pngPath), ""); err != nil {
		return fmt.Errorf("không chép được ảnh vào clipboard: %w", err)
	}
	return nil
}
