package engine

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"autoclicker/internal/device"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Node Chụp màn hình (bước 35 N1): chụp vùng theo ShotArea, lưu PNG, tuỳ chọn chép vào clipboard. Không di chuột.

func (e *Engine) screenshot(ctx context.Context, s model.Step) (map[string]any, error) {
	area, err := e.shotArea(s)
	if err != nil {
		return nil, err
	}
	if err := device.CaptureArea(area, e.displays); err != nil {
		return nil, err
	}
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	img, err := e.scr.Capture(area)
	if err != nil {
		return nil, err
	}
	path, err := shotPath(s.SaveDir, s.FileName, time.Now())
	if err != nil {
		return nil, err
	}
	if err := writePNG(path, img); err != nil {
		return nil, err
	}
	if s.ToClipboard {
		if err := e.clip.WriteImage(path); err != nil {
			return nil, err
		}
	}
	return map[string]any{"path": path, "x": area.X, "y": area.Y, "w": area.W, "h": area.H}, nil
}

// shotArea – vùng chụp theo toạ độ robotgo.
func (e *Engine) shotArea(s model.Step) (model.Rect, error) {
	switch s.ShotArea {
	case model.ShotScreen:
		if s.Display == 0 {
			return e.searchArea(model.Step{SearchArea: "screen"}), nil
		}
		if s.Display > len(e.displays) {
			return model.Rect{}, &invalidErr{fmt.Sprintf("không có màn hình thứ %d (máy đang có %d màn)", s.Display, len(e.displays))}
		}
		return e.displays[s.Display-1], nil
	case model.ShotRect:
		if s.ShotRect == nil {
			return model.Rect{}, &invalidErr{"chưa chọn vùng chụp"}
		}
		return *s.ShotRect, nil
	case model.ShotNode:
		_, box, err := e.nodeBox(s.ImageFrom)
		return box, err
	}
	return script.Bounds(e.sc, e.screen), nil
}

// nodeBox – tâm + khung chỗ thấy của node id lần chạy gần nhất trong vòng (output x, y = tâm, w, h).
func (e *Engine) nodeBox(id string) (model.Point, model.Rect, error) {
	r := e.results[id]
	ref := e.byID[id].Ref
	if r == nil || r.Status != model.StatusSuccess {
		return model.Point{}, model.Rect{}, fmt.Errorf("node %s chưa tìm thấy gì trong vòng này (chưa chạy hoặc lần gần nhất thất bại)", ref)
	}
	x, okX := r.Output["x"].(int)
	y, okY := r.Output["y"].(int)
	if !okX || !okY {
		return model.Point{}, model.Rect{}, fmt.Errorf("node %s không có chỗ thấy (output x / y)", ref)
	}
	w, _ := r.Output["w"].(int)
	h, _ := r.Output["h"].(int)
	return model.Point{X: x, Y: y}, model.Rect{X: x - w/2, Y: y - h/2, W: w, H: h}, nil
}

// userPath – đường dẫn người dùng gõ: ~ / đường dẫn tương đối tính từ thư mục người dùng (app không có thư mục làm việc
// rõ ràng – chạy từ Finder / Start). Node Chụp màn hình, Clipboard (ảnh), sau này Tệp / thư mục.
func userPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if filepath.IsAbs(p) {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("không tìm thấy thư mục người dùng: %w", err)
	}
	if p == "~" || strings.HasPrefix(p, "~/") || strings.HasPrefix(p, `~\`) {
		p = p[1:]
	}
	return filepath.Join(home, p), nil
}

// shotPath – đường dẫn file sẽ ghi: dir trống = Pictures/autoclicker, còn lại theo userPath; name trống =
// autoclicker-<ngày>-<giờ>; ký tự không hợp lệ trong tên → "-"; trùng file có sẵn thì thêm " (2)"…
func shotPath(dir, name string, now time.Time) (string, error) {
	if strings.TrimSpace(dir) == "" {
		dir = filepath.Join("Pictures", "autoclicker")
	}
	dir, err := userPath(dir)
	if err != nil {
		return "", err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "autoclicker-" + now.Format("2006-01-02-150405")
	}
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`/\:*?"<>|`, r) || r < 32 {
			return '-'
		}
		return r
	}, strings.TrimSuffix(strings.TrimSuffix(name, ".png"), ".PNG"))
	p := filepath.Join(dir, name+".png")
	for i := 2; fileExists(p); i++ {
		p = filepath.Join(dir, fmt.Sprintf("%s (%d).png", name, i))
	}
	return p, nil
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// writePNG ghi ảnh (tạo thư mục nếu chưa có); ghi file tạm rồi đổi tên – lỗi giữa chừng không để lại file hỏng.
func writePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("không tạo được thư mục lưu ảnh: %w", err)
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".shot-*.png")
	if err != nil {
		return fmt.Errorf("không ghi được ảnh: %w", err)
	}
	err = png.Encode(f, img)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), path)
	}
	if err != nil {
		os.Remove(f.Name())
		return fmt.Errorf("không ghi được ảnh: %w", err)
	}
	return nil
}
