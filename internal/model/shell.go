package model

// Node Chạy lệnh (bước 43 S1, engine/shell.go, package internal/shell): chạy một script bằng shell / PowerShell /
// AppleScript · JXA. Chỉ chạy khi người dùng đã đồng ý đúng nội dung lệnh (script.DangerOps, DangerShell).

// Ngôn ngữ script (Step.ShellLang).
const (
	ShellDefault     = ""            // shell của máy: zsh (macOS) / cmd (Windows)
	ShellPowerShell  = "powershell"  // PowerShell (Windows; macOS / Linux cần pwsh)
	ShellAppleScript = "applescript" // chỉ macOS
	ShellJXA         = "jxa"         // JavaScript for Automation – chỉ macOS
)

var ShellLangs = []string{ShellDefault, ShellPowerShell, ShellAppleScript, ShellJXA}

// ShellTimeoutMs – hết giờ mặc định của node Chạy lệnh.
const ShellTimeoutMs = 30000

// OSShortcutTimeoutMs – hết giờ mặc định của node Shortcuts của hệ điều hành.
const OSShortcutTimeoutMs = 60000
