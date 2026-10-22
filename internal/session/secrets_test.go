package session

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
	"ghosthand/internal/secrets"
)

func secretScript() model.Script {
	return scripttest.Linear(model.Script{Name: "bi mat", Loops: 1, Steps: []model.Step{
		modeltest.Node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "pw", model.VarSet, "{{secret.mk}}" }),
		modeltest.Node("e", "n2", model.StepEnd, func(s *model.Step) { s.EndStatus, s.Message = "success", "pw={{vars.pw}}" }),
	}})
}

func TestSecretMaskedInHistory(t *testing.T) {
	m, _, saved := testManager(t)
	m.deps.Secrets = &secrets.Store{Backend: &secrets.MemBackend{M: map[string]string{"mk": "Ghost-1234"}}, Dir: t.TempDir()}
	id, err := m.Start(secretScript(), nil, Options{})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := m.Wait(ctx, id); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(*saved)
	if strings.Contains(string(b), "Ghost-1234") || !strings.Contains(string(b), "pw=••••") {
		t.Errorf("lịch sử chạy: %s", b)
	}
}

func TestSecretMissingBlocksStart(t *testing.T) {
	m, _, _ := testManager(t)
	m.deps.Secrets = &secrets.Store{Backend: &secrets.MemBackend{}, Dir: t.TempDir()}
	_, err := m.Start(secretScript(), nil, Options{})
	if err == nil || !strings.Contains(err.Error(), "thiếu bí mật: mk") {
		t.Errorf("err=%v", err)
	}
	if err := m.Validate(secretScript(), nil); err == nil {
		t.Error("Validate phải báo thiếu bí mật")
	}
}
