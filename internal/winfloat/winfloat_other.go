//go:build !darwin

// Package winfloat: chế độ nổi trên app toàn màn hình – chỉ cần trên macOS (Space riêng); Windows always-on-top đã đủ.
package winfloat

func PrepareLaunch() {}

func Set(on, focus bool) {}
