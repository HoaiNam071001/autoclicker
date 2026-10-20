package store

import (
	"errors"
	"testing"
)

func TestRevisions(t *testing.T) {
	r := NewRevisions()
	ok := func() error { return nil }
	if n, err := r.Write("a.json", -1, ok); err != nil || n != 1 {
		t.Fatalf("lưu từ app: %d %v", n, err)
	}
	if n, err := r.Write("./a.json", 1, ok); err != nil || n != 2 || r.Get("a.json") != 2 {
		t.Fatalf("ghi đúng revision (đường dẫn tương đối cùng file): %d %v", n, err)
	}
	var ce *ConflictError
	if _, err := r.Write("a.json", 1, ok); !errors.As(err, &ce) || ce.Have != 2 || ce.Code() != "conflict" {
		t.Fatalf("revision cũ phải xung đột: %v", err)
	}
	if n, _ := r.Write("a.json", 2, func() error { return errors.New("đĩa đầy") }); n != 2 || r.Get("a.json") != 2 {
		t.Fatal("ghi lỗi không được tăng revision")
	}
	if r.Get("b.json") != 0 {
		t.Fatal("file chưa ghi = 0")
	}
}
