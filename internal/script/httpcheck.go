package script

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"ghosthand/internal/model"
)

// checkHTTPStep – node HTTP đủ dữ liệu để chạy (bước 35 N3). URL còn biến node {{n2…}} / body còn biến thì chỉ kiểm tra được lúc chạy.
func checkHTTPStep(i int, s model.Step) error {
	u := strings.TrimSpace(s.URL)
	if u == "" || u == "https://" || u == "http://" {
		return stepErr(i, s, "url", "chưa nhập URL cần gọi")
	}
	if !NodeVarRe.MatchString(u) {
		if p, err := url.Parse(u); err != nil || (p.Scheme != "http" && p.Scheme != "https") || p.Host == "" {
			return stepErr(i, s, "url", "URL “%s” sai dạng – cần http:// hoặc https://", u)
		}
	}
	for _, t := range []struct {
		name string
		rows []model.KV
	}{{"query", s.Query}, {"headers", s.Headers}, {"form", s.Form}} {
		for j, kv := range t.rows {
			if strings.TrimSpace(kv.Key) == "" && kv.Value != "" {
				return stepErr(i, s, fmt.Sprintf("%s[%d].key", t.name, j), "dòng %d của %s có giá trị nhưng chưa có tên", j+1, t.name)
			}
		}
	}
	if s.BodyType == model.BodyJSON && strings.TrimSpace(s.Body) != "" && !strings.Contains(s.Body, "{{") && !json.Valid([]byte(s.Body)) {
		return stepErr(i, s, "body", "thân JSON sai cú pháp")
	}
	return nil
}
