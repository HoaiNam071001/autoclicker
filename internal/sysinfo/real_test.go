package sysinfo

import (
	"testing"
	"time"
)

// Máy thật: các khoá cốt lõi phải có và hợp lý (không phụ thuộc có pin / mạng).
func TestSystemReadReal(t *testing.T) {
	start := time.Now()
	out, err := System{}.Read(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"batteryPercent", "cpuPercent", "memTotalMB", "diskFreeGB", "hostname", "os", "online"} {
		if _, ok := out[k]; !ok {
			t.Errorf("thiếu %s: %v", k, out)
		}
	}
	if p := out["batteryPercent"].(int); p < 0 || p > 100 {
		t.Errorf("pin %d", p)
	}
	t.Logf("%v (%s)", out, time.Since(start))
}
