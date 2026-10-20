package session

import (
	"context"
	"errors"
	"image"
	"sync"
	"testing"
	"time"

	"autoclicker/internal/device"
	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/web"
)

type events struct {
	mu   sync.Mutex
	list []string
}

func (e *events) emit(ev string, _ any) {
	e.mu.Lock()
	e.list = append(e.list, ev)
	e.mu.Unlock()
}

func (e *events) has(ev string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, x := range e.list {
		if x == ev {
			return true
		}
	}
	return false
}

func testManager(t *testing.T) (*Manager, *events, *[]model.RunRecord) {
	t.Helper()
	ev := &events{}
	var saved []model.RunRecord
	var mu sync.Mutex
	m := New(ev.emit, Deps{
		Input:       func() device.Input { return &devicetest.FakeInput{} },
		Screen:      func() model.Rect { return model.Rect{W: 1920, H: 1080} },
		Displays:    func() []model.Rect { return []model.Rect{{W: 1920, H: 1080}} },
		AllowDanger: func(model.Script, string, string) (bool, error) { return false, nil },
		Browser:     func(func(x, y int) error) (web.Browser, func(), error) { return nil, func() {}, nil },
		Corners:     func(context.Context, *device.TrackingInput, func()) {},
		SaveRun: func(rec model.RunRecord, _ string, _ map[string]image.Image) error {
			mu.Lock()
			saved = append(saved, rec)
			mu.Unlock()
			return nil
		},
	})
	return m, ev, &saved
}

// Kịch bản một node Chờ ms.
func waitScript(t *testing.T, ms int) model.Script {
	t.Helper()
	sc := model.NewScript()
	s, err := model.NewStep("wait")
	if err != nil {
		t.Fatal(err)
	}
	s.ID, s.Ref, s.DelayAfter = "s1", "n1", ms
	sc.Steps = []model.Step{s}
	sc.Start = "s1"
	sc.Loops = 1
	return sc
}

func TestStartWaitSuccess(t *testing.T) {
	m, ev, saved := testManager(t)
	id, err := m.Start(waitScript(t, 1), nil, Options{Source: SourceCLI})
	if err != nil {
		t.Fatal(err)
	}
	st, err := m.Wait(context.Background(), id)
	if err != nil || st.State != StateSuccess || st.Source != SourceCLI || st.RunID != id {
		t.Fatalf("status %+v err %v", st, err)
	}
	if len(*saved) != 1 || (*saved)[0].ID != id {
		t.Fatalf("lịch sử phải ghi một lần với runId: %+v", *saved)
	}
	for _, e := range []string{"run:started", "engine:state", "engine:done", "run:ended"} {
		if !ev.has(e) {
			t.Errorf("thiếu event %s", e)
		}
	}
	if m.Running() != "" {
		t.Error("xong rồi mà vẫn đang chạy")
	}
}

func TestBusyAndStop(t *testing.T) {
	m, _, saved := testManager(t)
	id, err := m.Start(waitScript(t, 1), nil, Options{Countdown: 5})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.Start(waitScript(t, 1), nil, Options{}); !errors.Is(err, ErrBusy) {
		t.Fatalf("phiên thứ hai phải báo bận, được %v", err)
	}
	if st, _ := m.Status(""); st.RunID != id || st.State != StateCountdown {
		t.Fatalf("status lúc đếm ngược %+v", st)
	}
	if m.Stop("khac") {
		t.Error("dừng nhầm runId khác")
	}
	if !m.Stop(id) {
		t.Fatal("không dừng được")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	st, err := m.Wait(ctx, id)
	if err != nil || st.State != StateStopped {
		t.Fatalf("dừng lúc đếm ngược: %+v %v", st, err)
	}
	if len(*saved) != 0 {
		t.Error("dừng lúc đếm ngược không ghi lịch sử")
	}
}

func TestValidateErrorAndUnknown(t *testing.T) {
	m, _, _ := testManager(t)
	sc := waitScript(t, 1)
	sc.Start = "khong-co"
	if _, err := m.Start(sc, nil, Options{}); err == nil {
		t.Fatal("start sai phải báo lỗi ngay")
	}
	if _, err := m.Status("x"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("runId lạ: %v", err)
	}
	if _, err := m.Wait(context.Background(), "x"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("runId lạ: %v", err)
	}
}
