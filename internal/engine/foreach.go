package engine

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/unicode/norm"

	"ghosthand/internal/model"
)

// Node Lặp qua danh sách (bước 42 G2): lần đầu đi vào thì nạp danh sách, mỗi lần đi vào (cạnh quay lại từ cuối thân
// vòng) lấy một phần tử ra cổng item; hết thì ra done và quên danh sách – đi vào lần nữa (vd vòng ngoài) nạp lại từ đầu.

// listCursor – danh sách đã nạp của một node forEach và phần tử sẽ lấy tiếp.
type listCursor struct {
	items []any
	next  int
}

func (e *Engine) forEachStep(s model.Step, logf func(level, format string, a ...any)) (string, map[string]any, error) {
	cur := e.lists[s.ID]
	if cur == nil {
		items, _, err := loadList(s)
		if err != nil {
			return "", nil, err
		}
		cur = &listCursor{items: items}
		e.lists[s.ID] = cur
		logf("info", "danh sách có %d phần tử", len(items))
	}
	n := len(cur.items)
	if cur.next >= n {
		delete(e.lists, s.ID)
		logf("info", "hết danh sách (%d phần tử)", n)
		return model.PortDone, map[string]any{"count": n}, nil
	}
	item := cur.items[cur.next]
	cur.next++
	return model.PortItem, map[string]any{"item": item, "index": cur.next, "count": n}, nil
}

// ListPreview – đọc thử danh sách cho inspector (App.PreviewList): số phần tử, cột của bảng (tên gốc → khoá dùng
// trong biến), phần tử đầu.
type ListPreview struct {
	Count   int          `json:"count"`
	Columns []ListColumn `json:"columns"`
	First   any          `json:"first,omitempty"`
}

type ListColumn struct {
	Name string `json:"name"` // tên cột trong file
	Key  string `json:"key"`  // khoá: {{nX.output.item.<key>}}
}

// PreviewList đọc danh sách của node forEach như lúc chạy (chưa thay biến – có biến thì báo chỉ biết khi chạy).
func PreviewList(s model.Step) (ListPreview, error) {
	src := s.Value
	if s.ListFrom == model.ListFile {
		src = s.FilePath + s.Sheet
	}
	if strings.Contains(src, "{{") { // biến node / input chỉ có giá trị lúc chạy
		return ListPreview{}, errors.New("danh sách có biến – chỉ biết khi chạy")
	}
	items, cols, err := loadList(s)
	if err != nil {
		return ListPreview{}, err
	}
	pv := ListPreview{Count: len(items), Columns: cols}
	if len(items) > 0 {
		pv.First = items[0]
	}
	return pv, nil
}

// loadList – phần tử của danh sách theo ListFrom (Value đã thay biến); cols = cột của bảng (nguồn file).
func loadList(s model.Step) ([]any, []ListColumn, error) {
	if s.ListFrom != model.ListFile {
		items, err := parseList(s.Value)
		return items, nil, err
	}
	p, err := userPath(s.FilePath)
	if err != nil {
		return nil, nil, err
	}
	var rows [][]string
	switch strings.ToLower(filepath.Ext(p)) {
	case ".csv":
		rows, err = readCSV(p)
	case ".xlsx":
		rows, err = readXLSX(p, s.Sheet)
	default:
		return nil, nil, &invalidErr{fmt.Sprintf("chỉ đọc được file .csv hoặc .xlsx (file này: %q)", filepath.Base(p))}
	}
	if err != nil {
		return nil, nil, err
	}
	items, cols := tableItems(rows)
	return items, cols, nil
}

// parseList: mảng JSON (biến trỏ output mảng được in ra JSON) hoặc chữ – mỗi dòng không trống một phần tử.
func parseList(v string) ([]any, error) {
	t := strings.TrimSpace(v)
	if strings.HasPrefix(t, "[") {
		var arr []any
		if err := json.Unmarshal([]byte(t), &arr); err != nil {
			return nil, &invalidErr{fmt.Sprintf("danh sách bắt đầu bằng [ nhưng không phải mảng JSON hợp lệ: %v", err)}
		}
		return arr, nil
	}
	if strings.HasPrefix(t, "{") {
		return nil, &invalidErr{"danh sách là một object JSON, cần mảng – trỏ biến vào trường mảng bên trong (vd {{n3.output.json.items}})"}
	}
	items := []any{}
	for _, line := range strings.Split(t, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			items = append(items, line)
		}
	}
	return items, nil
}

// readCSV – các dòng của file CSV: bỏ BOM UTF-8, tự nhận dấu ngăn cột (, ; hoặc tab – Excel tiếng Việt hay lưu bằng ;).
func readCSV(p string) ([][]string, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, &invalidErr{fmt.Sprintf("không đọc được file: %v", err)}
	}
	text := strings.TrimPrefix(string(b), utf8BOM)
	r := csv.NewReader(strings.NewReader(text))
	r.Comma, r.FieldsPerRecord, r.LazyQuotes = csvComma(text), -1, true
	rows, err := r.ReadAll()
	if err != nil {
		return nil, &invalidErr{fmt.Sprintf("file CSV hỏng: %v", err)}
	}
	return rows, nil
}

// csvComma – dấu ngăn cột xuất hiện nhiều nhất ở dòng đầu (mặc định phẩy).
func csvComma(text string) rune {
	head, _, _ := strings.Cut(text, "\n")
	best, n := ',', strings.Count(head, ",")
	for _, c := range []rune{';', '\t'} {
		if k := strings.Count(head, string(c)); k > n {
			best, n = c, k
		}
	}
	return best
}

func readXLSX(p, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(p)
	if err != nil {
		return nil, &invalidErr{fmt.Sprintf("không mở được file Excel: %v", err)}
	}
	defer f.Close()
	if sheet == "" {
		sheet = f.GetSheetName(0)
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, &invalidErr{fmt.Sprintf("không đọc được sheet “%s”: %v", sheet, err)}
	}
	return rows, nil
}

// tableItems – dòng đầu là tên cột, mỗi dòng sau (bỏ dòng trống) thành {khoá cột: giá trị}. Khoá = colKey của tên
// (biến chỉ nhận chữ không dấu); cột không tên → cot<số thứ tự từ 1>, trùng khoá → thêm _<số thứ tự>.
func tableItems(rows [][]string) ([]any, []ListColumn) {
	items := []any{}
	if len(rows) == 0 {
		return items, nil
	}
	cols := make([]ListColumn, len(rows[0]))
	seen := map[string]bool{}
	for i, h := range rows[0] {
		h = strings.TrimSpace(h)
		k := colKey(h)
		if k == "" {
			k = fmt.Sprintf("cot%d", i+1)
		} else if seen[k] {
			k = fmt.Sprintf("%s_%d", k, i+1)
		}
		cols[i], seen[k] = ListColumn{Name: h, Key: k}, true
	}
	for _, row := range rows[1:] {
		if strings.TrimSpace(strings.Join(row, "")) == "" {
			continue
		}
		item := map[string]any{}
		for i, c := range cols {
			v := ""
			if i < len(row) {
				v = strings.TrimSpace(row[i])
			}
			item[c.Key] = v
		}
		items = append(items, item)
	}
	return items, cols
}

// colKey – tên cột thành khoá dùng được trong biến: bỏ dấu, chữ thường, cụm ký tự khác chữ / số → _
// ("Số điện thoại" → so_dien_thoai).
func colKey(name string) string {
	var b strings.Builder
	gap := false
	for _, r := range norm.NFD.String(name) {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		r = unicode.ToLower(r)
		if r == 'đ' {
			r = 'd'
		}
		if r < 128 && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
			if gap && b.Len() > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r)
			gap = false
		} else {
			gap = true
		}
	}
	return b.String()
}
