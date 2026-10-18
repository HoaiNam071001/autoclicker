package script

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// Giờ hẹn của node Chờ đến giờ (bước 35 N6). Cùng công thức parseAt (frontend/src/flow/waitUntil.ts), test chung
// testdata/wait_until.json.

var atRe = regexp.MustCompile(`^\s*(?:(\d{4})-(\d{1,2})-(\d{1,2})|(\d{1,2})/(\d{1,2})/(\d{4}))?\s*(\d{1,2}):(\d{2})(?::(\d{2})(?:[.,](\d{1,3}))?)?\s*$`)

// ParseAt đọc giờ hẹn theo giờ địa phương của now: "HH:mm[:ss[.mmm]]" (hôm nay), có thể kèm ngày "YYYY-MM-DD" /
// "DD/MM/YYYY" phía trước. dated: có ghi ngày.
func ParseAt(at string, now time.Time) (t time.Time, dated bool, err error) {
	m := atRe.FindStringSubmatch(at)
	if m == nil {
		return time.Time{}, false, errors.New("giờ sai dạng – vd 20:00, 20:00:00.500, 2026-11-11 20:00 hoặc 11/11/2026 20:00")
	}
	n := func(s string) int { v, _ := strconv.Atoi(s); return v }
	y, mo, d := now.Date()
	switch {
	case m[1] != "":
		y, mo, d, dated = n(m[1]), time.Month(n(m[2])), n(m[3]), true
	case m[4] != "":
		y, mo, d, dated = n(m[6]), time.Month(n(m[5])), n(m[4]), true
	}
	h, mi, sec := n(m[7]), n(m[8]), n(m[9])
	ms := 0
	if m[10] != "" {
		ms = n((m[10] + "00")[:3]) // .5 = 500 ms
	}
	if h > 23 || mi > 59 || sec > 59 || mo < 1 || mo > 12 || d < 1 || d > 31 {
		return time.Time{}, false, errors.New("giờ / ngày ngoài khoảng (giờ 0–23, phút / giây 0–59)")
	}
	t = time.Date(y, mo, d, h, mi, sec, ms*int(time.Millisecond), now.Location())
	if t.Day() != d { // 31/02 → Go tự lăn sang tháng sau
		return time.Time{}, false, errors.New("ngày không có trong tháng")
	}
	return t, dated, nil
}
