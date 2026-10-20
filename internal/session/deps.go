package session

import (
	"context"
	"image"
	"time"

	"autoclicker/internal/device"
	"autoclicker/internal/model"
	"autoclicker/internal/store"
	"autoclicker/internal/web"
)

// Deps – thiết bị + lưu trữ một phiên dùng; trường nil = bản thật. Test thay bằng bản giả.
type Deps struct {
	Input       func() device.Input
	Screen      func() model.Rect
	Displays    func() []model.Rect
	Now         func() time.Time
	AllowDanger func(sc model.Script, path, confirm string) (bool, error)
	// Browser – Chrome của app cho node web (nối lười); close gọi khi phiên xong.
	Browser func(move func(x, y int) error) (b web.Browser, close func(), err error)
	// Corners – dừng khẩn khi chuột chạm góc màn hình (chạy tới khi ctx huỷ).
	Corners func(ctx context.Context, in *device.TrackingInput, onHit func())
	SaveRun func(rec model.RunRecord, path string, shots map[string]image.Image) error
}

func (d Deps) withDefaults() Deps {
	if d.Input == nil {
		d.Input = func() device.Input { return device.RobotgoInput{} }
	}
	if d.Screen == nil {
		d.Screen = device.MainScreen
	}
	if d.Displays == nil {
		d.Displays = device.Displays
	}
	if d.Now == nil {
		d.Now = time.Now
	}
	if d.AllowDanger == nil {
		d.AllowDanger = AllowDanger
	}
	if d.Browser == nil {
		d.Browser = func(move func(x, y int) error) (web.Browser, func(), error) {
			profile, err := store.BrowserProfile()
			if err != nil {
				return nil, nil, err
			}
			br := web.NewChrome(profile, move) // nối lười – kịch bản không có node web thì không đụng trình duyệt
			return br, br.Close, nil
		}
	}
	if d.Corners == nil {
		d.Corners = func(ctx context.Context, in *device.TrackingInput, onHit func()) {
			device.WatchCorners(ctx, device.Location, device.HardCorners(device.Displays()), in, onHit)
		}
	}
	if d.SaveRun == nil {
		d.SaveRun = func(rec model.RunRecord, path string, shots map[string]image.Image) error {
			base, err := store.ConfigDir()
			if err != nil {
				return err
			}
			return store.SaveRun(store.RunsDir(base, rec.Script, path), rec, shots)
		}
	}
	return d
}
