package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Lưu / mở kịch bản (bước 21): file .json, bản nháp tự lưu và danh sách mở gần đây trong configDir().

// ScriptVersion – định dạng file kịch bản hiện tại. Tăng khi đổi model theo cách file cũ không đọc đúng được,
// và thêm bước chuyển đổi trong decodeScript.
// 2 (bước 21b): thêm inputs – app cũ không biết inputs nên không được mở file mới.
const ScriptVersion = 2

const (
	draftFile  = "draft.json"
	maxRecent  = 5
	scriptExt  = ".json"
	dialogName = "Kịch bản Auto Clicker"
)

// settingsMu: đọc – sửa – ghi settings.json (theme, mở gần đây) từ nhiều lời gọi Wails cùng lúc.
var settingsMu sync.Mutex

// ScriptFile – kịch bản vừa mở và đường dẫn file của nó. Path "" = người dùng huỷ hộp thoại.
type ScriptFile struct {
	Script Script `json:"script"`
	Path   string `json:"path"`
}

// Draft – bản nháp tự lưu: kịch bản đang soạn kèm file nó thuộc về, mở lại app thì khôi phục nguyên trạng.
type Draft struct {
	Script  Script `json:"script"`
	Path    string `json:"path"`    // file đang mở, "" = chưa lưu lần nào
	Dirty   bool   `json:"dirty"`   // có thay đổi chưa lưu xuống file
	SavedAt int64  `json:"savedAt"` // thời điểm ghi nháp, unix ms
}

func encodeScript(sc Script) ([]byte, error) {
	sc.Version = ScriptVersion
	return json.MarshalIndent(sc, "", "  ")
}

// decodeScript đọc file kịch bản: kiểm tra version, cấu trúc (checkScript) và bổ sung giá trị trống.
// Không chạy validateScript – bản đang soạn dở (vd bước gõ chữ chưa nhập) vẫn mở được, lúc Chạy mới kiểm tra.
func decodeScript(b []byte) (Script, error) {
	var sc Script
	if err := json.Unmarshal(b, &sc); err != nil {
		return Script{}, fmt.Errorf("file hỏng hoặc không phải JSON: %w", err)
	}
	switch {
	case sc.Version == 0:
		return Script{}, errors.New("không phải file kịch bản (thiếu version)")
	case sc.Version > ScriptVersion:
		return Script{}, fmt.Errorf("file tạo bởi phiên bản app mới hơn (version %d, app đọc được tới %d) – hãy cập nhật app", sc.Version, ScriptVersion)
	}
	if err := checkScript(&sc); err != nil {
		return Script{}, err
	}
	return sc, nil
}

// checkScript kiểm tra giá trị nằm trong miền hợp lệ và bổ sung phần trống (slice nil, id trùng, nút / cách bấm mặc định)
// để frontend dùng ngay được.
func checkScript(sc *Script) error {
	if sc.Loops < 0 || sc.MaxMinutes < 0 {
		return errors.New("số vòng / số phút không được âm")
	}
	if sc.Region != nil && (sc.Region.W <= 0 || sc.Region.H <= 0) {
		return errors.New("vùng có kích thước không hợp lệ")
	}
	if sc.AnchorMode != AnchorNearest && sc.AnchorMode != AnchorTopLeft {
		return fmt.Errorf("cách neo điểm không hợp lệ: %q", sc.AnchorMode)
	}
	if sc.Steps == nil {
		sc.Steps = []Step{}
	}
	ids := map[string]bool{}
	n := 0
	for i := range sc.Steps {
		s := &sc.Steps[i]
		if err := checkStep(s); err != nil {
			return fmt.Errorf("bước %d: %w", i+1, err)
		}
		for s.ID == "" || ids[s.ID] { // id trống / trùng (sửa tay file) → đặt id mới
			n++
			s.ID = fmt.Sprintf("f%d", n)
		}
		ids[s.ID] = true
	}
	return checkInputs(sc)
}

func checkStep(s *Step) error {
	need := 0
	switch s.Type {
	case StepClick, StepDouble, StepHold:
		need = 1
	case StepDrag:
		need = 2
	case StepWait, StepText, StepKey, StepOpen, StepScroll:
	default:
		return fmt.Errorf("loại bước không hợp lệ: %q", s.Type)
	}
	if s.Direction != "" && !slices.Contains(scrollDirections, s.Direction) {
		return fmt.Errorf("hướng cuộn không hợp lệ: %q", s.Direction)
	}
	if len(s.Points) < need {
		return fmt.Errorf("thiếu điểm (cần %d)", need)
	}
	for _, p := range s.Points {
		if !validAnchor(p.Ax) || !validAnchor(p.Ay) {
			return fmt.Errorf("neo điểm không hợp lệ: %q / %q", p.Ax, p.Ay)
		}
	}
	switch s.Button {
	case "":
		s.Button = "left"
	case "left", "right", "center":
	default:
		return fmt.Errorf("nút chuột không hợp lệ: %q", s.Button)
	}
	switch s.KeyMode {
	case "":
		s.KeyMode = KeyChord
	case KeyChord, KeySequence:
	default:
		return fmt.Errorf("cách bấm phím không hợp lệ: %q", s.KeyMode)
	}
	for _, m := range s.Modifiers {
		if m != "cmd" && m != "ctrl" && m != "alt" && m != "shift" {
			return fmt.Errorf("phím kèm không hợp lệ: %q", m)
		}
	}
	if s.Clicks < 0 || s.ClickInterval < 0 || s.HoldMs < 0 || s.DragMs < 0 || s.DelayAfter < 0 || s.JitterMs < 0 || s.JitterPx < 0 || s.Amount < 0 {
		return errors.New("số lần / thời gian / ngẫu nhiên không được âm")
	}
	if s.Points == nil {
		s.Points = []Point{}
	}
	if s.Keys == nil {
		s.Keys = []string{}
	}
	if s.Modifiers == nil {
		s.Modifiers = []string{}
	}
	return nil
}

func validAnchor(a string) bool {
	return a == "" || a == AnchorStart || a == AnchorEnd || a == AnchorCenter || a == AnchorRatio
}

func readScriptFile(path string) (Script, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Script{}, err
	}
	sc, err := decodeScript(b)
	if err != nil {
		return Script{}, err
	}
	if sc.Name == "" {
		sc.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	return sc, nil
}

func writeScriptFile(path string, sc Script) error {
	b, err := encodeScript(sc)
	if err != nil {
		return err
	}
	return writeFileAtomic(path, b)
}

// pushRecent đưa path lên đầu danh sách (không trùng, tối đa maxRecent).
func pushRecent(list []string, path string) []string {
	out := []string{path}
	for _, p := range list {
		if p != path && len(out) < maxRecent {
			out = append(out, p)
		}
	}
	return out
}

func removeRecent(list []string, path string) []string {
	out := []string{}
	for _, p := range list {
		if p != path {
			out = append(out, p)
		}
	}
	return out
}

// editRecent sửa danh sách mở gần đây trong settings.json; lỗi chỉ ghi log – không chặn lưu / mở.
func editRecent(f func([]string) []string) {
	settingsMu.Lock()
	defer settingsMu.Unlock()
	dir, err := configDir()
	if err != nil {
		return
	}
	s := loadSettings(dir)
	s.Recent = f(s.Recent)
	if err := saveSettings(dir, s); err != nil {
		fmt.Println("không lưu được danh sách mở gần đây:", err)
	}
}

// ---- Gọi từ frontend ----

var scriptFilters = []runtime.FileFilter{{DisplayName: dialogName + " (*.json)", Pattern: "*" + scriptExt}}

// SaveScript ghi kịch bản vào path; path "" (chưa lưu lần nào) thì hỏi chỗ lưu như SaveScriptAs.
// Trả kịch bản đã ghi + đường dẫn; Path "" nếu người dùng huỷ hộp thoại.
func (a *App) SaveScript(sc Script, path string) (ScriptFile, error) {
	if path == "" {
		return a.SaveScriptAs(sc, "")
	}
	sc.Version = ScriptVersion
	if err := writeScriptFile(path, sc); err != nil {
		return ScriptFile{}, fmt.Errorf("không lưu được %s: %w", filepath.Base(path), err)
	}
	editRecent(func(l []string) []string { return pushRecent(l, path) })
	return ScriptFile{Script: sc, Path: path}, nil
}

// SaveScriptAs hỏi chỗ lưu (gợi ý thư mục / tên từ path hoặc tên kịch bản) rồi ghi; tên kịch bản đổi theo tên file.
func (a *App) SaveScriptAs(sc Script, path string) (ScriptFile, error) {
	opts := runtime.SaveDialogOptions{
		Title:                "Lưu kịch bản",
		DefaultFilename:      safeFileName(sc.Name) + scriptExt,
		Filters:              scriptFilters,
		CanCreateDirectories: true,
	}
	if path != "" {
		opts.DefaultDirectory = filepath.Dir(path)
		opts.DefaultFilename = filepath.Base(path)
	}
	dst, err := runtime.SaveFileDialog(a.ctx, opts)
	if err != nil || dst == "" {
		return ScriptFile{}, err
	}
	if filepath.Ext(dst) == "" {
		dst += scriptExt
	}
	sc.Name = strings.TrimSuffix(filepath.Base(dst), filepath.Ext(dst))
	return a.SaveScript(sc, dst)
}

// OpenScript hỏi file rồi đọc. Path "" nếu người dùng huỷ.
func (a *App) OpenScript() (ScriptFile, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "Mở kịch bản", Filters: scriptFilters})
	if err != nil || path == "" {
		return ScriptFile{}, err
	}
	return a.OpenScriptPath(path)
}

// OpenScriptPath đọc file kịch bản (mục "Mở gần đây"); file không còn thì bỏ khỏi danh sách.
func (a *App) OpenScriptPath(path string) (ScriptFile, error) {
	sc, err := readScriptFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			editRecent(func(l []string) []string { return removeRecent(l, path) })
			return ScriptFile{}, fmt.Errorf("không tìm thấy file %s", path)
		}
		return ScriptFile{}, fmt.Errorf("không mở được %s: %w", filepath.Base(path), err)
	}
	editRecent(func(l []string) []string { return pushRecent(l, path) })
	return ScriptFile{Script: sc, Path: path}, nil
}

// RecentScripts trả danh sách file mở / lưu gần đây, mới nhất trước.
func (a *App) RecentScripts() []string {
	dir, err := configDir()
	if err != nil {
		return []string{}
	}
	if r := loadSettings(dir).Recent; r != nil {
		return r
	}
	return []string{}
}

// SaveDraft ghi bản nháp tự lưu (frontend gọi sau mỗi lần sửa, có giãn cách).
func (a *App) SaveDraft(d Draft) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	d.SavedAt = time.Now().UnixMilli()
	b, err := encodeDraft(d)
	if err != nil {
		return err
	}
	return writeFileAtomic(filepath.Join(dir, draftFile), b)
}

func encodeDraft(d Draft) ([]byte, error) {
	d.Script.Version = ScriptVersion
	return json.Marshal(d)
}

// LoadDraft trả bản nháp lần trước; nil nếu chưa có hoặc hỏng (bỏ qua, mở kịch bản trống).
func (a *App) LoadDraft() *Draft {
	dir, err := configDir()
	if err != nil {
		return nil
	}
	return loadDraft(dir)
}

func loadDraft(dir string) *Draft {
	b, err := os.ReadFile(filepath.Join(dir, draftFile))
	if err != nil {
		return nil
	}
	var raw struct {
		Draft
		Script json.RawMessage `json:"script"`
	}
	if json.Unmarshal(b, &raw) != nil {
		return nil
	}
	sc, err := decodeScript(raw.Script)
	if err != nil {
		fmt.Println("bỏ qua bản nháp:", err)
		return nil
	}
	d := raw.Draft
	d.Script = sc
	return &d
}

// ClearDraft xoá bản nháp (đóng app chọn "Không lưu").
func (a *App) ClearDraft() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(dir, draftFile)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// ---- Hỏi lưu khi đóng app ----

// SetDirty: frontend báo kịch bản đang soạn có thay đổi chưa lưu – đóng app khi đó thì hỏi trước.
func (a *App) SetDirty(dirty bool) {
	a.mu.Lock()
	a.dirty = dirty
	a.mu.Unlock()
}

// ConfirmQuit đóng app sau khi frontend đã hỏi (Lưu / Không lưu).
func (a *App) ConfirmQuit() {
	a.mu.Lock()
	a.quitting = true
	a.mu.Unlock()
	runtime.Quit(a.ctx)
}

// beforeClose (OnBeforeClose): còn thay đổi chưa lưu → giữ app, nhờ frontend hỏi (event app:confirmClose).
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	a.mu.Lock()
	ask := a.dirty && !a.quitting
	a.mu.Unlock()
	if ask {
		runtime.EventsEmit(ctx, "app:confirmClose")
	}
	return ask
}

// safeFileName bỏ ký tự không dùng được trong tên file (Windows khắt khe nhất).
func safeFileName(name string) string {
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) || r < 32 {
			return '-'
		}
		return r
	}, strings.TrimSpace(name))
	if name == "" {
		return "kich-ban"
	}
	return name
}
