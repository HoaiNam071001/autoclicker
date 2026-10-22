package secrets

import (
	"context"
	"testing"
)

func TestStoreNamesAndResolve(t *testing.T) {
	s := &Store{Backend: &MemBackend{}, Dir: t.TempDir()}
	ctx := context.Background()
	for _, n := range []string{"zalo", "email_pw"} {
		if err := s.Set(ctx, n, "v-"+n); err != nil {
			t.Fatal(err)
		}
	}
	if got := s.Names(); len(got) != 2 || got[0] != "email_pw" {
		t.Errorf("names=%v", got)
	}
	vals, missing, err := s.Resolve(ctx, []string{"zalo", "khong_co"})
	if err != nil || vals["zalo"] != "v-zalo" || len(missing) != 1 || missing[0] != "khong_co" {
		t.Errorf("vals=%v missing=%v err=%v", vals, missing, err)
	}
	if err := s.Delete(ctx, "zalo"); err != nil || s.Has("zalo") {
		t.Errorf("xoá: %v %v", err, s.Names())
	}
}

func TestStoreRejectsBadInput(t *testing.T) {
	s := &Store{Backend: &MemBackend{}, Dir: t.TempDir()}
	for _, n := range []string{"", "1abc", "a.b", "a b", "mật_khẩu"} {
		if s.Set(context.Background(), n, "x") == nil {
			t.Errorf("tên %q lẽ ra bị từ chối", n)
		}
	}
	if s.Set(context.Background(), "ok", "abc") == nil {
		t.Error("giá trị trống lẽ ra bị từ chối")
	}
}

// Kho thật của hệ điều hành: ghi – đọc – xoá một mục thử rồi dọn.
func TestSystemRoundTrip(t *testing.T) {
	var b Backend = System{}
	ctx := context.Background()
	name := "ghosthand_test_secret"
	if err := b.Set(ctx, name, "mật khẩu ✓ 123"); err != nil {
		t.Skipf("kho bí mật không dùng được ở đây: %v", err)
	}
	defer b.Delete(ctx, name)
	if v, err := b.Get(ctx, name); err != nil || v != "mật khẩu ✓ 123" {
		t.Fatalf("đọc lại: %q %v", v, err)
	}
	if err := b.Set(ctx, name, "đổi"); err != nil {
		t.Fatal(err)
	}
	if v, _ := b.Get(ctx, name); v != "đổi" {
		t.Errorf("sau khi đổi: %q", v)
	}
	b.Delete(ctx, name)
	if _, err := b.Get(ctx, name); err != ErrNotFound {
		t.Errorf("sau khi xoá: %v", err)
	}
}
