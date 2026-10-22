package store

import (
	"sync"

	"ghosthand/internal/secrets"
)

var (
	secretsOnce  sync.Once
	secretsStore *secrets.Store
)

// Secrets – kho bí mật của máy (Keychain / Credential Manager + chỉ mục tên trong thư mục cài đặt), bước 43 S10.
func Secrets() *secrets.Store {
	secretsOnce.Do(func() {
		dir, _ := ConfigDir()
		secretsStore = &secrets.Store{Backend: secrets.System{}, Dir: dir}
	})
	return secretsStore
}
