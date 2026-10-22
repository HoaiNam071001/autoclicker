package sysinfo

import (
	"context"
	"os/exec"
)

func readBattery(ctx context.Context) Battery {
	out, err := exec.CommandContext(ctx, "pmset", "-g", "batt").Output()
	if err != nil {
		return Battery{}
	}
	return parsePmset(string(out))
}
