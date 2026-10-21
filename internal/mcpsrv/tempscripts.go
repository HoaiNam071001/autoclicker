package mcpsrv

import (
	"sync"

	"ghosthand/internal/model"
)

// tempScripts – kịch bản tạm AI gửi kèm run_script (46 B1), theo runId: wait_run / resume_run cần lại kịch bản mà không
// có file. Chỉ trong bộ nhớ, giữ vài lần gần nhất – không sinh file rác.
type tempScripts struct {
	mu  sync.Mutex
	m   map[string]model.Script
	ids []string
}

const maxTempScripts = 8

func (t *tempScripts) put(id string, sc model.Script) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.m == nil {
		t.m = map[string]model.Script{}
	}
	t.m[id] = sc
	t.ids = append(t.ids, id)
	for len(t.ids) > maxTempScripts {
		delete(t.m, t.ids[0])
		t.ids = t.ids[1:]
	}
}

func (t *tempScripts) get(id string) (model.Script, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	sc, ok := t.m[id]
	return sc, ok
}
