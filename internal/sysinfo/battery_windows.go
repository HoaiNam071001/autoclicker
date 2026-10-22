package sysinfo

import (
	"context"
	"syscall"
	"unsafe"
)

var procGetSystemPowerStatus = syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemPowerStatus")

type powerStatus struct {
	ACLineStatus        byte
	BatteryFlag         byte
	BatteryLifePercent  byte
	SystemStatusFlag    byte
	BatteryLifeTime     uint32
	BatteryFullLifeTime uint32
}

func readBattery(context.Context) Battery {
	var ps powerStatus
	if r, _, _ := procGetSystemPowerStatus.Call(uintptr(unsafe.Pointer(&ps))); r == 0 {
		return Battery{}
	}
	// BatteryFlag: 128 = không có pin, 255 = không rõ; BatteryLifePercent 255 = không rõ
	if ps.BatteryFlag&128 != 0 || ps.BatteryFlag == 255 || ps.BatteryLifePercent > 100 {
		return Battery{}
	}
	return Battery{Has: true, Percent: int(ps.BatteryLifePercent), Charging: ps.ACLineStatus == 1}
}
