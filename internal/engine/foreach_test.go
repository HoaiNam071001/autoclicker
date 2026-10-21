package engine

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

// forEachScript: n1 Lặp qua danh sách → item: n2 gõ {{n1.output.item…}} → quay lại n1; done: n3 gõ "|xong".
func forEachScript(set func(*model.Step), text string) model.Script {
	return model.Script{Start: "f", Steps: []model.Step{
		modeltest.Node("f", "n1", model.StepForEach, set),
		modeltest.Node("t", "n2", model.StepText, func(s *model.Step) { s.Text = text }),
		modeltest.Node("d", "n3", model.StepText, func(s *model.Step) { s.Text = "|xong {{n1.output.count}}" }),
	}, Edges: []model.Edge{modeltest.Edge("f", model.PortItem, "t"), modeltest.Edge("t", model.PortSuccess, "f"), modeltest.Edge("f", model.PortDone, "d")}}
}

func TestForEachLinesAndJSON(t *testing.T) {
	for _, c := range []struct{ name, value, text, want string }{
		{"dòng", "táo\n\n  cam \r\nxoài", "{{n1.output.index}}={{n1.output.item}},", "1=táo,2=cam,3=xoài,|xong 3"},
		{"mảng JSON", `[{"ten":"An","tuoi":3},{"ten":"Bình","tuoi":5}]`, "{{n1.output.item.ten}}:{{n1.output.item.tuoi}}/{{n1.output.count}},", "An:3/2,Bình:5/2,|xong 2"},
		{"rỗng", "[]", "x", "|xong 0"},
	} {
		e, in := ctlEngine(forEachScript(func(s *model.Step) { s.Value = c.value }, c.text), nil)
		if err := e.Run(context.Background()); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if typed(in) != c.want {
			t.Errorf("%s: gõ %q, muốn %q", c.name, typed(in), c.want)
		}
	}
}

// Danh sách từ biến node trước (output mảng in ra JSON); object JSON / mảng hỏng → lỗi invalid_input.
func TestForEachFromVarAndErrors(t *testing.T) {
	sc := forEachScript(func(s *model.Step) { s.Value = "{{vars.ds}}" }, "{{n1.output.item}};")
	sc.Steps = append(sc.Steps, modeltest.Node("v", "n0", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "ds", model.VarSet, "a\nb" }))
	sc.Edges = append(sc.Edges, modeltest.Edge("v", model.PortSuccess, "f"))
	sc.Start = "v"
	e, in := ctlEngine(sc, nil)
	if err := e.Run(context.Background()); err != nil || typed(in) != "a;b;|xong 2" {
		t.Fatalf("lỗi %v, gõ %q", err, typed(in))
	}
	for _, v := range []string{`{"items":[1]}`, `[1,`} {
		e, _ := ctlEngine(forEachScript(func(s *model.Step) { s.Value = v }, "x"), nil)
		if err := e.Run(context.Background()); err == nil || e.results["f"].Error.Code != model.ErrInvalidInput {
			t.Errorf("%s: phải lỗi invalid_input, có %v", v, err)
		}
	}
}

// Hết danh sách rồi đi vào lần nữa (vòng ngoài) → nạp lại từ đầu; mỗi vòng kịch bản cũng nạp lại.
func TestForEachRestartsAfterDone(t *testing.T) {
	sc := forEachScript(func(s *model.Step) { s.Value = "a\nb" }, "{{n1.output.item}}")
	sc.Steps = append(sc.Steps, modeltest.Node("k", "n4", model.StepCounter, func(s *model.Step) { s.Max = 1 }))
	sc.Edges = append(sc.Edges[:2], modeltest.Edge("f", model.PortDone, "k"), modeltest.Edge("k", model.PortNext, "f"))
	sc.Loops = 2
	e, in := ctlEngine(sc, nil)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if typed(in) != "abababab" {
		t.Errorf("gõ %q, muốn abababab", typed(in))
	}
}

func TestForEachCSVAndXLSX(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "khach.csv")
	os.WriteFile(csvPath, []byte(utf8BOM+"Tên;Số điện thoại;\r\nNguyễn Văn A;0901\r\n;;\r\n\"Trần; B\";0902;thừa\r\n"), 0o644)
	xlsxPath := filepath.Join(dir, "khach.xlsx")
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]any{"Tên", "Tên", "Tuổi"})
	f.SetSheetRow("Sheet1", "A2", &[]any{"Lê C", "x", 30})
	if err := f.SaveAs(xlsxPath); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ path, text, want string }{
		{csvPath, "{{n1.output.item.ten}}/{{n1.output.item.so_dien_thoai}}/{{n1.output.item.cot3}},", "Nguyễn Văn A/0901/,Trần; B/0902/thừa,|xong 2"},
		{xlsxPath, "{{n1.output.item.ten}}-{{n1.output.item.ten_2}}-{{n1.output.item.tuoi}},", "Lê C-x-30,|xong 1"},
	} {
		e, in := ctlEngine(forEachScript(func(s *model.Step) { s.ListFrom, s.FilePath = model.ListFile, c.path }, c.text), nil)
		if err := e.Run(context.Background()); err != nil {
			t.Fatalf("%s: %v", filepath.Base(c.path), err)
		}
		if typed(in) != c.want {
			t.Errorf("%s: gõ %q, muốn %q", filepath.Base(c.path), typed(in), c.want)
		}
	}
	e, _ := ctlEngine(forEachScript(func(s *model.Step) { s.ListFrom, s.FilePath = model.ListFile, filepath.Join(dir, "khong-co.csv") }, "x"), nil)
	if err := e.Run(context.Background()); err == nil || !strings.Contains(err.Error(), "không đọc được file") {
		t.Errorf("file không có phải lỗi rõ ràng: %v", err)
	}
}

// Thanh mini: "lượt n1 i / n" khi đang ở thân vòng, hết danh sách thì bỏ.
func TestForEachProgress(t *testing.T) {
	e, _ := ctlEngine(forEachScript(func(s *model.Step) { s.Value = "a\nb" }, "x"), nil)
	var got []string
	e.onProgress = func(p Progress) {
		c := "-"
		if p.Counter != nil {
			c = p.Counter.Ref + " " + string(rune('0'+p.Counter.Count)) + "/" + string(rune('0'+p.Counter.Max))
		}
		if p.Wait == nil && (len(got) == 0 || got[len(got)-1] != p.Ref+" "+c) {
			got = append(got, p.Ref+" "+c)
		}
	}
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := []string{"n1 -", "n2 n1 1/2", "n1 n1 1/2", "n2 n1 2/2", "n1 n1 2/2", "n3 -"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tiến độ %q, muốn %q", got, want)
	}
}

func TestColKeyAndPreview(t *testing.T) {
	for in, want := range map[string]string{"Số điện thoại": "so_dien_thoai", "  Đơn giá (VNĐ) ": "don_gia_vnd", "email": "email", "Họ & tên": "ho_ten", "日本": ""} {
		if got := colKey(in); got != want {
			t.Errorf("colKey(%q) = %q, muốn %q", in, got, want)
		}
	}
	pv, err := PreviewList(model.Step{Value: "a\nb"})
	if err != nil || pv.Count != 2 || pv.First != "a" {
		t.Errorf("xem trước dòng: %+v %v", pv, err)
	}
	if _, err := PreviewList(model.Step{Value: "{{n2.output.json}}"}); err == nil {
		t.Error("có biến node phải báo chỉ biết khi chạy")
	}
}
