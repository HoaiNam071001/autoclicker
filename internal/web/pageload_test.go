package web

import (
	"testing"
	"time"
)

func TestLoadWatch(t *testing.T) {
	t0 := time.Unix(0, 0)
	at := func(ms int) time.Time { return t0.Add(time.Duration(ms) * time.Millisecond) }
	type sample struct {
		ms    int
		state string
		nodes int
	}
	cases := []struct {
		name    string
		samples []sample
		wantAt  int // ms lần đầu settled; -1 = không bao giờ
	}{
		{"complete, đứng yên", []sample{{0, "complete", 200}, {900, "complete", 200}, {1500, "complete", 200}}, 1500},
		{"không bao giờ complete nhưng đã đứng yên", []sample{{0, "interactive", 500}, {1000, "interactive", 501}, {1600, "interactive", 502}}, 1600},
		{"còn loading thì chờ", []sample{{0, "loading", 100}, {1500, "loading", 100}, {3000, "loading", 100}}, -1},
		{"SPA đang dựng thêm phần tử", []sample{{0, "complete", 20}, {1000, "complete", 300}, {2000, "complete", 800}, {3600, "complete", 805}}, 3600},
		{"đồng hồ / carousel đổi ít", []sample{{0, "complete", 1000}, {800, "complete", 1010}, {1600, "complete", 995}}, 1600},
	}
	for _, c := range cases {
		var w loadWatch
		got := -1
		for _, s := range c.samples {
			if w.settled(loadProbe{State: s.state, Nodes: s.nodes}, at(s.ms)) {
				got = s.ms
				break
			}
		}
		if got != c.wantAt {
			t.Errorf("%s: settled lúc %d ms, muốn %d", c.name, got, c.wantAt)
		}
	}
}
