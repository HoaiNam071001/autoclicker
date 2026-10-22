package model

// Node Nguồn (bước 43 S3, engine/power.go): khoá màn, tắt màn hình, ngủ, đăng xuất, tắt, khởi động lại.

// Hành động của node Nguồn (Step.PowerOp) – khớp power.Lock…
const (
	PowerLock         = "lock"
	PowerDisplaySleep = "displaySleep"
	PowerSleep        = "sleep"
	PowerLogout       = "logout"
	PowerShutdown     = "shutdown"
	PowerRestart      = "restart"
)

var PowerOps = []string{PowerLock, PowerDisplaySleep, PowerSleep, PowerLogout, PowerShutdown, PowerRestart}

// PowerNeedsConsent – hành động làm mất phiên làm việc (người dùng phải đồng ý trước lần chạy đầu): ngủ, đăng xuất,
// tắt, khởi động lại. Khoá màn / tắt màn hình thì không.
func PowerNeedsConsent(op string) bool {
	return op == PowerSleep || op == PowerLogout || op == PowerShutdown || op == PowerRestart
}

// PowerCountdownSec – đếm ngược mặc định (giây) trước hành động cần đồng ý, để kịp bấm Dừng.
const PowerCountdownSec = 30
