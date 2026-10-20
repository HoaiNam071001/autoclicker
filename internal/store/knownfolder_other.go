//go:build !windows

package store

import "path/filepath"

func knownFolders(home string) (docs, pics string) {
	return filepath.Join(home, "Documents"), filepath.Join(home, "Pictures")
}
