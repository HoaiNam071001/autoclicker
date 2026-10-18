package model

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

// Thứ tự hiển thị: cùng dữ liệu với frontend/src/flow.test.ts (testdata/flow_order.json) – hai bên phải cùng công thức.
func TestFlowOrderShared(t *testing.T) {
	b, err := os.ReadFile("../../testdata/flow_order.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name   string `json:"name"`
		Script Script `json:"script"`
		Order  []string
	}
	if err := json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		if got := FlowOrder(c.Script); !reflect.DeepEqual(got, c.Order) {
			t.Errorf("%s: %q, muốn %q", c.Name, got, c.Order)
		}
	}
	if len(cases) == 0 {
		t.Fatal("không có ca nào")
	}
}
