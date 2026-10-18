//go:build !darwin && !windows

package device

import "errors"

var errNoTrash = errors.New("Thùng rác chưa hỗ trợ trên hệ điều hành này")

func (SystemTrash) MoveToTrash([]string) error { return errNoTrash }
func (SystemTrash) EmptyTrash() error          { return errNoTrash }
