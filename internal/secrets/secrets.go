// Package secrets – kho bí mật (bước 43 S10): giá trị nằm trong Keychain (macOS) / Credential Manager (Windows),
// không bao giờ ghi vào file kịch bản, lịch sử chạy hay cài đặt. Chỉ danh sách TÊN được lưu ở file chỉ mục để liệt kê.
package secrets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sync"
)

var (
	ErrNotFound    = errors.New("chưa có bí mật tên này")
	ErrUnsupported = errors.New("kho bí mật chỉ có trên macOS và Windows")
)

// MinLen – độ dài tối thiểu của giá trị; engine chỉ che giá trị từ độ dài này (engine/mask.go).
const MinLen = 4

// nameRe – tên dùng được trong {{secret.<tên>}}.
var nameRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

// ValidName báo tên có dùng được trong {{secret.<tên>}} không.
func ValidName(name string) error {
	if !nameRe.MatchString(name) {
		return fmt.Errorf("tên bí mật %q không hợp lệ – chữ cái đầu, sau đó chữ / số / gạch dưới (vd mat_khau_email)", name)
	}
	return nil
}

// Backend – nơi giữ giá trị. Test thay bản giả (MemBackend).
type Backend interface {
	Get(ctx context.Context, name string) (string, error) // không có → ErrNotFound
	Set(ctx context.Context, name, value string) error
	Delete(ctx context.Context, name string) error // không có thì thôi
}

// Store ghép Backend với chỉ mục tên (file JSON trong Dir).
type Store struct {
	Backend Backend
	Dir     string
	mu      sync.Mutex
}

func (s *Store) indexPath() string { return filepath.Join(s.Dir, "secret-names.json") }

// Names – tên các bí mật đã lưu, theo thứ tự chữ cái.
func (s *Store) Names() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.names()
}

func (s *Store) names() []string {
	var out []string
	if b, err := os.ReadFile(s.indexPath()); err == nil {
		_ = json.Unmarshal(b, &out)
	}
	slices.Sort(out)
	return out
}

func (s *Store) saveNames(names []string) error {
	b, _ := json.Marshal(names)
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.indexPath(), b, 0o600)
}

// Set lưu (hoặc đổi) giá trị của bí mật name.
func (s *Store) Set(ctx context.Context, name, value string) error {
	if err := ValidName(name); err != nil {
		return err
	}
	if len(value) < MinLen {
		return fmt.Errorf("giá trị bí mật tối thiểu %d ký tự (ngắn hơn không che được trong nhật ký)", MinLen)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.Backend.Set(ctx, name, value); err != nil {
		return err
	}
	names := s.names()
	if !slices.Contains(names, name) {
		names = append(names, name)
		slices.Sort(names)
	}
	return s.saveNames(names)
}

// Delete xoá bí mật name (không có thì thôi).
func (s *Store) Delete(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.Backend.Delete(ctx, name); err != nil {
		return err
	}
	return s.saveNames(slices.DeleteFunc(s.names(), func(n string) bool { return n == name }))
}

// Resolve lấy giá trị các bí mật names; thiếu → trả danh sách tên thiếu (giá trị khác vẫn có trong map).
func (s *Store) Resolve(ctx context.Context, names []string) (values map[string]string, missing []string, err error) {
	values = map[string]string{}
	for _, n := range names {
		v, gerr := s.Backend.Get(ctx, n)
		switch {
		case errors.Is(gerr, ErrNotFound):
			missing = append(missing, n)
		case gerr != nil:
			return nil, nil, fmt.Errorf("không đọc được bí mật %q: %w", n, gerr)
		default:
			values[n] = v
		}
	}
	return values, missing, nil
}

// Has báo bí mật name đã có chưa (không đọc giá trị ra ngoài).
func (s *Store) Has(name string) bool { return slices.Contains(s.Names(), name) }

// MemBackend – Backend trong bộ nhớ, cho test.
type MemBackend struct {
	mu sync.Mutex
	M  map[string]string
}

func (b *MemBackend) Get(_ context.Context, name string) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if v, ok := b.M[name]; ok {
		return v, nil
	}
	return "", ErrNotFound
}

func (b *MemBackend) Set(_ context.Context, name, value string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.M == nil {
		b.M = map[string]string{}
	}
	b.M[name] = value
	return nil
}

func (b *MemBackend) Delete(_ context.Context, name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.M, name)
	return nil
}
