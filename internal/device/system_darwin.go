package device

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func (SystemAudio) Volume() (int, bool, error) {
	out, err := run(exec.Command("osascript", "-e", "set s to get volume settings", "-e",
		`return ((output volume of s) as text) & "," & ((output muted of s) as text)`), "")
	if err != nil {
		return 0, false, fmt.Errorf("không đọc được âm lượng: %w", err)
	}
	return parseVolume(out)
}

func (SystemAudio) SetVolume(level int) error {
	if _, err := run(exec.Command("osascript", "-e", fmt.Sprintf("set volume output volume %d without output muted", level)), ""); err != nil {
		return fmt.Errorf("không đặt được âm lượng: %w", err)
	}
	return nil
}

func (SystemAudio) SetMuted(muted bool) error {
	if _, err := run(exec.Command("osascript", "-e", fmt.Sprintf("set volume output muted %v", muted)), ""); err != nil {
		return fmt.Errorf("không đổi được tắt tiếng: %w", err)
	}
	return nil
}

// parseVolume đọc "38,false" (máy không có loa ra: "missing value,…").
func parseVolume(out string) (int, bool, error) {
	v, m, _ := strings.Cut(strings.TrimSpace(out), ",")
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false, fmt.Errorf("máy không có loa ra để chỉnh âm lượng (%q)", strings.TrimSpace(out))
	}
	return n, m == "true", nil
}
