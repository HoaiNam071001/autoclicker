//go:build !darwin && !windows

package inputsource

import "errors"

var errUnsupported = errors.New("hệ điều hành này chưa hỗ trợ đọc / đổi bộ gõ")

func (System) Current() (Source, error) { return Source{}, errUnsupported }
func (System) List() ([]Source, error)  { return nil, errUnsupported }
func (System) Select(string) error      { return errUnsupported }
