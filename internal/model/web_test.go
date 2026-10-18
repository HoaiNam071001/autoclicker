package model

import (
	"encoding/json"
	"os"
	"testing"
)

// Cùng dữ liệu với frontend/src/flow/web.test.ts (testdata/selectors.json) – hai bên phải đọc selector giống nhau.
func TestParseSelector(t *testing.T) {
	b, err := os.ReadFile("../../testdata/selectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		In   string   `json:"in"`
		Want Selector `json:"want"`
		Bad  bool     `json:"bad"`
	}
	if err := json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		got, err := ParseSelector(c.In)
		if c.Bad {
			if err == nil {
				t.Errorf("%q: cần lỗi, có %+v", c.In, got)
			}
			continue
		}
		if err != nil || got != c.Want {
			t.Errorf("%q: %+v, %v – cần %+v", c.In, got, err, c.Want)
		}
	}
}
