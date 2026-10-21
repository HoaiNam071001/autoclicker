package device

import (
	"os"
	"testing"
	"time"
)

// Đọc chữ cả màn hình chính bằng OCR của hệ điều hành – chỉ khi GHOSTHAND_REAL=1. Lần 2 đo tốc độ khi engine đã dựng.
func TestRealOCR(t *testing.T) {
	if os.Getenv("GHOSTHAND_REAL") != "1" {
		t.Skip("đặt GHOSTHAND_REAL=1 để đọc chữ trên màn hình thật")
	}
	img, err := RobotgoScreen{}.Capture(MainScreen())
	if err != nil {
		t.Fatal(err)
	}
	for round := 1; round <= 2; round++ {
		start := time.Now()
		lines, err := SystemOCR{}.ReadText(img)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("lần %d: %d dòng trong %v", round, len(lines), time.Since(start))
		if round == 2 {
			for i, l := range lines[:min(len(lines), 60)] {
				t.Logf("  %2d %v %q (%d từ)", i, l.Box, l.Text, len(l.Words))
			}
		}
		if len(lines) == 0 {
			t.Error("không đọc được dòng nào trên cả màn hình")
		}
	}
}
