package engine

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"ghosthand/internal/element"
	"ghosthand/internal/model"
)

// Node Đọc bảng (bước 43 B5): tìm bảng (table / grid) trong cửa sổ app – hoặc trong phần tử selectors khớp – rồi đọc
// từng dòng / ô qua cây accessibility. Dòng đầu là tên cột. FilePath khác trống thì ghi luôn ra .csv / .xlsx (ghi đè).

func (e *Engine) readTable(ctx context.Context, s model.Step) (map[string]any, error) {
	q, scoped, err := pageScope(s)
	if err != nil {
		return nil, err
	}
	var rows [][]string
	why := ""
	var miss func() map[string]any
	_, err = e.waitLoop(ctx, s, "element", func() (bool, int, error) {
		nodes, err := e.elems.Elements(s.ElementApp)
		if waitable(err) {
			why = err.Error()
			return false, 0, nil
		}
		if err != nil {
			return false, 0, err
		}
		from, to, nth := 0, len(nodes), s.Occurrence
		if scoped {
			r := element.Find(nodes, q)
			if !r.Found {
				why, miss = missWhy(r, nodes, max(s.Occurrence, 1)), func() map[string]any { return elementMiss(s.Selectors, nodes, r, q) }
				return false, 0, nil
			}
			from, to, nth = r.Hit.Index, element.SubtreeEnd(nodes, r.Hit.Index), 1
		}
		t := element.FindTable(nodes, from, to, nth, s.PageAll)
		if t < 0 {
			why = fmt.Sprintf("đã đọc %d phần tử nhưng không thấy bảng (vai trò table / grid)", to-from)
			return false, 0, nil
		}
		if rows = element.TableRows(nodes, t, s.PageAll); len(rows) == 0 {
			why = "bảng chưa có dòng nào (trang còn tải?)"
			return false, 0, nil
		}
		return true, 100, nil
	})
	if errors.Is(err, errTimeout) {
		return nil, &notFoundErr{msg: fmt.Sprintf("không đọc được bảng sau %s – %s", fmtSeconds(s.TimeoutMs), why), detail: lazyDetail(miss)}
	}
	if err != nil {
		return nil, err
	}
	items, cols := tableItems(rows)
	out := map[string]any{"rows": rowsAny(rows), "items": items, "columns": colNames(cols), "count": len(items)}
	if p := strings.TrimSpace(s.FilePath); p != "" {
		path, err := saveTable(p, s.Sheet, rows)
		if err != nil {
			return nil, err
		}
		out["path"] = path
	}
	return out, nil
}

func rowsAny(rows [][]string) []any {
	out := make([]any, len(rows))
	for i, r := range rows {
		cells := make([]any, len(r))
		for j, c := range r {
			cells[j] = c
		}
		out[i] = cells
	}
	return out
}

func colNames(cols []ListColumn) []any {
	out := make([]any, len(cols))
	for i, c := range cols {
		out[i] = c.Name
	}
	return out
}

// saveTable – ghi đè p (.csv có BOM UTF-8 cho Excel | .xlsx, ô trông như số ghi thành số).
func saveTable(p, sheet string, rows [][]string) (string, error) {
	p, err := userPath(p)
	if err != nil {
		return "", err
	}
	if err := guardPath(p); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", fmt.Errorf("không tạo được thư mục chứa file: %w", err)
	}
	switch strings.ToLower(filepath.Ext(p)) {
	case ".csv":
		var b strings.Builder
		b.WriteString(utf8BOM)
		if err := csv.NewWriter(&b).WriteAll(rows); err != nil {
			return "", fmt.Errorf("không ghi được file: %w", err)
		}
		return p, writeFileAtomic(p, []byte(b.String()))
	case ".xlsx":
		f := excelize.NewFile()
		defer f.Close()
		if sheet != "" {
			f.SetSheetName(f.GetSheetName(0), sheet)
		} else {
			sheet = f.GetSheetName(0)
		}
		for i, r := range rows {
			vals := make([]any, len(r))
			for j, c := range r {
				if n, err := strconv.ParseFloat(strings.TrimSpace(c), 64); err == nil && strings.TrimSpace(c) != "" {
					vals[j] = n
				} else {
					vals[j] = c
				}
			}
			cell, _ := excelize.CoordinatesToCellName(1, i+1)
			if err := f.SetSheetRow(sheet, cell, &vals); err != nil {
				return "", fmt.Errorf("không ghi được dòng %d: %w", i+1, err)
			}
		}
		tmp := filepath.Join(filepath.Dir(p), ".tmp-"+filepath.Base(p))
		if err := f.SaveAs(tmp); err != nil {
			os.Remove(tmp)
			return "", fmt.Errorf("không lưu được file Excel (đang mở trong Excel?): %w", err)
		}
		if err := os.Rename(tmp, p); err != nil {
			os.Remove(tmp)
			return "", fmt.Errorf("không lưu được file Excel (đang mở trong Excel?): %w", err)
		}
		return p, nil
	}
	return "", &invalidErr{fmt.Sprintf("chỉ ghi được file .csv hoặc .xlsx (file này: %q)", filepath.Base(p))}
}
