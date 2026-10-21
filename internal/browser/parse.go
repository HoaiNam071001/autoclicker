package browser

import (
	"encoding/json"
	"slices"
	"strings"

	"ghosthand/internal/model"
)

// rawBrowser – một trình duyệt trong kết quả lệnh list của tabs.js: mảng theo cửa sổ (trên cùng trước).
type rawBrowser struct {
	Browser string     `json:"browser"`
	Front   bool       `json:"front"`
	Windows []int64    `json:"windows"`
	Active  []int      `json:"active"` // tab đang hiện của từng cửa sổ, từ 1
	IDs     [][]int64  `json:"ids"`
	Titles  [][]string `json:"titles"`
	URLs    [][]string `json:"urls"`
}

// parseTabs đổi kết quả list ra danh sách tab phẳng: trình duyệt đang ở trước lên đầu (còn lại giữ thứ tự).
func parseTabs(out string) ([]model.TabInfo, error) {
	var raw []rawBrowser
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, err
	}
	slices.SortStableFunc(raw, func(a, b rawBrowser) int {
		switch {
		case a.Front == b.Front:
			return 0
		case a.Front:
			return -1
		}
		return 1
	})
	tabs := []model.TabInfo{}
	for _, b := range raw {
		for i, w := range b.Windows {
			if i >= len(b.IDs) || i >= len(b.Titles) || i >= len(b.URLs) {
				break
			}
			for j, id := range b.IDs[i] {
				t := model.TabInfo{Browser: b.Browser, ID: id, Window: w, Index: j + 1}
				if j < len(b.Titles[i]) && j < len(b.URLs[i]) {
					t.Title, t.URL = b.Titles[i][j], b.URLs[i][j]
				}
				t.Active = i < len(b.Active) && b.Active[i] == j+1
				tabs = append(tabs, t)
			}
		}
	}
	return tabs, nil
}

// osaError đổi lỗi của osascript (mã AppleEvent trong stderr) ra lỗi chung của package; khác → nil.
func osaError(stderr string) error {
	switch {
	case strings.Contains(stderr, "-1743"):
		return ErrNoAutomation
	case strings.Contains(stderr, "-600"):
		return ErrNotRunning
	case strings.Contains(stderr, "-1728"), strings.Contains(stderr, "-1719"):
		return ErrTabGone
	}
	return nil
}
