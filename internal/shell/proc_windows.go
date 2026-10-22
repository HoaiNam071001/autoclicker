//go:build windows

package shell

import (
	"os/exec"
	"strconv"
	"syscall"
)

// shellCommand – cmd /C: giữ nguyên chuỗi lệnh (CmdLine) thay vì để Go bọc ngoặc kép theo kiểu CommandLineToArgv.
func shellCommand(script string) *exec.Cmd {
	cmd := exec.Command("cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "cmd.exe /C " + script}
	return cmd
}

// own – không bật cửa sổ console; nhóm riêng để Ctrl+C của app không lan vào.
func own(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
	cmd.SysProcAttr.CreationFlags |= 0x08000000 | 0x00000200 // CREATE_NO_WINDOW | CREATE_NEW_PROCESS_GROUP
}

func killTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	k := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	k.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return k.Run()
}
