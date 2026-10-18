package samples

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

var testScreen = model.Rect{W: 1512, H: 982}

func TestSamplesValid(t *testing.T) {
	list, err := List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) == 0 {
		t.Fatal("không có mẫu nào")
	}
	for _, info := range list {
		sc, err := Open(info.ID, "file:///sandbox")
		if err != nil {
			t.Errorf("%s: %v", info.ID, err)
			continue
		}
		if err := script.ValidateScript(sc, testScreen); err != nil {
			t.Errorf("%s: %v", info.ID, err)
		}
		if info.Name == "" || info.Description == "" || info.App == "" || info.Level < 1 || info.Level > 3 {
			t.Errorf("%s: thiếu tên / mô tả / kiểu app / độ khó 1–3: %+v", info.ID, info)
		}
		if !info.Real {
			if _, err := pageFS.ReadFile("pages/" + info.Page); err != nil {
				t.Errorf("%s: trang sân thử %q không có", info.ID, info.Page)
			}
		}
		for _, st := range sc.Steps {
			if strings.Contains(st.Target, SandboxVar) || strings.HasPrefix(st.Template, templatePrefix) {
				t.Errorf("%s/%s: chưa thay {{sandbox}} / ảnh mẫu", info.ID, st.ID)
			}
			if page, ok := strings.CutPrefix(st.Target, "file:///sandbox/"); ok {
				page, _, _ = strings.Cut(page, "?")
				if _, err := pageFS.ReadFile("pages/" + page); err != nil {
					t.Errorf("%s/%s: mở trang %q không có trong sân thử", info.ID, st.ID, page)
				}
			}
		}
	}
}

// Mỗi loại bước app hỗ trợ phải có ít nhất một mẫu – thêm loại bước mới thì thêm mẫu (bước 33).
func TestEveryStepTypeHasSample(t *testing.T) {
	list, err := List()
	if err != nil {
		t.Fatal(err)
	}
	covered := map[string]bool{}
	for _, info := range list {
		for _, typ := range info.Steps {
			covered[typ] = true
		}
	}
	for _, d := range model.StepTypeDocs() {
		if !covered[d.Type] {
			t.Errorf("loại bước %q (%s) chưa có kịch bản mẫu nào trong internal/samples/scripts", d.Type, d.Label)
		}
	}
}

func TestWriteSandbox(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "sân thử")
	if err := WriteSandbox(dir); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"index.html", "click.html", "sandbox.css", "logo.png"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Error(err)
		}
	}
}

func TestFileURL(t *testing.T) {
	cases := map[string]string{
		"/Users/a/Library/Application Support/autoclicker/sandbox": "file:///Users/a/Library/Application%20Support/autoclicker/sandbox",
		`C:\Users\a\AppData\Roaming\autoclicker\sandbox`:          "file:///C:/Users/a/AppData/Roaming/autoclicker/sandbox",
	}
	for in, want := range cases {
		if filepath.Separator == '/' && strings.Contains(in, `\`) {
			continue // đường dẫn Windows chỉ đổi \ → / khi chạy trên Windows
		}
		if got := FileURL(in); got != want {
			t.Errorf("FileURL(%q) = %q, muốn %q", in, got, want)
		}
	}
}

func TestOpenUnknown(t *testing.T) {
	if _, err := Open("khong-co", ""); err == nil {
		t.Error("mẫu không có phải lỗi")
	}
}
