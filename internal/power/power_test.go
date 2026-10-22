package power

import (
	"strings"
	"testing"
)

func TestCommandCoversEveryActionOnBothSystems(t *testing.T) {
	for _, goos := range []string{"darwin", "windows"} {
		for _, a := range []string{Lock, DisplaySleep, Sleep, Logout, Shutdown, Restart} {
			argv, err := command(goos, a)
			if err != nil || len(argv) == 0 {
				t.Errorf("%s/%s: %v %v", goos, a, argv, err)
			}
		}
	}
}

func TestCommandRejectsUnknown(t *testing.T) {
	if _, err := command("darwin", "format"); err == nil || !strings.Contains(err.Error(), "không hợp lệ") {
		t.Errorf("err=%v", err)
	}
	if _, err := command("linux", Lock); err != ErrUnsupported {
		t.Errorf("err=%v", err)
	}
}
