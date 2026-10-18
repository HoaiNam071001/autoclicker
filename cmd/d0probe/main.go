//go:build probe

// D0 (bước 25d): thử đo toạ độ trang → màn hình bằng chuột thật, chưa vào model.
// Chạy: go run -tags probe ./cmd/d0probe [-scenario all|base|bookmarks|zoom125|moved|devtools|behind|display2]
// Cần quyền Accessibility; chuột sẽ tự di chuyển + bấm trong cửa sổ Chrome thử (profile tạm).
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/remote"
	"github.com/go-vgo/robotgo"
)

var chromePath = flag.String("chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "đường dẫn Chrome")

func main() {
	scenario := flag.String("scenario", "all", "all|base|bookmarks|zoom125|moved|devtools|behind|display2")
	flag.Parse()

	names := []string{*scenario}
	if *scenario == "all" {
		names = []string{"base", "bookmarks", "zoom125", "moved", "devtools", "behind", "display2"}
	}
	fmt.Printf("robotgo ScaleF=%.2f, số màn hình=%d\n\n", robotgo.ScaleF(), robotgo.DisplaysNum())
	for _, n := range names {
		if n == "display2" && robotgo.DisplaysNum() < 2 {
			fmt.Println("## display2 – bỏ qua (chỉ có 1 màn hình)")
			continue
		}
		if err := runScenario(n); err != nil {
			fmt.Printf("## %s – LỖI: %v\n\n", n, err)
		}
	}
}

func runScenario(name string) error {
	fmt.Printf("## %s\n", name)
	dir, err := os.MkdirTemp("", "d0probe-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	pagePath := filepath.Join(dir, "probe.html")
	if err := os.WriteFile(pagePath, []byte(probePage), 0o644); err != nil {
		return err
	}
	args := []string{
		"--remote-debugging-port=0", // cổng cố định có thể đã bị tiến trình khác giữ → để Chrome chọn, đọc DevToolsActivePort
		"--user-data-dir=" + filepath.Join(dir, "profile"),
		"--no-first-run", "--no-default-browser-check",
		"--window-position=100,80", "--window-size=1200,850",
	}
	if name == "devtools" {
		args = append(args, "--auto-open-devtools-for-tabs")
	}
	cmd := exec.Command(*chromePath, append(args, "file://"+pagePath)...)
	if err := cmd.Start(); err != nil {
		return err
	}
	defer func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() }()

	ctx, cancel, err := attach(filepath.Join(dir, "profile"), "probe.html")
	if err != nil {
		return err
	}
	defer cancel()
	time.Sleep(time.Second) // chờ cửa sổ vẽ xong
	if name == "behind" { // app khác lên trên Chrome trước khi đo
		_ = exec.Command("osascript", "-e", `tell application "Visual Studio Code" to activate`).Run()
		time.Sleep(800 * time.Millisecond)
	}
	if err := bringToFront(ctx, cmd.Process.Pid); err != nil {
		return err
	}

	switch name {
	case "bookmarks":
		robotgo.KeyTap("b", "cmd", "shift")
		time.Sleep(600 * time.Millisecond)
	case "zoom125":
		robotgo.KeyTap("=", "cmd") // 110%
		time.Sleep(300 * time.Millisecond)
		robotgo.KeyTap("=", "cmd") // 125%
		time.Sleep(600 * time.Millisecond)
	case "display2":
		x, y, _, _ := robotgo.GetDisplayBounds(1)
		if err := moveWindow(ctx, x+100, y+80); err != nil {
			return err
		}
	}

	cal, err := calibrate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(cal)
	if err := clickTargets(ctx, cal, ""); err != nil {
		return err
	}
	if name == "moved" {
		if err := checkMoved(ctx, cal); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

// checkMoved – dời cửa sổ: bấm bằng số đo cũ (để thấy sai), phát hiện qua screenX/Y rồi đo lại.
func checkMoved(ctx context.Context, old calibration) error {
	if err := moveWindow(ctx, 300, 160); err != nil {
		return err
	}
	if err := clickTargets(ctx, old, "số đo cũ sau khi dời"); err != nil {
		return err
	}
	g, err := readGeom(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("phát hiện dời: screenX/Y %v,%v → %v,%v\n", old.geom.ScreenX, old.geom.ScreenY, g.ScreenX, g.ScreenY)
	cal, err := calibrate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(cal)
	return clickTargets(ctx, cal, "đo lại")
}

// attach – chờ Chrome ghi cổng debug vào profile rồi nối vào tab đang mở trang thử (không mở tab mới).
func attach(profile, urlPart string) (context.Context, context.CancelFunc, error) {
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		port, err := activePort(profile)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		id, err := findTarget(port, urlPart)
		if err == nil {
			alloc, cancelAlloc := remote.NewAllocator(context.Background(), fmt.Sprintf("http://127.0.0.1:%d", port))
			ctx, cancelCtx := chromedp.NewContext(alloc, chromedp.WithTargetID(target.ID(id)))
			if err := chromedp.Do(ctx); err != nil {
				cancelCtx()
				cancelAlloc()
				return nil, nil, err
			}
			return ctx, func() { cancelCtx(); cancelAlloc() }, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil, nil, fmt.Errorf("không nối được Chrome (profile %s)", profile)
}

// activePort – dòng đầu của <profile>/DevToolsActivePort là cổng Chrome đã chọn.
func activePort(profile string) (int, error) {
	b, err := os.ReadFile(filepath.Join(profile, "DevToolsActivePort"))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(strings.SplitN(string(b), "\n", 2)[0]))
}

func findTarget(port int, urlPart string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/json/list", port))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var list []struct{ ID, Type, URL string }
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return "", err
	}
	for _, t := range list {
		if t.Type == "page" && strings.Contains(t.URL, urlPart) {
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("chưa thấy tab")
}
