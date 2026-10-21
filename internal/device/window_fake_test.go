package device

import "ghosthand/internal/device/devicetest"

// Bản giả phải khớp interface – đổi interface mà quên bản giả thì không biên dịch được.
var (
	_ Windows = (*devicetest.FakeWindows)(nil)
	_ Apps    = (*devicetest.FakeApps)(nil)
	_ Windows = SystemWindows{}
	_ Apps    = SystemApps{}
)
