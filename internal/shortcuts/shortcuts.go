// Package shortcuts gọi ứng dụng Phím tắt (Shortcuts) của macOS qua lệnh `shortcuts` cho node Shortcuts của hệ điều
// hành (bước 43 S2). Windows chưa có tương đương → ErrUnsupported.
package shortcuts

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ErrUnsupported = errors.New("Shortcuts chỉ có trên macOS")
	ErrNotFound    = errors.New("không có shortcut tên này")
)

// Backend – danh sách và chạy shortcut. Test thay bản giả.
type Backend interface {
	List(ctx context.Context) ([]string, error)
	// Run chạy shortcut name với input chữ (rỗng = không input), trả chữ shortcut xuất ra ("" nếu không có).
	Run(ctx context.Context, name, input string) (string, error)
}

// System – Backend thật: `shortcuts list` / `shortcuts run`. Tên và input không đi qua shell (đối số + file tạm).
type System struct{}

func (System) List(ctx context.Context) ([]string, error) {
	if runtime.GOOS != "darwin" {
		return nil, ErrUnsupported
	}
	out, err := exec.CommandContext(ctx, "shortcuts", "list").Output()
	if err != nil {
		return nil, cmdErr(err)
	}
	return parseList(string(out)), nil
}

func (System) Run(ctx context.Context, name, input string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", ErrUnsupported
	}
	dir, err := os.MkdirTemp("", "ghosthand-shortcut-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	outPath := filepath.Join(dir, "out.txt")
	args := []string{"run", name, "--output-path", outPath, "--output-type", "public.plain-text"}
	if input != "" {
		inPath := filepath.Join(dir, "in.txt")
		if err := os.WriteFile(inPath, []byte(input), 0o600); err != nil {
			return "", err
		}
		args = append(args, "--input-path", inPath)
	}
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "shortcuts", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		if msg := strings.TrimSpace(stderr.String()); strings.Contains(msg, "Couldn’t find") || strings.Contains(msg, "Couldn't find") {
			return "", fmt.Errorf("%w: %s", ErrNotFound, name)
		} else if msg != "" {
			return "", fmt.Errorf("shortcut %q lỗi: %s", name, msg)
		}
		return "", cmdErr(err)
	}
	b, err := os.ReadFile(outPath)
	if err != nil && !os.IsNotExist(err) { // shortcut không xuất gì → không có file
		return "", err
	}
	return string(b), nil
}

func cmdErr(err error) error {
	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("không thấy lệnh shortcuts (cần macOS 12 trở lên): %w", err)
	}
	return err
}

// parseList – mỗi dòng một tên shortcut (bỏ dòng trống, trùng).
func parseList(s string) []string {
	var out []string
	seen := map[string]bool{}
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" && !seen[l] {
			seen[l], out = true, append(out, l)
		}
	}
	return out
}
