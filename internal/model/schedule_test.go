package model

import (
	"testing"
	"time"
)

func TestNextRun(t *testing.T) {
	loc := time.FixedZone("ICT", 7*3600)
	at := func(d, h, m int) time.Time { return time.Date(2026, 10, d, h, m, 0, 0, loc) } // 2026-10-19 là thứ Hai
	cases := []struct {
		name  string
		s     Schedule
		after time.Time
		want  time.Time
	}{
		{"hôm nay chưa tới giờ", Schedule{Time: "08:00", Repeat: RepeatDaily}, at(20, 7, 59), at(20, 8, 0)},
		{"đúng giờ thì sang mai", Schedule{Time: "08:00", Repeat: RepeatDaily}, at(20, 8, 0), at(21, 8, 0)},
		{"một lần đã qua giờ → mai", Schedule{Time: "08:00", Repeat: RepeatOnce}, at(20, 10, 0), at(21, 8, 0)},
		{"chọn ngày: T2 → T6", Schedule{Time: "09:30", Repeat: RepeatWeekdays, Days: []int{5}}, at(19, 12, 0), at(23, 9, 30)},
		{"chọn ngày: CN", Schedule{Time: "00:00", Repeat: RepeatWeekdays, Days: []int{7}}, at(19, 0, 0), at(25, 0, 0)},
		{"chọn ngày hôm nay chưa tới giờ", Schedule{Time: "23:59", Repeat: RepeatWeekdays, Days: []int{1}}, at(19, 23, 0), at(19, 23, 59)},
	}
	for _, c := range cases {
		got, ok := NextRun(c.s, c.after)
		if !ok || !got.Equal(c.want) {
			t.Errorf("%s: %v %v, cần %v", c.name, got, ok, c.want)
		}
	}
	if _, ok := NextRun(Schedule{Time: "25:00", Repeat: RepeatDaily}, at(19, 0, 0)); ok {
		t.Error("giờ sai phải false")
	}
}

func TestScheduleCheck(t *testing.T) {
	ok := []Schedule{
		{Time: "08:00", Repeat: RepeatDaily},
		{Time: "8:05", Repeat: RepeatWeekdays, Days: []int{1, 7}, Notify: NotifyTelegram, NotifyTarget: "123:ABC@456"},
		{Time: "23:59", Repeat: RepeatOnce, Notify: NotifyWebhook, NotifyTarget: "https://x.test/hook"},
	}
	for _, s := range ok {
		if err := s.Check(); err != nil {
			t.Errorf("%+v: %v", s, err)
		}
	}
	bad := []Schedule{
		{Time: "8h", Repeat: RepeatDaily},
		{Time: "08:00", Repeat: "monthly"},
		{Time: "08:00", Repeat: RepeatWeekdays},
		{Time: "08:00", Repeat: RepeatWeekdays, Days: []int{8}},
		{Time: "08:00", Repeat: RepeatDaily, Notify: NotifyTelegram, NotifyTarget: "thiếu-chat"},
		{Time: "08:00", Repeat: RepeatDaily, Notify: NotifyDiscord, NotifyTarget: "discord.com"},
	}
	for _, s := range bad {
		if s.Check() == nil {
			t.Errorf("%+v phải lỗi", s)
		}
	}
}
