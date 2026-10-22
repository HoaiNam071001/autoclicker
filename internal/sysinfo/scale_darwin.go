package sysinfo

import (
	"github.com/go-vgo/robotgo"

	"ghosthand/internal/model"
)

func screenScale(model.Rect) float64 { return robotgo.ScaleF() }
