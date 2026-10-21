package browser

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"ghosthand/internal/model"
)

// macOS: mỗi thao tác một lần osascript chạy tabs.js. Lần đầu điều khiển mỗi trình duyệt, macOS hỏi quyền Automation.

//go:embed tabs.js
var tabsJS string

// osaTimeout – một lệnh osascript tối đa (trình duyệt treo / hộp hỏi quyền chưa trả lời).
const osaTimeout = 20 * time.Second

func osa(ctx context.Context, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, osaTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "osascript", append([]string{"-l", "JavaScript", "-e", tabsJS}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		if known := osaError(stderr.String()); known != nil {
			return "", known
		}
		return "", fmt.Errorf("điều khiển trình duyệt lỗi: %s", strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(out)), nil
}

func (System) Tabs(ctx context.Context, b string) ([]model.TabInfo, error) {
	out, err := osa(ctx, "list", b)
	if err != nil {
		return nil, err
	}
	return parseTabs(out)
}

func (System) Activate(ctx context.Context, t model.TabInfo) error {
	_, err := osa(ctx, "activate", t.Browser, id(t.Window), id(t.ID))
	return err
}

func (s System) NewTab(ctx context.Context, b, url string) (model.TabInfo, error) {
	if b == "" {
		tabs, err := s.Tabs(ctx, "")
		if err != nil {
			return model.TabInfo{}, err
		}
		if len(tabs) == 0 {
			return model.TabInfo{}, ErrNotRunning
		}
		b = tabs[0].Browser
	}
	out, err := osa(ctx, "new", b, url)
	if err != nil {
		return model.TabInfo{}, err
	}
	t := model.TabInfo{Browser: b, Active: true}
	return t, json.Unmarshal([]byte(out), &t)
}

func (System) Close(ctx context.Context, t model.TabInfo) error {
	_, err := osa(ctx, "close", t.Browser, id(t.Window), id(t.ID))
	return err
}

func (System) Nav(ctx context.Context, t model.TabInfo, op string) error {
	_, err := osa(ctx, op, t.Browser, id(t.Window), id(t.ID))
	return err
}

func id(n int64) string { return strconv.FormatInt(n, 10) }
