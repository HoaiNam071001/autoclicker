package model

import (
	"encoding/json"
	"os"
	"testing"
)

// Cùng dữ liệu với frontend/src/flow/window.test.ts (testdata/windows.json, window_presets.json).
func TestMatchWindowShared(t *testing.T) {
	var data struct {
		Wins  []WindowInfo
		Cases []struct {
			App, Title, Match string
			Want              any
		}
	}
	b, err := os.ReadFile("../../testdata/windows.json")
	if err == nil {
		err = json.Unmarshal(b, &data)
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range data.Cases {
		i, err := MatchWindow(data.Wins, c.App, c.Title, c.Match)
		if c.Want == "error" {
			if err == nil {
				t.Errorf("%q/%q %s: muốn lỗi", c.App, c.Title, c.Match)
			}
			continue
		}
		if err != nil || float64(i) != c.Want {
			t.Errorf("%q/%q %s: %d %v, muốn %v", c.App, c.Title, c.Match, i, err, c.Want)
		}
	}
}

func TestWindowPresetShared(t *testing.T) {
	var data struct {
		Cases []struct {
			Preset          string
			Work, Want, Out Rect
		}
	}
	b, err := os.ReadFile("../../testdata/window_presets.json")
	if err == nil {
		err = json.Unmarshal(b, &data)
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range data.Cases {
		if got := WindowPreset(c.Preset, c.Work, c.Want); got != c.Out {
			t.Errorf("%q %v %v: %v, muốn %v", c.Preset, c.Work, c.Want, got, c.Out)
		}
	}
}
