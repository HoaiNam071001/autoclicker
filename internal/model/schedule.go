package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Hẹn giờ chạy kịch bản (bước 22): lưu ở schedules.json (store) theo đường dẫn file, không ghi vào file kịch bản.

const (
	RepeatOnce     = "once"
	RepeatDaily    = "daily"
	RepeatWeekdays = "weekdays"

	NotifyNone     = "none"
	NotifyTelegram = "telegram"
	NotifyDiscord  = "discord"
	NotifyWebhook  = "webhook"

	NotifyOnDone  = "done"
	NotifyOnError = "error"
	NotifyOnBoth  = "both"
)

type Schedule struct {
	Path     string `json:"path"`
	Time     string `json:"time"`     // HH:mm giờ máy
	Repeat   string `json:"repeat"`   // once | daily | weekdays
	Days     []int  `json:"days"`     // 1 = T2 … 7 = CN (repeat = weekdays)
	Preset   string `json:"preset"`   // bộ giá trị input ("" = giá trị mặc định)
	Minimize bool   `json:"minimize"` // thu nhỏ cửa sổ khi chờ tới giờ (Wails v2 không có khay)
	// Thông báo sau lần chạy hẹn giờ (internal/notify): telegram "botToken@chatId", discord / webhook = URL.
	Notify       string `json:"notify"`
	NotifyTarget string `json:"notifyTarget"`
	NotifyOn     string `json:"notifyOn"` // done | error | both
	Screenshot   bool   `json:"screenshot"`
}

// clock – giờ, phút của s.Time.
func (s Schedule) clock() (int, int, error) {
	h, m, ok := strings.Cut(s.Time, ":")
	hh, e1 := strconv.Atoi(h)
	mm, e2 := strconv.Atoi(m)
	if !ok || e1 != nil || e2 != nil || hh < 0 || hh > 23 || mm < 0 || mm > 59 {
		return 0, 0, fmt.Errorf("giờ hẹn %q không hợp lệ (cần HH:mm)", s.Time)
	}
	return hh, mm, nil
}

// Check – lỗi cấu hình hẹn giờ (trước khi lưu).
func (s Schedule) Check() error {
	if _, _, err := s.clock(); err != nil {
		return err
	}
	switch s.Repeat {
	case RepeatOnce, RepeatDaily:
	case RepeatWeekdays:
		if len(s.Days) == 0 {
			return fmt.Errorf("chưa chọn ngày nào trong tuần")
		}
		for _, d := range s.Days {
			if d < 1 || d > 7 {
				return fmt.Errorf("ngày %d không hợp lệ (1 = T2 … 7 = CN)", d)
			}
		}
	default:
		return fmt.Errorf("lặp lại %q không hợp lệ", s.Repeat)
	}
	switch s.Notify {
	case "", NotifyNone:
	case NotifyTelegram:
		if t, c, ok := strings.Cut(s.NotifyTarget, "@"); !ok || t == "" || c == "" {
			return fmt.Errorf("Telegram cần dạng botToken@chatId")
		}
	case NotifyDiscord, NotifyWebhook:
		if !strings.HasPrefix(s.NotifyTarget, "http://") && !strings.HasPrefix(s.NotifyTarget, "https://") {
			return fmt.Errorf("cần URL bắt đầu bằng https://")
		}
	default:
		return fmt.Errorf("kiểu thông báo %q không hợp lệ", s.Notify)
	}
	return nil
}

// isoDay – 1 = T2 … 7 = CN.
func isoDay(t time.Time) int { return (int(t.Weekday())+6)%7 + 1 }

// NextRun – lần chạy đầu tiên sau after (không tính đúng after) theo giờ địa phương của after; lỗi cấu hình → false.
func NextRun(s Schedule, after time.Time) (time.Time, bool) {
	hh, mm, err := s.clock()
	if err != nil {
		return time.Time{}, false
	}
	day := time.Date(after.Year(), after.Month(), after.Day(), hh, mm, 0, 0, after.Location())
	for i := 0; i < 8; i++ {
		t := day.AddDate(0, 0, i)
		if !t.After(after) {
			continue
		}
		if s.Repeat != RepeatWeekdays || containsInt(s.Days, isoDay(t)) {
			return t, true
		}
	}
	return time.Time{}, false
}

func containsInt(l []int, v int) bool {
	for _, x := range l {
		if x == v {
			return true
		}
	}
	return false
}
