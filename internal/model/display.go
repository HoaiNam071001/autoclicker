package model

// Node Màn hình / âm thanh (bước 43 S7, engine/display.go).

// Thao tác của node Màn hình / âm thanh (Step.DisplayOp).
const (
	DisplayRead       = ""           // đọc: sáng / tối, độ sáng, loa + micro mặc định, danh sách thiết bị
	DisplayTheme      = "theme"      // DisplayValue: light | dark | toggle
	DisplayBrightness = "brightness" // DisplayValue: 0–100 (%)
	DisplayAudioOut   = "audioOut"   // đổi loa mặc định sang thiết bị tên / mã DisplayValue
	DisplayAudioIn    = "audioIn"    // đổi micro mặc định
	DisplaySound      = "sound"      // phát file âm thanh hoặc âm hệ thống tên DisplayValue
)

var DisplayOps = []string{DisplayRead, DisplayTheme, DisplayBrightness, DisplayAudioOut, DisplayAudioIn, DisplaySound}

// Giá trị DisplayValue của DisplayTheme.
var DisplayThemes = []string{"light", "dark", "toggle"}
