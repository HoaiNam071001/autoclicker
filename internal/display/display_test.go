package display

import "testing"

func TestFindAudio(t *testing.T) {
	list := []Audio{
		{ID: "a", Name: "MacBook Pro Speakers"}, {ID: "b", Name: "AirPods Pro"}, {ID: "c", Name: "AirPods Pro", Input: true},
		{ID: "d", Name: "Pro"}, {ID: "e", Name: "MacBook Pro Microphone", Input: true},
	}
	for _, c := range []struct {
		q     string
		input bool
		want  int
	}{
		{"airpods", false, 1}, {"airpods", true, 2}, {"pro", false, 3}, {"c", true, 2}, {"microphone", true, 4}, {"zzz", false, -1}, {"", false, -1},
	} {
		if got := FindAudio(list, c.q, c.input); got != c.want {
			t.Errorf("FindAudio(%q, input=%v) = %d, muốn %d", c.q, c.input, got, c.want)
		}
	}
}
