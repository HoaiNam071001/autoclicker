//go:build !windows

package winexclude

// ExcludeOwnWindows chỉ có tác dụng trên Windows.
func ExcludeOwnWindows() int { return 0 }
