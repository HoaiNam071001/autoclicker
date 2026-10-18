//go:build !darwin && !windows

package device

import (
	"errors"
	"os/exec"
)

var errNoClipboard = errors.New("clipboard chưa hỗ trợ trên hệ điều hành này")

func hideWindow(*exec.Cmd) {}

func (SystemClipboard) ReadText() (string, error) { return "", errNoClipboard }
func (SystemClipboard) WriteText(string) error    { return errNoClipboard }
func (SystemClipboard) WriteImage(string) error   { return errNoClipboard }
