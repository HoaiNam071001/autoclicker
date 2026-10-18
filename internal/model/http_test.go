package model

import (
	"encoding/json"
	"os"
	"testing"
)

// Cùng dữ liệu với frontend/src/flow/http.test.ts (testdata/secret_keys.json).
func TestIsSecretKey(t *testing.T) {
	b, err := os.ReadFile("../../testdata/secret_keys.json")
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Cases []struct {
			Name   string
			Secret bool
		}
	}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatal(err)
	}
	for _, c := range data.Cases {
		if got := IsSecretKey(c.Name); got != c.Secret {
			t.Errorf("IsSecretKey(%q) = %v, muốn %v", c.Name, got, c.Secret)
		}
	}
}
