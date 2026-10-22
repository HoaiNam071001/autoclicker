//go:build !darwin && !windows

package secrets

import "context"

// System – hệ khác chưa có kho bí mật.
type System struct{}

func (System) Get(context.Context, string) (string, error) { return "", ErrUnsupported }
func (System) Set(context.Context, string, string) error   { return ErrUnsupported }
func (System) Delete(context.Context, string) error        { return nil }
