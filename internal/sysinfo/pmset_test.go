package sysinfo

import "testing"

func TestParsePmset(t *testing.T) {
	cases := []struct {
		in   string
		want Battery
	}{
		{"Now drawing from 'Battery Power'\n -InternalBattery-0 (id=1)\t87%; discharging; 4:10 remaining present: true\n", Battery{true, 87, false}},
		{"Now drawing from 'AC Power'\n -InternalBattery-0 (id=1)\t45%; charging; 1:00 remaining present: true\n", Battery{true, 45, true}},
		{" -InternalBattery-0 (id=1)\t100%; charged; 0:00 remaining present: true", Battery{true, 100, true}},
		{"Now drawing from 'AC Power'\n", Battery{}},
	}
	for _, c := range cases {
		if got := parsePmset(c.in); got != c.want {
			t.Errorf("%q → %+v, muốn %+v", c.in, got, c.want)
		}
	}
}

func TestResolveDiskPath(t *testing.T) {
	if resolveDiskPath("") == "" || resolveDiskPath("/") != "/" {
		t.Error("đường dẫn ổ đĩa")
	}
}
