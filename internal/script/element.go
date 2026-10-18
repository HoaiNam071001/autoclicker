package script

import (
	"errors"
	"strings"

	"autoclicker/internal/element"
	"autoclicker/internal/model"
)

// checkElementFind – node tìm phần tử (Chờ / Click theo phần tử, Đọc phần tử) đủ dữ liệu để chạy.
func checkElementFind(i int, s model.Step) error {
	if err := checkElementSelectors(s.Selectors); err != nil {
		return stepErr(i, s, "selectors", "%v", err)
	}
	return checkFindWait(i, s)
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
