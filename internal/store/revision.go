package store

import (
	"fmt"
	"path/filepath"
	"sync"
)

// Revision kịch bản (bước 22): số tăng mỗi lần file kịch bản được ghi (app lưu, API / AI ghi) – chỉ trong bộ nhớ,
// không ghi vào file. Bên ghi từ ngoài gửi kèm revision đã đọc; khác revision hiện tại → ConflictError (đã có người
// sửa sau lần đọc đó, ghi đè sẽ mất).

// ConflictError – ghi kịch bản bằng revision cũ.
type ConflictError struct {
	Path      string
	Have, Got int // revision hiện tại / bên ghi gửi
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("kịch bản %s đã đổi (revision %d, bạn gửi %d) – đọc lại rồi sửa", filepath.Base(e.Path), e.Have, e.Got)
}

// Code – mã lỗi cho API / MCP.
func (e *ConflictError) Code() string { return "conflict" }

type Revisions struct {
	mu sync.Mutex
	m  map[string]int
}

func NewRevisions() *Revisions { return &Revisions{m: map[string]int{}} }

func revKey(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		return filepath.Clean(abs)
	}
	return path
}

// Get – revision hiện tại của file (0 = chưa ghi lần nào trong lần mở app này).
func (r *Revisions) Get(path string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.m[revKey(path)]
}

// Write ghi file qua write rồi tăng revision; base < 0 = không kiểm tra (app tự lưu), base khác hiện tại → ConflictError.
func (r *Revisions) Write(path string, base int, write func() error) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := revKey(path)
	if base >= 0 && base != r.m[k] {
		return r.m[k], &ConflictError{Path: path, Have: r.m[k], Got: base}
	}
	if err := write(); err != nil {
		return r.m[k], err
	}
	r.m[k]++
	return r.m[k], nil
}
