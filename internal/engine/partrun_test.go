package engine

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Chạy đến node (31b): chạy xong StopAt thì dừng, không lỗi, Reached báo đã tới.
func TestRunStopAt(t *testing.T) {
	e, in := sendFlow(0, nil)
	e.stopAt = "c"
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !e.Reached() || len(in.Lines()) != 2 || e.results["t"] != nil {
		t.Errorf("reached %v, lệnh %q, gõ đã chạy %v", e.Reached(), in.Lines(), e.results["t"] != nil)
	}
}

// Chạy tiếp từ node (resume_run): biến node trước lấy từ Carry của lần chạy trước.
func TestRunSeed(t *testing.T) {
	e, _ := sendFlow(0, nil)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	carry := e.Carry()
	if carry.Results["w"].Output == nil {
		t.Fatalf("Carry thiếu output w: %+v", carry)
	}

	e2, in := sendFlow(0, nil)
	e2.sc = script.RunFrom(e2.sc, "t", false)
	e2.seed = &carry
	if err := e2.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	typed := ""
	for _, l := range in.Lines() {
		typed += strings.TrimPrefix(l, "type ")
	}
	if typed != "ở 220,160" { // chỗ thấy chữ "Gửi" ở lần chạy trước
		t.Errorf("gõ %q", typed)
	}

	e3, _ := sendFlow(0, nil) // không seed: biến node trước trống → node lỗi
	e3.sc = script.RunFrom(e3.sc, "t", false)
	if err := e3.Run(context.Background()); err == nil {
		t.Error("không seed mà vẫn chạy được")
	}
}

// Lỗi không thấy chữ có detail: chữ cần tìm, các dòng đọc được (kèm khung), vùng, ngữ cảnh; không dòng nào giống → không closest.
func TestFailDetail(t *testing.T) {
	e, _ := sendFlow(-1, nil)
	e.sc = script.RunFrom(e.sc, "w", true)
	_ = e.Run(context.Background())
	r := e.results["w"]
	if r == nil || r.Error == nil {
		t.Fatal("không có lỗi")
	}
	d := r.Error.Detail
	if d["kind"] != "text" || d["want"] != "Gửi" || !reflect.DeepEqual(d["read"], []map[string]any{{"text": "Đang tải", "box": model.Rect{X: 100, Y: 50, W: 80, H: 20}}}) || d["closest"] != nil || d["context"] == nil {
		t.Errorf("detail = %+v", d)
	}
}
