//go:build !windows

package power

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// KeepAwake giữ máy (và màn hình) thức tới khi gọi stop. macOS: `caffeinate` bám theo tiến trình của app, nên app thoát
// đột ngột cũng tự nhả. Hệ khác: không làm gì.
func KeepAwake(ctx context.Context) (stop func(), err error) {
	if runtime.GOOS != "darwin" {
		return func() {}, nil
	}
	cmd := exec.CommandContext(ctx, "caffeinate", "-dims", "-w", strconv.Itoa(os.Getpid()))
	if err := cmd.Start(); err != nil {
		return func() {}, err
	}
	go cmd.Wait()
	return func() { _ = cmd.Process.Kill() }, nil
}
