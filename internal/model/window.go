package model

// Cửa sổ / app đang chạy (bước 26) – kiểu thuần cho device.Windows / device.Apps, node Cửa sổ / Ứng dụng và tool MCP
// list_windows / list_apps.

// WindowInfo – một cửa sổ cấp cao đang mở (không tính cửa sổ của autoclicker).
type WindowInfo struct {
	ID        uint64 `json:"id"`        // HWND (Windows) / CGWindowID (macOS) – chỉ dùng trong lần chạy, không lưu vào kịch bản
	Title     string `json:"title"`     // tiêu đề
	App       string `json:"app"`       // tên app: Windows tên tệp chạy không .exe (chrome), macOS tên hiển thị
	AppID     string `json:"appId"`     // đường dẫn exe (Windows) / bundle id (macOS)
	PID       int    `json:"pid"`       // tiến trình
	Box       Rect   `json:"box"`       // khung nhìn thấy, cùng hệ toạ độ với điểm click (Windows px thật, macOS point)
	Minimized bool   `json:"minimized"` // đang thu nhỏ – Box rỗng
}

// AppInfo – một tiến trình / app đang chạy; Windows = số cửa sổ đang mở (0 = chạy nền / chưa có cửa sổ).
type AppInfo struct {
	Name    string `json:"name"` // như WindowInfo.App
	ID      string `json:"id"`   // như WindowInfo.AppID (có thể trống – tiến trình hệ thống không đọc được đường dẫn)
	PID     int    `json:"pid"`
	Windows int    `json:"windows"`
}
