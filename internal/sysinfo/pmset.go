package sysinfo

import (
	"regexp"
	"strconv"
	"strings"
)

var pmsetPct = regexp.MustCompile(`(\d+)%`)

// parsePmset đọc đầu ra `pmset -g batt`, vd:
//
//	Now drawing from 'Battery Power'
//	 -InternalBattery-0 (id=1)	87%; discharging; 4:10 remaining present: true
func parsePmset(s string) Battery {
	for _, line := range strings.Split(s, "\n") {
		if !strings.Contains(line, "InternalBattery") {
			continue
		}
		m := pmsetPct.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		p, _ := strconv.Atoi(m[1])
		l := strings.ToLower(line)
		// “discharging” chứa “charging” → loại trước; “charged” = đầy, đang cắm sạc
		charging := !strings.Contains(l, "discharging") && (strings.Contains(l, "charging") || strings.Contains(l, "charged") || strings.Contains(l, "finishing charge"))
		return Battery{Has: true, Percent: p, Charging: charging}
	}
	return Battery{}
}
