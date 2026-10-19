//go:build !windows

package device

import (
	"errors"

	"autoclicker/internal/model"
)

// Cầu nối tạm (bước 26 W0): macOS (CGWindowList + AX) làm trên máy Mac – xoá file này khi có window_darwin.go
// (Linux vẫn cần bản báo chưa hỗ trợ – đổi build tag thành !windows && !darwin).
var errNoWindows = errors.New("điều khiển cửa sổ / ứng dụng hiện chỉ có trên Windows")

func (SystemWindows) List() ([]model.WindowInfo, error)  { return nil, errNoWindows }
func (SystemWindows) Focus(uint64) error                 { return errNoWindows }
func (SystemWindows) SetBounds(uint64, model.Rect) error { return errNoWindows }
func (SystemWindows) SetState(uint64, string) error      { return errNoWindows }
func (SystemWindows) Close(uint64) error                 { return errNoWindows }
func (SystemWindows) WorkArea(model.Rect) model.Rect     { return MainScreen() }
func (SystemApps) List() ([]model.AppInfo, error)        { return nil, errNoWindows }
func (SystemApps) Quit(int) error                        { return errNoWindows }
func (SystemApps) ForceQuit(int) error                   { return errNoWindows }
