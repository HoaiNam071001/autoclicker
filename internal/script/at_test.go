package script

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

// Cùng dữ liệu với frontend/src/flow/waitUntil.test.ts (testdata/wait_until.json).
func TestParseAt(t *testing.T) {
	b, err := os.ReadFile("../../testdata/wait_until.json")
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Now   struct{ Y, M, D int }
		Cases []struct {
			At                    string
			Ok, Dated             bool
			Y, M, D, H, Mi, S, Ms int
		}
	}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatal(err)
	}
	now := time.Date(data.Now.Y, time.Month(data.Now.M), data.Now.D, 15, 0, 0, 0, time.Local)
	for _, c := range data.Cases {
		got, dated, err := ParseAt(c.At, now)
		if (err == nil) != c.Ok {
			t.Errorf("ParseAt(%q): lỗi %v, muốn ok = %v", c.At, err, c.Ok)
			continue
		}
		want := time.Date(c.Y, time.Month(c.M), c.D, c.H, c.Mi, c.S, c.Ms*int(time.Millisecond), time.Local)
		if c.Ok && (!got.Equal(want) || dated != c.Dated) {
			t.Errorf("ParseAt(%q) = %v %v, muốn %v %v", c.At, got, dated, want, c.Dated)
		}
	}
}
