package model

import "testing"

func TestRunRecordAddKeepsLastAndEndMessage(t *testing.T) {
	var r RunRecord
	for i := 0; i < MaxRunResults+3; i++ {
		r.Add(StepResult{NodeID: "a", Attempts: i}, StepText)
	}
	if len(r.Results) != MaxRunResults || r.Truncated != 3 || r.Results[0].Attempts != 3 || r.Nodes != MaxRunResults+3 {
		t.Errorf("giữ %d kết quả, bỏ %d, đầu %d", len(r.Results), r.Truncated, r.Results[0].Attempts)
	}
	r.Add(StepResult{Status: StatusSuccess, Output: map[string]any{"message": "Đã gửi"}}, StepEnd)
	r.Add(StepResult{Status: StatusSuccess, Output: map[string]any{"message": "không phải end"}}, StepText)
	if r.Message != "Đã gửi" {
		t.Errorf("lời nhắn = %q", r.Message)
	}
}
