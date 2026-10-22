// Package shell chạy lệnh / script của node Chạy lệnh (bước 43 S1): script đi vào stdin hoặc đối số của shell, biến của
// kịch bản chỉ đi qua biến môi trường – không bao giờ ghép vào chuỗi lệnh (tránh chèn lệnh).
package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"ghosthand/internal/model"
)

// ErrUnsupported – ngôn ngữ script không có trên hệ điều hành này (AppleScript · JXA chỉ macOS).
var ErrUnsupported = errors.New("ngôn ngữ script không chạy được trên hệ điều hành này")

// OutputMax – mỗi luồng stdout / stderr giữ tối đa chừng này byte (phần sau bị bỏ, lệnh vẫn chạy tiếp).
const OutputMax = 64 << 10

// Spec – một lần chạy. Env là "KEY=value" thêm vào môi trường của app.
type Spec struct {
	Lang   string // model.ShellLangs
	Script string
	Dir    string // "" = thư mục hiện tại của app
	Env    []string
}

// Result – kết quả lệnh đã chạy xong (ExitCode ≠ 0 vẫn là kết quả, không phải lỗi).
type Result struct {
	Stdout, Stderr string
	ExitCode       int
}

// Runner – chạy một lệnh tới khi xong hoặc ctx huỷ (huỷ = giết cả cây tiến trình, trả ctx.Err()). Test thay bản giả.
type Runner interface {
	Run(ctx context.Context, sp Spec) (Result, error)
}

// System – Runner thật bằng os/exec.
type System struct{}

func (System) Run(ctx context.Context, sp Spec) (Result, error) {
	cmd, stdin, err := command(sp.Lang, sp.Script)
	if err != nil {
		return Result{}, err
	}
	cmd.Dir = sp.Dir
	cmd.Env = append(os.Environ(), sp.Env...)
	if stdin != "" {
		cmd.Stdin = bytes.NewBufferString(stdin)
	}
	var out, errb capBuf
	cmd.Stdout, cmd.Stderr = &out, &errb
	own(cmd)
	cmd.WaitDelay = 2 * time.Second
	if err := cmd.Start(); err != nil {
		return Result{}, fmt.Errorf("không chạy được lệnh: %w", err)
	}
	// ctx huỷ: giết cả cây (con cháu giữ pipe thì Wait mới không treo)
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = killTree(cmd)
		case <-done:
		}
	}()
	werr := cmd.Wait()
	if ctx.Err() != nil {
		return Result{}, ctx.Err()
	}
	res := Result{Stdout: out.String(), Stderr: errb.String(), ExitCode: cmd.ProcessState.ExitCode()}
	var exit *exec.ExitError
	if werr != nil && !errors.As(werr, &exit) && !errors.Is(werr, exec.ErrWaitDelay) {
		return res, werr
	}
	return res, nil
}

// command – tiến trình chạy script lang và nội dung gửi vào stdin ("" = không có).
func command(lang, script string) (*exec.Cmd, string, error) {
	switch lang {
	case model.ShellPowerShell:
		if runtime.GOOS == "windows" {
			return exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", "-"), script, nil
		}
		if p, err := exec.LookPath("pwsh"); err == nil {
			return exec.Command(p, "-NoProfile", "-NonInteractive", "-Command", "-"), script, nil
		}
		return nil, "", ErrUnsupported
	case model.ShellAppleScript, model.ShellJXA:
		if runtime.GOOS != "darwin" {
			return nil, "", ErrUnsupported
		}
		if lang == model.ShellJXA {
			return exec.Command("osascript", "-l", "JavaScript", "-"), script, nil
		}
		return exec.Command("osascript", "-"), script, nil
	}
	return shellCommand(script), "", nil
}

// capBuf – giữ OutputMax byte đầu, nuốt phần còn lại (không trả lỗi để lệnh không chết vì pipe đầy).
type capBuf struct{ b bytes.Buffer }

func (c *capBuf) Write(p []byte) (int, error) {
	if room := OutputMax - c.b.Len(); room > 0 {
		c.b.Write(p[:min(len(p), room)])
	}
	return len(p), nil
}

func (c *capBuf) String() string { return c.b.String() }
