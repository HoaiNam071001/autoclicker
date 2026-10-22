//go:build !darwin && !windows

package display

import "context"

// System – hệ khác chưa hỗ trợ.
type System struct{}

func (System) Dark() (bool, error)                           { return false, ErrUnsupported }
func (System) SetDark(bool) error                            { return ErrUnsupported }
func (System) Brightness() (int, error)                      { return 0, ErrUnsupported }
func (System) SetBrightness(int) error                       { return ErrUnsupported }
func (System) AudioDevices() ([]Audio, error)                { return nil, ErrUnsupported }
func (System) SetDefaultAudio(string, bool) error            { return ErrUnsupported }
func (System) PlaySound(context.Context, string, bool) error { return ErrUnsupported }
