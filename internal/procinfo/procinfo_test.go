package procinfo

import (
	"os"
	"testing"
)

func TestCPUPercent(t *testing.T) {
	if got := cpuPercent(0.15, 0.3); got != 50 {
		t.Errorf("got %v", got)
	}
	if cpuPercent(-1, 0.3) != 0 || cpuPercent(1, 0) != 0 {
		t.Error("số âm / khoảng 0 phải ra 0")
	}
}

func TestSystemListFindsSelf(t *testing.T) {
	list, err := System{}.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range list {
		if p.PID == os.Getpid() {
			if p.Name == "" || p.MemMB <= 0 {
				t.Errorf("tiến trình của chính test: %+v", p)
			}
			return
		}
	}
	t.Error("không thấy tiến trình của chính test")
}
