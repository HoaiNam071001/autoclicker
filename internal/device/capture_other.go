//go:build !darwin

package device

import (
	"image"
	goruntime "runtime"

	"autoclicker/internal/model"
	"autoclicker/internal/winexclude"

	"github.com/go-vgo/robotgo"
)

func ScreenCaptureAllowed() bool { return true }

func RequestScreenCapture() bool { return true }

func screenCaptureCheck() error { return nil }

// selfExcluded: Windows đã đặt cửa sổ app không lọt vào ảnh chụp (ExcludeSelfFromCapture).
var selfExcluded bool

// ExcludeSelfFromCapture (gọi lúc khởi động): Windows 10 2004+ – cửa sổ app (thanh mini, overlay) không lọt vào
// ảnh chụp; robotgo chụp bằng BitBlt không có con trỏ chuột.
func ExcludeSelfFromCapture() {
	if goruntime.GOOS == "windows" {
		selfExcluded = winexclude.ExcludeOwnWindows() > 0
	}
}

// captureNoSelf: ok = false khi không loại được cửa sổ app (Linux, Windows cũ) – nơi gọi phải ẩn app trước khi chụp.
func captureNoSelf(r model.Rect) (image.Image, bool, error) {
	if !selfExcluded {
		return nil, false, nil
	}
	img, err := robotgo.CaptureImg(r.X, r.Y, r.W, r.H)
	return img, true, err
}

func CanCaptureWithoutSelf() bool { return selfExcluded }
