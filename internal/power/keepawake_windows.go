package power

import (
	"context"
	"runtime"
	"sync"
	"syscall"
)

var procSetThreadExecutionState = syscall.NewLazyDLL("kernel32.dll").NewProc("SetThreadExecutionState")

const (
	esContinuous      = 0x80000000
	esSystemRequired  = 0x00000001
	esDisplayRequired = 0x00000002
)

// KeepAwake giữ máy thức: SetThreadExecutionState có hiệu lực theo luồng nên giữ một luồng OS riêng cho tới khi stop
// (viết nhưng chưa chạy thử trên máy thật).
func KeepAwake(ctx context.Context) (stop func(), err error) {
	done := make(chan struct{})
	ready := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		procSetThreadExecutionState.Call(uintptr(esContinuous | esSystemRequired | esDisplayRequired))
		close(ready)
		select {
		case <-done:
		case <-ctx.Done():
		}
		procSetThreadExecutionState.Call(uintptr(esContinuous))
	}()
	<-ready
	var once sync.Once
	return func() { once.Do(func() { close(done) }) }, nil
}
