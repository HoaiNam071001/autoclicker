package engine

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script/scripttest"
)

func fileStep(f func(*model.Step)) model.Step { return modeltest.Step(model.StepFile, f) }

// runFile chạy các node, trả kết quả theo mã node.
func runFile(t *testing.T, trash *devicetest.FakeTrash, allow bool, steps ...model.Step) map[string]model.StepResult {
	t.Helper()
	for i := range steps {
		steps[i].Ref = "n" + string(rune('1'+i))
	}
	e := New(&devicetest.FakeInput{}, scripttest.Linear(model.Script{Loops: 1, Steps: steps}), Options{AllowDanger: allow})
	e.trash = trash
	res := map[string]model.StepResult{}
	e.onNode = func(r model.StepResult) { res[r.Ref] = r }
	e.Run(context.Background())
	return res
}

func TestFileWriteMkdirRows(t *testing.T) {
	dir := t.TempDir()
	txt, csvP, xlsx := filepath.Join(dir, "a", "ghi.txt"), filepath.Join(dir, "bang.csv"), filepath.Join(dir, "bang.xlsx")
	res := runFile(t, &devicetest.FakeTrash{}, false,
		fileStep(func(s *model.Step) { s.FileOp, s.FilePath = model.FileMkdir, filepath.Join(dir, "x", "y") }),
		fileStep(func(s *model.Step) { s.FilePath, s.Text = txt, "dòng 1\n" }),
		fileStep(func(s *model.Step) { s.FilePath, s.Text, s.FileAppend = txt, "dòng 2\n", true }),
		fileStep(func(s *model.Step) { s.FileOp, s.FilePath, s.Cells = model.FileRow, csvP, []string{"Tên", "Ghi chú"} }),
		fileStep(func(s *model.Step) {
			s.FileOp, s.FilePath, s.Cells = model.FileRow, csvP, []string{"Lan", `có "nháy", phẩy`}
		}),
		fileStep(func(s *model.Step) { s.FileOp, s.FilePath, s.Cells = model.FileRow, xlsx, []string{"Lan", "42"} }),
		fileStep(func(s *model.Step) {
			s.FileOp, s.FilePath, s.Cells, s.Sheet = model.FileRow, xlsx, []string{"Hà", "3.5"}, "Đơn"
		}),
	)
	for ref, r := range res {
		if r.Status != model.StatusSuccess {
			t.Fatalf("%s: %+v", ref, r.Error)
		}
	}
	if st, err := os.Stat(filepath.Join(dir, "x", "y")); err != nil || !st.IsDir() {
		t.Error("không tạo thư mục")
	}
	if b, _ := os.ReadFile(txt); string(b) != "dòng 1\ndòng 2\n" {
		t.Errorf("ghi + thêm vào cuối: %q", b)
	}
	if b, _ := os.ReadFile(csvP); string(b) != utf8BOM+"Tên,Ghi chú\nLan,\"có \"\"nháy\"\", phẩy\"\n" || res["n5"].Output["rows"] != 2 {
		t.Errorf("csv: %q, rows %v", b, res["n5"].Output["rows"])
	}
	f, err := excelize.OpenFile(xlsx)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if v, _ := f.GetCellValue("Sheet1", "B1"); v != "42" {
		t.Errorf("xlsx B1 = %q", v)
	}
	if typ, _ := f.GetCellType("Sheet1", "B1"); typ != excelize.CellTypeNumber && typ != excelize.CellTypeUnset {
		t.Errorf("ô số phải ghi thành số, có kiểu %v", typ)
	}
	if v, _ := f.GetCellValue("Đơn", "A1"); v != "Hà" || res["n7"].Output["rows"] != 1 {
		t.Errorf("sheet mới: %q, rows %v", v, res["n7"].Output["rows"])
	}
}

func TestFileTrash(t *testing.T) {
	dir := t.TempDir()
	for _, n := range []string{"a.png", "b.PNG", "c.png", ".an.png", "d.txt"} {
		os.WriteFile(filepath.Join(dir, n), []byte("x"), 0o644)
	}
	os.Mkdir(filepath.Join(dir, "thu-muc.png"), 0o755)
	trash := &devicetest.FakeTrash{}
	step := fileStep(func(s *model.Step) { s.FileOp, s.FilePath, s.FilePattern = model.FileTrash, dir, "*.png" })
	if r := runFile(t, trash, false, step)["n1"]; r.Error == nil || r.Error.Code != model.ErrPermission || len(trash.Trashed) > 0 {
		t.Fatalf("chưa xác nhận mà vẫn xoá: %+v %v", r.Error, trash.Trashed)
	}
	r := runFile(t, trash, true, step, fileStep(func(s *model.Step) { s.FileOp = model.FileEmptyTrash }))
	want := []string{filepath.Join(dir, "a.png"), filepath.Join(dir, "c.png")} // b.PNG: khác hoa thường; file ẩn, thư mục bỏ qua
	if strings.Join(trash.Trashed, "|") != strings.Join(want, "|") || r["n1"].Output["count"] != 2 || trash.Emptied != 1 {
		t.Errorf("xoá %v (count %v), dọn %d lần", trash.Trashed, r["n1"].Output["count"], trash.Emptied)
	}
	// không file nào khớp → vẫn thành công
	step.FilePattern = "*.khong"
	if r := runFile(t, &devicetest.FakeTrash{}, true, step)["n1"]; r.Status != model.StatusSuccess || r.Output["count"] != 0 {
		t.Errorf("không khớp: %+v", r)
	}
}

func TestGuardPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	blocked := []string{"/", home, filepath.Join(home, "Library", "x"), "/System/x", "/usr/local/bin/a.txt", "/etc/hosts"}
	if runtime.GOOS == "windows" {
		blocked = []string{`C:\`, home, filepath.Join(os.Getenv("SystemRoot"), "System32", "x.dll"),
			filepath.Join(os.Getenv("ProgramFiles"), "a.txt"), filepath.Join(os.Getenv("ProgramData"), "x")}
	}
	for _, p := range blocked {
		if guardPath(p) == nil {
			t.Errorf("%s phải bị chặn", p)
		}
	}
	for _, p := range []string{filepath.Join(home, "Documents", "a.txt"), filepath.Join(home, "Desktop"), filepath.Join(t.TempDir(), "x")} {
		if err := guardPath(p); err != nil {
			t.Errorf("%s bị chặn nhầm: %v", p, err)
		}
	}
	// chặn cả khi đi vòng qua symlink
	link := filepath.Join(t.TempDir(), "lib")
	if os.Symlink(filepath.Join(home, "Library"), link) == nil && guardPath(filepath.Join(link, "x")) == nil {
		t.Error("symlink tới ~/Library phải bị chặn")
	}
}
