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
	"time"

	"autoclicker/internal/model"
	"autoclicker/internal/store"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Lưu / mở kịch bản (bước 21): file .json, bản nháp tự lưu và danh sách mở gần đây trong configDir().

const (
	draftFile  = "draft.json"
	scriptExt  = ".json"
	dialogName = "Kịch bản Auto Clicker"
)

// ScriptFile – kịch bản vừa mở và đường dẫn file của nó. Path "" = người dùng huỷ hộp thoại.
type ScriptFile struct {
	Script model.Script `json:"script"`
	Path   string       `json:"path"`
}

// Draft – bản nháp tự lưu: kịch bản đang soạn kèm file nó thuộc về, mở lại app thì khôi phục nguyên trạng.
type Draft struct {
	Script  model.Script `json:"script"`
	Path    string       `json:"path"`    // file đang mở, "" = chưa lưu lần nào
	Dirty   bool         `json:"dirty"`   // có thay đổi chưa lưu xuống file
	SavedAt int64        `json:"savedAt"` // thời điểm ghi nháp, unix ms
}

func encodeScript(sc model.Script) ([]byte, error) {
	sc.Version = model.ScriptVersion
	return json.MarshalIndent(sc, "", "  ")
}

// decodeScript đọc file kịch bản: kiểm tra version, cấu trúc (checkScript) và bổ sung giá trị trống.
// Không chạy validateScript – bản đang soạn dở (vd bước gõ chữ chưa nhập) vẫn mở được, lúc Chạy mới kiểm tra.
func decodeScript(b []byte) (model.Script, error) {
	var sc model.Script
	if err := json.Unmarshal(b, &sc); err != nil {
		return model.Script{}, fmt.Errorf("file hỏng hoặc không phải JSON: %w", err)
	}
	switch {
	case sc.Version == 0:
		return model.Script{}, errors.New("không phải file kịch bản (thiếu version)")
	case sc.Version > model.ScriptVersion:
		return model.Script{}, fmt.Errorf("file tạo bởi phiên bản app mới hơn (version %d, app đọc được tới %d) – hãy cập nhật app", sc.Version, model.ScriptVersion)
	case sc.Version < model.ScriptVersion:
		return model.Script{}, fmt.Errorf("kịch bản tạo bằng bản cũ (version %d, dạng danh sách bước), không mở được – bản này dùng dạng flow (version %d)", sc.Version, model.ScriptVersion)
	}
	if err := checkScript(&sc); err != nil {
		return model.Script{}, err
	}
	return sc, nil
}

// checkScript kiểm tra giá trị nằm trong miền hợp lệ và bổ sung phần trống (slice nil, id / ref trùng, nút / cách bấm mặc định,
// cạnh trỏ node không còn) để frontend dùng ngay được. Không kiểm tra đủ flow – validateFlow làm lúc Chạy.
func checkScript(sc *model.Script) error {
	if sc.Loops < 0 || sc.MaxMinutes < 0 {
		return errors.New("số vòng / số phút không được âm")
	}
	if sc.Region != nil && (sc.Region.W <= 0 || sc.Region.H <= 0) {
		return errors.New("vùng có kích thước không hợp lệ")
	}
	if sc.AnchorMode != model.AnchorNearest && sc.AnchorMode != model.AnchorTopLeft {
		return fmt.Errorf("cách neo điểm không hợp lệ: %q", sc.AnchorMode)
	}
	if sc.Steps == nil {
		sc.Steps = []model.Step{}
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
	fixRefs(sc.Steps)
	checkEdges(sc, ids)
	return checkInputs(sc)
}

// fixRefs đặt ref cho node chưa có / trùng / sai dạng: n1, n2… (số nhỏ nhất chưa dùng).
func fixRefs(steps []model.Step) {
	used := map[string]bool{}
	bad := []int{}
	for i, s := range steps {
		if !model.ValidRef(s.Ref) || used[s.Ref] {
			bad = append(bad, i)
			continue
		}
		used[s.Ref] = true
	}
	n := 0
	for _, i := range bad {
		for n++; used[fmt.Sprintf("n%d", n)]; n++ {
		}
		steps[i].Ref = fmt.Sprintf("n%d", n)
		used[steps[i].Ref] = true
	}
}

// checkEdges bỏ cạnh hỏng (trỏ node không còn, thiếu cổng, trùng cổng – sửa tay file), đặt id cạnh trống / trùng;
// start trỏ node không còn thì bỏ trống.
func checkEdges(sc *model.Script, ids map[string]bool) {
	if sc.Edges == nil {
		sc.Edges = []model.Edge{}
	}
	out := sc.Edges[:0]
	seenPort := map[[2]string]bool{}
	seenID := map[string]bool{}
	n := 0
	for _, e := range sc.Edges {
		key := [2]string{e.From, e.Port}
		if !ids[e.From] || !ids[e.To] || e.Port == "" || seenPort[key] {
			continue
		}
		seenPort[key] = true
		for e.ID == "" || seenID[e.ID] {
			n++
			e.ID = fmt.Sprintf("e%d", n)
		}
		seenID[e.ID] = true
		out = append(out, e)
	}
	sc.Edges = out
	if !ids[sc.Start] {
		sc.Start = ""
	}
}

func checkStep(s *model.Step) error {
	need := 0
	switch s.Type {
	case model.StepClick, model.StepDouble, model.StepHold, model.StepWaitPixel:
		need = 1
	case model.StepDrag:
		need = 2
	case model.StepWait, model.StepText, model.StepKey, model.StepOpen, model.StepScroll, model.StepWaitImage, model.StepClickImage:
	case model.StepIf, model.StepSwitch, model.StepSwitchScreen, model.StepCounter, model.StepEnd, model.StepSetVar:
		if err := checkControl(s); err != nil {
			return err
		}
	default:
		return fmt.Errorf("loại bước không hợp lệ: %q", s.Type)
	}
	if s.Direction != "" && !slices.Contains(model.ScrollDirections, s.Direction) {
		return fmt.Errorf("hướng cuộn không hợp lệ: %q", s.Direction)
	}
	if s.ImageSource != "" && s.ImageSource != model.ImageFromNode && s.ImageSource != model.ImageFromSearch {
		return fmt.Errorf("nguồn vị trí ảnh không hợp lệ: %q", s.ImageSource)
	}
	if r := s.Retry; r != nil && (r.Count < 0 || r.Count > model.MaxRetries || r.DelayMs < 0) {
		return fmt.Errorf("số lần thử lại phải trong 0–%d, thời gian nghỉ không âm", model.MaxRetries)
	}
	if s.FindBy != "" && s.FindBy != model.FindByText {
		return fmt.Errorf("cách tìm không hợp lệ: %q", s.FindBy)
	}
	if s.SearchArea != "" && s.SearchArea != "screen" {
		return fmt.Errorf("vùng tìm không hợp lệ: %q", s.SearchArea)
	}
	if s.Threshold < 0 || s.Threshold > 100 || s.Tolerance < 0 || s.Tolerance > 255 || s.TemplateScale < 0 {
		return errors.New("độ khớp / sai số màu / tỷ lệ ảnh mẫu ngoài khoảng cho phép")
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
		s.KeyMode = model.KeyChord
	case model.KeyChord, model.KeySequence:
	default:
		return fmt.Errorf("cách bấm phím không hợp lệ: %q", s.KeyMode)
	}
	for _, m := range s.Modifiers {
		if m != "cmd" && m != "ctrl" && m != "alt" && m != "shift" {
			return fmt.Errorf("phím kèm không hợp lệ: %q", m)
		}
	}
	if s.Clicks < 0 || s.ClickInterval < 0 || s.HoldMs < 0 || s.DragMs < 0 || s.DelayAfter < 0 || s.JitterMs < 0 || s.JitterPx < 0 || s.Amount < 0 || s.TimeoutMs < 0 || s.IntervalMs < 0 {
		return errors.New("số lần / thời gian / ngẫu nhiên không được âm")
	}
	if s.Points == nil {
		s.Points = []model.Point{}
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

func readScriptFile(path string) (model.Script, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return model.Script{}, err
	}
	sc, err := decodeScript(b)
	if err != nil {
		return model.Script{}, err
	}
	if sc.Name == "" {
		sc.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	return sc, nil
}

func writeScriptFile(path string, sc model.Script) error {
	b, err := encodeScript(sc)
	if err != nil {
		return err
	}
	return store.WriteFileAtomic(path, b)
}

// ---- Gọi từ frontend ----

var scriptFilters = []runtime.FileFilter{{DisplayName: dialogName + " (*.json)", Pattern: "*" + scriptExt}}

// SaveScript ghi kịch bản vào path; path "" (chưa lưu lần nào) thì hỏi chỗ lưu như SaveScriptAs.
// Trả kịch bản đã ghi + đường dẫn; Path "" nếu người dùng huỷ hộp thoại.
func (a *App) SaveScript(sc model.Script, path string) (ScriptFile, error) {
	if path == "" {
		return a.SaveScriptAs(sc, "")
	}
	sc.Version = model.ScriptVersion
	if err := writeScriptFile(path, sc); err != nil {
		return ScriptFile{}, fmt.Errorf("không lưu được %s: %w", filepath.Base(path), err)
	}
	store.EditRecent(func(l []string) []string { return store.PushRecent(l, path) })
	return ScriptFile{Script: sc, Path: path}, nil
}

// SaveScriptAs hỏi chỗ lưu (gợi ý thư mục / tên từ path hoặc tên kịch bản) rồi ghi; tên kịch bản đổi theo tên file.
func (a *App) SaveScriptAs(sc model.Script, path string) (ScriptFile, error) {
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
			store.EditRecent(func(l []string) []string { return store.RemoveRecent(l, path) })
			return ScriptFile{}, fmt.Errorf("không tìm thấy file %s", path)
		}
		return ScriptFile{}, fmt.Errorf("không mở được %s: %w", filepath.Base(path), err)
	}
	store.EditRecent(func(l []string) []string { return store.PushRecent(l, path) })
	return ScriptFile{Script: sc, Path: path}, nil
}

// RecentScripts trả danh sách file mở / lưu gần đây, mới nhất trước.
func (a *App) RecentScripts() []string { return store.RecentScripts() }

// SaveDraft ghi bản nháp tự lưu (frontend gọi sau mỗi lần sửa, có giãn cách).
func (a *App) SaveDraft(d Draft) error {
	dir, err := store.ConfigDir()
	if err != nil {
		return err
	}
	d.SavedAt = time.Now().UnixMilli()
	b, err := encodeDraft(d)
	if err != nil {
		return err
	}
	return store.WriteFileAtomic(filepath.Join(dir, draftFile), b)
}

func encodeDraft(d Draft) ([]byte, error) {
	d.Script.Version = model.ScriptVersion
	return json.Marshal(d)
}

// LoadDraft trả bản nháp lần trước; nil nếu chưa có hoặc hỏng (bỏ qua, mở kịch bản trống).
func (a *App) LoadDraft() *Draft {
	dir, err := store.ConfigDir()
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
	dir, err := store.ConfigDir()
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
