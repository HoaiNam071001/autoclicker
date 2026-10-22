//go:build !darwin

package sysinfo

import (
	"ghosthand/internal/device"
	"ghosthand/internal/model"
)

func screenScale(r model.Rect) float64 { return device.DPIScale(r) }
