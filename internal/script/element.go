package script

import (
	"errors"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"ghosthand/internal/element"
	"ghosthand/internal/model"
)

// checkElementFind – node tìm phần tử (Chờ / Click theo phần tử, Đọc phần tử) đủ dữ liệu để chạy.
func checkElementFind(i int, s model.Step) error {
	if err := checkElementSelectors(s.Selectors); err != nil {
		return stepErr(i, s, "selectors", "%v", err)
	}
	return checkFindWait(i, s)
}

// checkElementRead – node Đọc phần tử: tìm bằng cây phần tử, hoặc tìm nhãn bằng OCR (37 E4 – khi đó phải đọc cạnh nhãn);
// mẫu lọc phải là regex đọc được.
func checkElementRead(i int, s model.Step) error {
	if s.FindBy == model.FindByText {
		if strings.TrimSpace(s.FindText) == "" {
			return stepErr(i, s, "findText", "chưa nhập chữ của nhãn cần tìm")
		}
		if s.ReadNear == "" {
			return stepErr(i, s, "readNear", "tìm nhãn bằng chữ (OCR) thì phải chọn đọc bên phải / bên dưới nhãn")
		}
		if err := checkFindWait(i, s); err != nil {
			return err
		}
	} else if err := checkElementFind(i, s); err != nil {
		return err
	}
	if s.ReadSpan < 0 || s.ReadSpan > 2000 {
		return stepErr(i, s, "readSpan", "khoảng đọc phải trong 0–2000 px")
	}
	if _, err := regexp.Compile(s.ReadPattern); err != nil {
		return stepErr(i, s, "readPattern", "mẫu lọc sai: %v", err)
	}
	return nil
}

// checkElementSelectors – cách tìm phần tử màn hình (bước 37): có ít nhất một dòng, dòng không có {{biến}} phải đọc được
// và là text= / aria= / role= (css / xpath chỉ có ở node web). Dòng có biến chỉ kiểm được lúc chạy.
func checkElementSelectors(lines []string) error {
	n := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		n++
		if strings.Contains(line, "{{") {
			continue
		}
		sel, err := model.ParseSelector(line)
		if err != nil {
			return err
		}
		if err := element.CheckSelector(sel); err != nil {
			return err
		}
	}
	if n == 0 {
		return errors.New("chưa nhập cách tìm phần tử")
	}
	return nil
}

// checkPageText – node Đọc chữ trang: selectors có thể trống (đọc cả cửa sổ); có thì phải là text= / aria= / role=.
func checkPageText(i int, s model.Step) error {
	if slices.ContainsFunc(s.Selectors, func(l string) bool { return strings.TrimSpace(l) != "" }) {
		if err := checkElementSelectors(s.Selectors); err != nil {
			return stepErr(i, s, "selectors", "%v", err)
		}
	}
	return checkFindWait(i, s)
}

// checkTableFile – node Đọc bảng: file ghi ra (nếu có) phải là .csv hoặc .xlsx.
func checkTableFile(i int, s model.Step) error {
	p := strings.TrimSpace(s.FilePath)
	if p == "" || strings.Contains(p, "{{") {
		return nil
	}
	if ext := strings.ToLower(filepath.Ext(p)); ext != ".csv" && ext != ".xlsx" {
		return stepErr(i, s, "filePath", "chỉ ghi được file .csv hoặc .xlsx (file này: %q)", filepath.Base(p))
	}
	return nil
}

// checkElementAction – node Thao tác phần tử: tìm được phần tử, thao tác hợp lệ, đặt giá trị thì có giá trị.
func checkElementAction(i int, s model.Step) error {
	if s.ElementAct != "" && !slices.Contains(model.ElementActs, s.ElementAct) {
		return stepErr(i, s, "elementAct", "thao tác không hợp lệ: %q (press | setValue | check | uncheck | expand | collapse | select | scroll)", s.ElementAct)
	}
	return checkElementFind(i, s)
}

// checkWaitUIEvent – node Chờ sự kiện giao diện: sự kiện hợp lệ; sự kiện về phần tử phải có selector để biết theo dõi gì.
func checkWaitUIEvent(i int, s model.Step) error {
	ev := s.UIEvent
	if ev == "" {
		ev = model.UIEventNewWindow
	}
	if !slices.Contains(model.UIEvents, ev) {
		return stepErr(i, s, "uiEvent", "sự kiện không hợp lệ: %q (newWindow | focus | newElement | valueChange)", s.UIEvent)
	}
	if s.TimeoutMs <= 0 {
		return stepErr(i, s, "timeoutMs", "thời gian chờ tối đa phải > 0")
	}
	if ev == model.UIEventNewElement || ev == model.UIEventValueChange {
		if err := checkElementSelectors(s.Selectors); err != nil {
			return stepErr(i, s, "selectors", "%v", err)
		}
	}
	return nil
}
