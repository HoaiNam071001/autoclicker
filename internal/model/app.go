package model

import (
	"slices"
	"strings"
)

// Node Ứng dụng (bước 26, engine/app.go): kiểm tra / chờ / thoát app khớp AppName. Một app có thể nhiều tiến trình
// (Chrome) – coi cả nhóm là một app.

// Thao tác của node Ứng dụng (Step.AppOp).
const (
	AppCheck     = ""          // đang chạy → cổng yes, không → no
	AppWaitStart = "waitStart" // chờ app chạy và có ít nhất một cửa sổ
	AppWaitQuit  = "waitQuit"  // chờ app thoát hẳn
	AppQuit      = "quit"      // thoát bình thường (app có thể hỏi lưu); AppForce: quá giờ thì buộc thoát
	AppForceQuit = "forceQuit" // buộc thoát – mất dữ liệu chưa lưu
)

var AppOps = []string{AppCheck, AppWaitStart, AppWaitQuit, AppQuit, AppForceQuit}

// MatchApps – chỉ số mọi tiến trình trong apps khớp want (so như MatchApp). Cùng công thức matchApps
// (frontend/src/flow/app.ts), test chung testdata/apps.json.
func MatchApps(apps []AppInfo, want, mode string) ([]int, error) {
	var out []int
	for i, a := range apps {
		ok, err := MatchApp(a.Name, a.ID, want, mode)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, i)
		}
	}
	return out, nil
}

// MainApp – tiến trình đại diện của nhóm: nhiều cửa sổ nhất (bằng nhau → đầu tiên); windows = tổng số cửa sổ.
func MainApp(apps []AppInfo) (main AppInfo, windows int) {
	for i, a := range apps {
		if i == 0 || a.Windows > main.Windows {
			main = a
		}
		windows += a.Windows
	}
	return main, windows
}

// protectedApps – không cho thoát: chính autoclicker, vỏ hệ điều hành, tiến trình hệ thống.
var protectedApps = []string{
	"autoclicker",
	// macOS
	"finder", "dock", "systemuiserver", "loginwindow", "windowserver", "kernel_task", "launchd", "controlcenter",
	// Windows
	"explorer", "system", "registry", "smss", "csrss", "wininit", "winlogon", "services", "lsass", "svchost", "dwm",
	"fontdrvhost", "sihost", "ctfmon", "taskhostw", "runtimebroker", "startmenuexperiencehost", "shellexperiencehost",
	"searchhost", "textinputhost", "lockapp", "msmpeng", "conhost",
}

// ProtectedApp – tên app / tên exe / đường dẫn thuộc danh sách không được thoát (so cả chuỗi, bỏ thư mục, .exe / .app).
func ProtectedApp(name string) bool {
	n := foldName(name)
	if i := strings.LastIndexAny(n, `\/`); i >= 0 {
		n = n[i+1:]
	}
	n = strings.TrimSuffix(strings.TrimSuffix(n, ".exe"), ".app")
	return slices.Contains(protectedApps, n)
}
