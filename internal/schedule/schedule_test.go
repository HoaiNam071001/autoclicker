package schedule

import (
	"testing"
	"time"

	"ghosthand/internal/model"
)

func TestDue(t *testing.T) {
	at := func(h, m, s int) time.Time { return time.Date(2026, 10, 20, h, m, s, 0, time.Local) }
	list := []model.Schedule{
		{Path: "a", Time: "08:00", Repeat: model.RepeatDaily},
		{Path: "b", Time: "08:01", Repeat: model.RepeatDaily},
		{Path: "c", Time: "xx", Repeat: model.RepeatDaily},
	}
	if d := Due(list, at(7, 59, 50), at(8, 0, 5)); len(d) != 1 || d[0].Path != "a" {
		t.Fatalf("nhịp qua 08:00: %v", d)
	}
	if d := Due(list, at(8, 0, 5), at(8, 0, 20)); len(d) != 0 {
		t.Fatalf("nhịp sau không chạy lại: %v", d)
	}
	if d := Due(list, at(7, 0, 0), at(8, 1, 0)); len(d) != 2 {
		t.Fatalf("hai hẹn trong một nhịp dài: %v", d)
	}
}
