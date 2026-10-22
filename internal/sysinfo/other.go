//go:build !darwin && !windows

package sysinfo

import "context"

func readBattery(context.Context) Battery { return Battery{} }
func readWifi(context.Context) string     { return "" }
