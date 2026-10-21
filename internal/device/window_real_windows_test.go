package device

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/model"
)

// Mở Notepad thật rồi đặt khung / thu nhỏ / focus / đóng – chỉ khi GHOSTHAND_REAL=1 (cửa sổ thật nhảy trên màn).
// Notepad Windows 11 là app đóng gói: notepad.exe mở xong thoát, tiến trình thật khác pid – tìm theo cửa sổ.
func TestRealWindowsNotepad(t *testing.T) {
	if os.Getenv("GHOSTHAND_REAL") != "1" {
		t.Skip("đặt GHOSTHAND_REAL=1 để điều khiển cửa sổ thật")
	}
	ws, apps := SystemWindows{}, SystemApps{}
	before := map[uint64]bool{}
	list, _ := ws.List()
	for _, w := range list {
		before[w.ID] = true
	}
	if err := exec.Command("notepad.exe").Start(); err != nil {
		t.Fatal(err)
	}
	var win model.WindowInfo
	for end := time.Now().Add(10 * time.Second); win.ID == 0; time.Sleep(200 * time.Millisecond) {
		if time.Now().After(end) {
			t.Fatal("không thấy cửa sổ Notepad mới")
		}
		list, _ = ws.List()
		for _, w := range list {
			if !before[w.ID] && strings.EqualFold(w.App, "notepad") {
				win = w
			}
		}
	}
	t.Logf("cửa sổ %+v", win)
	defer apps.ForceQuit(win.PID) // dọn nếu Close không đóng được
	want := model.Rect{X: 120, Y: 90, W: 800, H: 600}
	if err := ws.SetBounds(win.ID, want); err != nil {
		t.Fatal(err)
	}
	if got := find(t, ws, win.ID); got.Box != want {
		t.Errorf("khung sau SetBounds %v, cần %v", got.Box, want)
	}
	t.Logf("vùng làm việc %v", ws.WorkArea(want))
	if err := ws.SetState(win.ID, WinMinimize); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	if !find(t, ws, win.ID).Minimized {
		t.Error("chưa thu nhỏ")
	}
	if err := ws.Focus(win.ID); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	if w := find(t, ws, win.ID); w.Minimized || w.Box != want {
		t.Errorf("sau Focus: %+v (cần hết thu nhỏ, khung %v)", w, want)
	}
	all, _ := apps.List()
	found := false
	for _, a := range all {
		found = found || a.PID == win.PID && a.Windows > 0
	}
	if !found {
		t.Errorf("Apps.List không có pid %d kèm cửa sổ", win.PID)
	}
	if err := ws.Close(win.ID); err != nil {
		t.Fatal(err)
	}
	for end := time.Now().Add(5 * time.Second); ; time.Sleep(200 * time.Millisecond) {
		list, _ = ws.List()
		gone := true
		for _, w := range list {
			gone = gone && w.ID != win.ID
		}
		if gone {
			break
		}
		if time.Now().After(end) {
			t.Error("Close chưa đóng cửa sổ sau 5 s")
			break
		}
	}
}

func find(t *testing.T, ws SystemWindows, id uint64) model.WindowInfo {
	t.Helper()
	list, _ := ws.List()
	for _, w := range list {
		if w.ID == id {
			return w
		}
	}
	t.Fatalf("mất cửa sổ %d", id)
	return model.WindowInfo{}
}
