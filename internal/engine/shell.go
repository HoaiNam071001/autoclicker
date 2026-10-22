package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/shell"
)

// Node Chạy lệnh (bước 43 S1): chạy Script bằng shell / PowerShell / AppleScript · JXA. Chỉ chạy khi người dùng đã
// đồng ý đúng nội dung lệnh (Options.AllowDanger – script.DangerShell). Script nguyên văn; biến chỉ đi qua ShellEnv
// (biến môi trường) nên giá trị biến không thể biến thành lệnh. Hết giờ / Dừng → giết cả cây tiến trình.

// exitErr – lệnh thoát với mã ≠ 0 (StepError code exit_code; output stdout / stderr / exitCode vẫn có).
type exitErr struct{ msg string }

func (e *exitErr) Error() string { return e.msg }

func (e *Engine) shellStep(ctx context.Context, s model.Step) (map[string]any, error) {
	if !e.allowDanger {
		return nil, dangerErr("chạy lệnh")
	}
	sp, err := shellSpec(s)
	if err != nil {
		return nil, err
	}
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	timeout := s.TimeoutMs
	if timeout <= 0 {
		timeout = model.ShellTimeoutMs
	}
	cctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()
	res, err := e.sh.Run(cctx, sp)
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case errors.Is(err, context.DeadlineExceeded):
		return nil, &timeoutErr{msg: fmt.Sprintf("lệnh chạy quá %s – đã dừng", fmtSeconds(timeout)), detail: map[string]any{"timeoutMs": timeout}}
	case errors.Is(err, shell.ErrUnsupported):
		return nil, &invalidErr{fmt.Sprintf("ngôn ngữ %q không chạy được trên hệ điều hành này", s.ShellLang)}
	case err != nil:
		return nil, err
	}
	out := map[string]any{"stdout": trimEOL(res.Stdout), "stderr": trimEOL(res.Stderr), "exitCode": res.ExitCode}
	if res.ExitCode != 0 {
		msg := fmt.Sprintf("lệnh thoát với mã %d", res.ExitCode)
		if line, _, _ := strings.Cut(strings.TrimSpace(res.Stderr), "\n"); line != "" {
			msg += ": " + line
		}
		return out, &exitErr{msg}
	}
	return out, nil
}

// shellSpec – thư mục chạy (có thật) + biến môi trường từ ShellEnv.
func shellSpec(s model.Step) (shell.Spec, error) {
	sp := shell.Spec{Lang: s.ShellLang, Script: s.Script}
	if strings.TrimSpace(s.Script) == "" {
		return sp, &invalidErr{"chưa có nội dung lệnh"}
	}
	if strings.TrimSpace(s.FilePath) != "" {
		dir, err := userPath(s.FilePath)
		if err != nil {
			return sp, err
		}
		if st, err := os.Stat(dir); err != nil || !st.IsDir() {
			return sp, &invalidErr{fmt.Sprintf("thư mục chạy lệnh không có: %s", dir)}
		}
		sp.Dir = dir
	}
	for _, kv := range s.ShellEnv {
		if k := strings.TrimSpace(kv.Key); k != "" {
			sp.Env = append(sp.Env, k+"="+kv.Value)
		}
	}
	return sp, nil
}

func trimEOL(s string) string { return strings.TrimRight(s, "\r\n") }
