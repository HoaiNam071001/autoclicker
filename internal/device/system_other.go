//go:build !darwin && !windows

package device

import "errors"

var errNoAudio = errors.New("chỉnh âm lượng chưa hỗ trợ trên hệ điều hành này")

func (SystemAudio) Volume() (int, bool, error) { return 0, false, errNoAudio }
func (SystemAudio) SetVolume(int) error        { return errNoAudio }
func (SystemAudio) SetMuted(bool) error        { return errNoAudio }
