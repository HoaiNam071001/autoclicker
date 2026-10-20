package cli

import (
	"os"
	"os/exec"
	"syscall"
)

// exe Windows build dạng GUI (không có console) – gắn vào console của terminal đang gọi để in được chữ.
// Chạy từ Explorer (không có console cha) thì bỏ qua.
func AttachConsole() {
	const attachParent = ^uintptr(0) // ATTACH_PARENT_PROCESS = (DWORD)-1
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("AttachConsole")
	if ok, _, _ := proc.Call(attachParent); ok == 0 {
		return
	}
	if f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0); err == nil {
		os.Stdout, os.Stderr = f, f
	}
	if f, err := os.OpenFile("CONIN$", os.O_RDONLY, 0); err == nil {
		os.Stdin = f
	}
}

func detach(cmd *exec.Cmd) {
	const detachedProcess, newGroup = 0x00000008, 0x00000200
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: detachedProcess | newGroup}
}
