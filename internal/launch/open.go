// Package launch – mở URL / file / app theo hệ điều hành (node Mở).
package launch

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	goruntime "runtime"
	"strings"
)

// Bước "Mở URL / app" (bước 23). Target có thể là:
//   - URL có scheme (https://…, mailto:…, spotify:…) hoặc dạng tên miền (facebook.com/abc → thêm https://)
//   - đường dẫn file / app (/Applications/Zalo.app, C:\…\app.exe, ~/x.pdf)
//   - tên app (Spotify, Zalo) – macOS open -a, Windows start
// App: mở Target bằng app này (thường là trình duyệt); "" = app mặc định của máy.

var (
	schemeRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)
	winPath  = regexp.MustCompile(`^[A-Za-z]:[\\/]`)
	domainRe = regexp.MustCompile(`^[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)+(:\d+)?([/?#].*)?$`)
)

// Kiểu target.
const (
	TargetURL  = "url"
	TargetPath = "path"
	TargetApp  = "app"
)

func TargetKind(t string) string {
	switch {
	case winPath.MatchString(t), strings.HasPrefix(t, "/"), strings.HasPrefix(t, "~"), strings.HasPrefix(t, `\\`):
		return TargetPath
	case schemeRe.MatchString(t), domainRe.MatchString(t):
		return TargetURL
	}
	return TargetApp
}

// normalizeURL thêm https:// cho dạng tên miền.
func normalizeURL(t string) string {
	if schemeRe.MatchString(t) {
		return t
	}
	return "https://" + t
}

// EscapeURLValue mã hoá giá trị biến chèn vào URL (dấu cách, tiếng Việt, &, ?…): "Tinh Hà" → "Tinh%20H%C3%A0".
func EscapeURLValue(v string) string {
	return strings.ReplaceAll(url.QueryEscape(v), "+", "%20")
}

// Tên trình duyệt hay dùng → tên lệnh trên Windows (start tìm theo App Paths).
var windowsApps = map[string]string{
	"google chrome":   "chrome",
	"chrome":          "chrome",
	"microsoft edge":  "msedge",
	"edge":            "msedge",
	"firefox":         "firefox",
	"mozilla firefox": "firefox",
}

// openCommand trả lệnh mở target bằng app trên hệ điều hành goos (tách riêng để test, không chạy gì).
func openCommand(goos, target, app string) (string, []string) {
	kind := TargetKind(target)
	if kind == TargetURL {
		target = normalizeURL(target)
	}
	switch goos {
	case "darwin":
		switch {
		case app != "":
			return "open", []string{"-a", app, target}
		case kind == TargetApp:
			return "open", []string{"-a", target}
		}
		return "open", []string{target}
	case "windows":
		// PowerShell Start-Process thay cho "cmd /c start": cmd hiểu & trong URL (?a=1&b=2) là nối lệnh.
		cmd := "Start-Process -FilePath " + psQuote(target)
		if app != "" {
			if w, ok := windowsApps[strings.ToLower(app)]; ok {
				app = w
			}
			cmd = "Start-Process -FilePath " + psQuote(app) + " -ArgumentList " + psQuote(`"`+target+`"`)
		}
		return "powershell", []string{"-NoProfile", "-NonInteractive", "-Command", cmd}
	}
	if app != "" {
		return app, []string{target}
	}
	return "xdg-open", []string{target}
}

// psQuote: chuỗi PowerShell trong nháy đơn (' viết thành ”).
func psQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

// BrowserURL – target của node Mở bằng Chrome của app (Step.WebBrowser) thành URL cho trình duyệt: tên miền thêm
// https://, đường dẫn file thành file://. Tên app không mở được bằng trình duyệt → "".
func BrowserURL(target string) string {
	switch TargetKind(target) {
	case TargetURL:
		return normalizeURL(target)
	case TargetPath:
		p := expandHome(target)
		if abs, err := filepath.Abs(p); err == nil {
			p = abs
		}
		p = filepath.ToSlash(p)
		if !strings.HasPrefix(p, "/") {
			p = "/" + p // C:/… → /C:/…
		}
		return (&url.URL{Scheme: "file", Path: p}).String()
	}
	return ""
}

func expandHome(target string) string {
	if strings.HasPrefix(target, "~/") || strings.HasPrefix(target, `~\`) {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, target[2:])
		}
	}
	return target
}

// OpenTarget mở target thật (không chờ app mở xong – bước sau nên có thời gian chờ).
func OpenTarget(target, app string) error {
	name, args := openCommand(goruntime.GOOS, expandHome(target), app)
	return exec.Command(name, args...).Start()
}
