package engine

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"autoclicker/internal/model"
)

// Node Tệp / thư mục (bước 35 N4): ghi file chữ, thêm dòng csv / xlsx, tạo thư mục, xoá file vào Thùng rác, dọn Thùng rác.
// Xoá chỉ chạy khi người dùng đã xác nhận (Options.AllowDelete – script.DeleteOps); thư mục hệ thống / thư mục người
// dùng bị chặn (guardPath).

func (e *Engine) fileStep(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	if (s.FileOp == model.FileTrash || s.FileOp == model.FileEmptyTrash) && !e.allowDelete {
		return nil, fmt.Errorf("%w: thao tác xoá chưa được xác nhận – bấm Chạy trong app để xác nhận", errDeleteNotAllowed)
	}
	if s.FileOp == model.FileEmptyTrash {
		return map[string]any{}, e.trash.EmptyTrash()
	}
	p, err := userPath(s.FilePath)
	if err != nil {
		return nil, err
	}
	if err := guardPath(p); err != nil {
		return nil, err
	}
	switch s.FileOp {
	case model.FileMkdir:
		if err := os.MkdirAll(p, 0o755); err != nil {
			return nil, fmt.Errorf("không tạo được thư mục: %w", err)
		}
		return map[string]any{"path": p}, nil
	case model.FileRow:
		if strings.EqualFold(filepath.Ext(p), ".xlsx") {
			return appendXLSX(p, s.Sheet, s.Cells)
		}
		if !strings.EqualFold(filepath.Ext(p), ".csv") {
			return nil, &invalidErr{fmt.Sprintf("thêm dòng chỉ cho file .csv hoặc .xlsx (file này: %q)", filepath.Base(p))}
		}
		return appendCSV(p, s.Cells)
	case model.FileTrash:
		return e.trashFiles(p, strings.TrimSpace(s.FilePattern))
	}
	return writeText(p, s.Text, s.FileAppend)
}

// utf8BOM – đầu file CSV mới: Excel đọc đúng UTF-8 (tiếng Việt).
const utf8BOM = "\xEF\xBB\xBF"

// errDeleteNotAllowed – xoá khi chưa xác nhận: lỗi code permission (flowrun errCode).
var errDeleteNotAllowed = fmt.Errorf("chưa cho phép xoá")

func writeText(p, text string, appendTo bool) (map[string]any, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, fmt.Errorf("không tạo được thư mục chứa file: %w", err)
	}
	if appendTo {
		f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("không mở được file: %w", err)
		}
		_, err = f.WriteString(text)
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return nil, fmt.Errorf("không ghi được file: %w", err)
		}
	} else if err := writeFileAtomic(p, []byte(text)); err != nil {
		return nil, err
	}
	return map[string]any{"path": p, "bytes": len(text)}, nil
}

// appendCSV thêm một dòng (tự đặt nháy khi ô có dấu phẩy / xuống dòng). File mới có BOM UTF-8 – Excel mới hiện đúng tiếng Việt.
func appendCSV(p string, cells []string) (map[string]any, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, fmt.Errorf("không tạo được thư mục chứa file: %w", err)
	}
	st, statErr := os.Stat(p)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("không mở được file: %w", err)
	}
	if statErr != nil || st.Size() == 0 {
		f.WriteString(utf8BOM)
	}
	w := csv.NewWriter(f)
	w.Write(cells)
	w.Flush()
	err = w.Error()
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return nil, fmt.Errorf("không ghi được file: %w", err)
	}
	b, _ := os.ReadFile(p)
	rows, _ := csv.NewReader(strings.NewReader(strings.TrimPrefix(string(b), utf8BOM))).ReadAll()
	return map[string]any{"path": p, "rows": len(rows)}, nil
}

// appendXLSX thêm dòng sau dòng cuối của sheet (chưa có file / sheet thì tạo). Ô trông như số ghi thành số (Excel cộng được).
func appendXLSX(p, sheet string, cells []string) (map[string]any, error) {
	f, err := excelize.OpenFile(p)
	if os.IsNotExist(err) {
		f, err = excelize.NewFile(), nil
		if sheet != "" {
			f.SetSheetName(f.GetSheetName(0), sheet)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("không mở được file Excel: %w", err)
	}
	defer f.Close()
	if sheet == "" {
		sheet = f.GetSheetName(0)
	} else if idx, _ := f.GetSheetIndex(sheet); idx < 0 {
		if _, err := f.NewSheet(sheet); err != nil {
			return nil, &invalidErr{fmt.Sprintf("tên sheet “%s” không dùng được: %v", sheet, err)}
		}
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("không đọc được sheet %s: %w", sheet, err)
	}
	vals := make([]any, len(cells))
	for i, c := range cells {
		if n, err := strconv.ParseFloat(strings.TrimSpace(c), 64); err == nil && strings.TrimSpace(c) != "" {
			vals[i] = n
		} else {
			vals[i] = c
		}
	}
	row := len(rows) + 1
	if err := f.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &vals); err != nil {
		return nil, fmt.Errorf("không ghi được dòng: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, fmt.Errorf("không tạo được thư mục chứa file: %w", err)
	}
	tmp := filepath.Join(filepath.Dir(p), ".tmp-"+filepath.Base(p))
	if err := f.SaveAs(tmp); err != nil {
		os.Remove(tmp)
		return nil, fmt.Errorf("không lưu được file Excel (đang mở trong Excel?): %w", err)
	}
	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return nil, fmt.Errorf("không lưu được file Excel (đang mở trong Excel?): %w", err)
	}
	return map[string]any{"path": p, "rows": row}, nil
}

// trashFiles chuyển file (không phải thư mục) khớp mẫu trong thư mục dir vào Thùng rác – không đệ quy; file ẩn chỉ khớp
// khi mẫu bắt đầu bằng ".". Không có file nào khớp vẫn thành công (count 0).
func (e *Engine) trashFiles(dir, pattern string) (map[string]any, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, &invalidErr{fmt.Sprintf("không đọc được thư mục %s: %v", dir, err)}
	}
	var paths, names []string
	for _, en := range entries {
		n := en.Name()
		if !en.Type().IsRegular() || (strings.HasPrefix(n, ".") && !strings.HasPrefix(pattern, ".")) {
			continue
		}
		if ok, _ := filepath.Match(pattern, n); ok {
			paths, names = append(paths, filepath.Join(dir, n)), append(names, n)
		}
	}
	if err := e.trash.MoveToTrash(paths); err != nil {
		return nil, err
	}
	if len(names) > 20 {
		names = append(names[:20], fmt.Sprintf("… (+%d)", len(names)-20))
	}
	return map[string]any{"count": len(paths), "files": strings.Join(names, ", "), "path": dir}, nil
}

// writeFileAtomic ghi file tạm rồi đổi tên – lỗi giữa chừng không để lại file hỏng.
func writeFileAtomic(p string, b []byte) error {
	f, err := os.CreateTemp(filepath.Dir(p), ".tmp-*")
	if err != nil {
		return fmt.Errorf("không ghi được file: %w", err)
	}
	_, err = f.Write(b)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), p)
	}
	if err != nil {
		os.Remove(f.Name())
		return fmt.Errorf("không ghi được file: %w", err)
	}
	return nil
}

// guardPath chặn thư mục gốc, thư mục người dùng (chính nó), thư mục hệ thống và ~/Library – ghi / xoá nhầm ở đó hỏng
// máy. Thư mục tạm của hệ điều hành được phép. Theo symlink trước khi xét.
func guardPath(p string) error {
	abs := realPath(filepath.Clean(p))
	home, _ := os.UserHomeDir()
	home = realPath(home)
	if tmp := realPath(os.TempDir()); within(abs, tmp) {
		return nil
	}
	bad := func() error {
		return &invalidErr{fmt.Sprintf("không thao tác trong %s – thư mục hệ thống / thư mục người dùng được bảo vệ", p)}
	}
	if abs == filepath.VolumeName(abs)+string(filepath.Separator) || abs == home || within(abs, filepath.Join(home, "Library")) {
		return bad()
	}
	for _, sys := range systemDirs() {
		if within(abs, realPath(sys)) {
			return bad()
		}
	}
	return nil
}

func systemDirs() []string {
	if runtime.GOOS == "windows" {
		var out []string
		for _, env := range []string{"SystemRoot", "ProgramFiles", "ProgramFiles(x86)", "ProgramData"} {
			if v := os.Getenv(env); v != "" {
				out = append(out, v)
			}
		}
		return out
	}
	return []string{"/System", "/Library", "/Applications", "/usr", "/bin", "/sbin", "/etc", "/private", "/var", "/opt", "/dev", "/cores"}
}

// realPath – đường dẫn sau khi theo symlink; phần chưa tồn tại giữ nguyên (theo symlink của thư mục cha gần nhất có thật).
func realPath(p string) string {
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	parent := filepath.Dir(p)
	if parent == p {
		return p
	}
	return filepath.Join(realPath(parent), filepath.Base(p))
}

// within: p là dir hoặc nằm trong dir (không phân biệt hoa thường trên Windows / macOS mặc định).
func within(p, dir string) bool {
	if dir == "" {
		return false
	}
	rel, err := filepath.Rel(strings.ToLower(dir), strings.ToLower(p))
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
