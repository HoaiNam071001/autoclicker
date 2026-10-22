package session

import (
	"context"
	"testing"
	"time"
)

func TestKeepAwakeHeldWhileRunning(t *testing.T) {
	m, _, _ := testManager(t)
	var started, stopped int
	m.deps.KeepAwake = func(context.Context) (func(), error) {
		started++
		return func() { stopped++ }, nil
	}
	for _, keep := range []bool{false, true} {
		sc := waitScript(t, 1)
		sc.KeepAwake = keep
		id, err := m.Start(sc, nil, Options{})
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err := m.Wait(ctx, id); err != nil {
			t.Fatal(err)
		}
		cancel()
	}
	if started != 1 || stopped != 1 {
		t.Errorf("started=%d stopped=%d, muốn mỗi cái 1 (chỉ lần có keepAwake)", started, stopped)
	}
}
