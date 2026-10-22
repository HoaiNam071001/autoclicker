//go:build !windows

package shell

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRunOutputAndExitCode(t *testing.T) {
	res, err := System{}.Run(t.Context(), Spec{Script: `echo "chao $TEN"; echo loi >&2; exit 3`, Env: []string{"TEN=Nam"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 3 || strings.TrimSpace(res.Stdout) != "chao Nam" || strings.TrimSpace(res.Stderr) != "loi" {
		t.Errorf("kết quả sai: %+v", res)
	}
}

func TestRunEnvIsNotInterpreted(t *testing.T) {
	res, err := System{}.Run(t.Context(), Spec{Script: `printf %s "$X"`, Env: []string{"X=$(echo hacked); `echo hacked`"}})
	if err != nil || res.Stdout != "$(echo hacked); `echo hacked`" {
		t.Errorf("giá trị biến bị hiểu như lệnh: %q %v", res.Stdout, err)
	}
}

func TestRunDir(t *testing.T) {
	dir := t.TempDir()
	res, err := System{}.Run(t.Context(), Spec{Script: "pwd -P", Dir: dir})
	if err != nil || !strings.HasSuffix(strings.TrimSpace(res.Stdout), strings.TrimPrefix(dir, "/private")) {
		t.Errorf("thư mục chạy sai: %q %v", res.Stdout, err)
	}
}

func TestRunKillsWholeTreeOnCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	// sleep con giữ pipe: không giết cả nhóm thì Wait treo tới 30 s
	_, err := System{}.Run(ctx, Spec{Script: "sleep 30 & sleep 30"})
	if err != context.DeadlineExceeded {
		t.Errorf("err = %v, muốn DeadlineExceeded", err)
	}
	if d := time.Since(start); d > 5*time.Second {
		t.Errorf("giết cây chậm: %v", d)
	}
}

func TestOutputCapped(t *testing.T) {
	res, err := System{}.Run(t.Context(), Spec{Script: "head -c 200000 /dev/zero | tr '\\0' a"})
	if err != nil || len(res.Stdout) != OutputMax {
		t.Errorf("len = %d (%v), muốn %d", len(res.Stdout), err, OutputMax)
	}
}
