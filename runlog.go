package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"

	"autoclicker/internal/model"
)

// Lịch sử chạy (bước 20 F5): mỗi lần chạy lưu một file runs/<kịch bản>/<id>.json (RunRecord – StepResult từng node,
// input đã dùng, kết quả cuối) + ảnh màn hình <id>-<n>.png lúc node thất bại. Giữ maxRuns lần gần nhất mỗi kịch bản.

const (
	maxRuns       = 50
	maxRunResults = 2000 // quá thì bỏ kết quả cũ nhất (vòng lặp dài) – Truncated đếm số đã bỏ
	maxShots      = 5    // ảnh lúc lỗi tối đa mỗi lần chạy (vòng thử lại qua cổng fail không chụp mãi)
	shotMaxW      = 1600 // ảnh rộng hơn thì thu nhỏ (Retina) cho file nhẹ
)

// RunRecord.Status
const (
	RunSuccess = "success"
	RunFail    = "fail"
	RunStopped = "stopped" // người dùng dừng (nút, F8, góc màn hình)
)

// Nguồn của lần chạy (RunOptions.Source).
const (
	SourceUser     = "user"
	SourceSchedule = "schedule"
	SourceCLI      = "cli"
	SourceAI       = "ai"
)

// RunOptions – cách chạy kịch bản (RunScript).
type RunOptions struct {
	Countdown int    `json:"countdown"`        // giây đếm ngược trước khi chạy
	Path      string `json:"path,omitempty"`   // file kịch bản – gom lịch sử theo kịch bản (chưa lưu: theo tên)
	Source    string `json:"source,omitempty"` // user (mặc định) | schedule | cli | ai
	From      string `json:"from,omitempty"`   // debug: chạy một vòng bắt đầu từ node này thay cho start
	Step      bool   `json:"step,omitempty"`   // debug: dừng trước mỗi node, Tiếp / F10 chạy node kế
}

// RunSummary – một dòng trong danh sách lịch sử (event engine:done, ListRuns).
type RunSummary struct {
	ID         string    `json:"id"`
	Script     string    `json:"script"`
	Source     string    `json:"source"`
	StartedAt  time.Time `json:"startedAt"`
	DurationMs int64     `json:"durationMs"`
	Status     string    `json:"status"`            // success | fail | stopped
	Error      string    `json:"error,omitempty"`   // lỗi làm dừng kịch bản
	Message    string    `json:"message,omitempty"` // lời nhắn node Kết thúc thành công (cuối cùng)
	Nodes      int       `json:"nodes"`             // số lần chạy node
}

// RunRecord – toàn bộ một lần chạy.
type RunRecord struct {
	RunSummary
	Inputs    map[string]any     `json:"inputs,omitempty"`
	Results   []model.StepResult `json:"results"`
	Truncated int                `json:"truncated,omitempty"`
}

// add ghi kết quả một node (giữ maxRunResults cuối); node Kết thúc thành công có lời nhắn thì nhớ làm Message.
func (r *RunRecord) add(res model.StepResult, typ string) {
	r.Nodes++
	if typ == model.StepEnd && res.Status == model.StatusSuccess {
		if m, ok := res.Output["message"].(string); ok {
			r.Message = m
		}
	}
	if len(r.Results) >= maxRunResults {
		r.Results = r.Results[1:]
		r.Truncated++
	}
	r.Results = append(r.Results, res)
}

// shoot chụp màn chứa vùng kịch bản lúc node thất bại (tối đa maxShots mỗi lần chạy); lỗi chụp thì bỏ qua.
func (e *engine) shoot(res *model.StepResult) {
	if e.shots == nil || len(e.shots) >= maxShots || res.Error == nil || res.Error.Code == model.ErrPermission || res.Error.Code == model.ErrEnd {
		return
	}
	img, err := e.scr.Capture(e.searchArea(model.Step{SearchArea: "screen"}))
	if err != nil {
		return
	}
	res.Shot = fmt.Sprintf("%d.png", len(e.shots)+1)
	e.shots[res.Shot] = shrink(img, shotMaxW)
}

// shrink thu nhỏ ảnh về rộng tối đa w (lấy trung bình từng ô – đủ để xem lại, không cần đẹp).
func shrink(img image.Image, w int) image.Image {
	b := img.Bounds()
	k := (b.Dx() + w - 1) / w
	if k <= 1 {
		return img
	}
	out := image.NewRGBA(image.Rect(0, 0, b.Dx()/k, b.Dy()/k))
	for y := 0; y < out.Rect.Dy(); y++ {
		for x := 0; x < out.Rect.Dx(); x++ {
			var r, g, bl uint32
			for dy := 0; dy < k; dy++ {
				for dx := 0; dx < k; dx++ {
					cr, cg, cb, _ := img.At(b.Min.X+x*k+dx, b.Min.Y+y*k+dy).RGBA()
					r, g, bl = r+cr, g+cg, bl+cb
				}
			}
			n := uint32(k * k)
			i := out.PixOffset(x, y)
			out.Pix[i], out.Pix[i+1], out.Pix[i+2], out.Pix[i+3] = uint8(r/n>>8), uint8(g/n>>8), uint8(bl/n>>8), 255
		}
	}
	return out
}

// runsDir – thư mục lịch sử của một kịch bản: runs/<tên>-<8 ký tự hash đường dẫn>. Chưa lưu thì hash theo tên.
func runsDir(base, name, path string) string {
	key := path
	if key == "" {
		key = "unsaved:" + name
	}
	h := sha1.Sum([]byte(key))
	slug := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return '-'
	}, name)
	if s := []rune(strings.Trim(slug, "-")); len(s) > 40 {
		slug = string(s[:40])
	}
	return filepath.Join(base, "runs", strings.Trim(slug, "-")+"-"+hex.EncodeToString(h[:4]))
}

// newRunID – id theo giờ bắt đầu (sắp theo tên = theo thời gian).
func newRunID(t time.Time) string { return t.Format("20060102-150405.000") }

// saveRun ghi lần chạy + ảnh lúc lỗi vào dir rồi xoá các lần cũ quá maxRuns.
func saveRun(dir string, rec RunRecord, shots map[string]image.Image) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, img := range shots {
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, rec.ID+"-"+name), buf.Bytes(), 0o644); err != nil {
			return err
		}
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, rec.ID+".json"), b, 0o644); err != nil {
		return err
	}
	ids := runIDs(dir)
	for _, id := range ids[min(len(ids), maxRuns):] {
		olds, _ := filepath.Glob(filepath.Join(dir, id+"-*.png"))
		for _, f := range append(olds, filepath.Join(dir, id+".json")) {
			os.Remove(f)
		}
	}
	return nil
}

// runIDs – id các lần chạy trong dir, mới nhất trước.
func runIDs(dir string) []string {
	files, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	ids := make([]string, len(files))
	for i, f := range files {
		ids[i] = strings.TrimSuffix(filepath.Base(f), ".json")
	}
	slices.Sort(ids)
	slices.Reverse(ids)
	return ids
}

// listRuns đọc tóm tắt các lần chạy trong dir, mới nhất trước (file hỏng thì bỏ qua).
func listRuns(dir string) []RunSummary {
	out := []RunSummary{}
	for _, id := range runIDs(dir) {
		var r struct{ RunSummary }
		if b, err := os.ReadFile(filepath.Join(dir, id+".json")); err == nil && json.Unmarshal(b, &r) == nil {
			out = append(out, r.RunSummary)
		}
	}
	return out
}

// safeName chặn tên file có đường dẫn (id / tên ảnh đến từ frontend).
func safeName(n string) error {
	if n == "" || n != filepath.Base(n) || strings.HasPrefix(n, ".") {
		return errors.New("tên file lịch sử không hợp lệ")
	}
	return nil
}

func loadRun(dir, id string) (RunRecord, error) {
	var r RunRecord
	if err := safeName(id); err != nil {
		return r, err
	}
	b, err := os.ReadFile(filepath.Join(dir, id+".json"))
	if err != nil {
		return r, fmt.Errorf("không đọc được lần chạy: %w", err)
	}
	return r, json.Unmarshal(b, &r)
}

// ListRuns – lịch sử chạy của kịch bản (tên + đường dẫn file như lúc chạy), mới nhất trước.
func (a *App) ListRuns(name, path string) []RunSummary {
	base, err := configDir()
	if err != nil {
		return []RunSummary{}
	}
	return listRuns(runsDir(base, name, path))
}

// GetRun – toàn bộ một lần chạy (nhật ký từng node).
func (a *App) GetRun(name, path, id string) (RunRecord, error) {
	base, err := configDir()
	if err != nil {
		return RunRecord{}, err
	}
	return loadRun(runsDir(base, name, path), id)
}

// RunShot – ảnh màn hình lúc node thất bại (StepResult.Shot) dạng data URL.
func (a *App) RunShot(name, path, id, shot string) (string, error) {
	base, err := configDir()
	if err != nil {
		return "", err
	}
	if err := safeName(id + "-" + shot); err != nil {
		return "", err
	}
	b, err := os.ReadFile(filepath.Join(runsDir(base, name, path), id+"-"+shot))
	if err != nil {
		return "", fmt.Errorf("không đọc được ảnh: %w", err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b), nil
}
