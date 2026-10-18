package main

import (
	"encoding/json"
	"testing"

	"autoclicker/internal/model"
)

// Kịch bản flow mẫu cho AI đọc được và chạy qua kiểm tra đầy đủ (input lấy mặc định).
func TestFlowExampleValidates(t *testing.T) {
	sc, err := decodeScript([]byte(model.Guide().Example))
	if err != nil {
		t.Fatal(err)
	}
	sc, err = resolveInputs(sc, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateScript(sc, model.Rect{W: 1920, H: 1080}); err != nil {
		t.Fatal(err)
	}
}

// Mọi loại bước NewStep nhận đều có mô tả; ví dụ đọc được qua checkStep và đúng loại.
func TestStepTypeDocs(t *testing.T) {
	docs := map[string]bool{}
	for _, d := range model.StepTypeDocs() {
		docs[d.Type] = true
		var s model.Step
		if err := json.Unmarshal([]byte(d.Example), &s); err != nil {
			t.Errorf("%s: ví dụ không phải JSON: %v", d.Type, err)
			continue
		}
		if err := checkStep(&s); err != nil || s.Type != d.Type {
			t.Errorf("%s: ví dụ không hợp lệ: %v (type %q)", d.Type, err, s.Type)
		}
		if s.Template != "" {
			if _, err := decodePNG(s.Template); err != nil {
				t.Errorf("%s: ảnh mẫu trong ví dụ hỏng: %v", d.Type, err)
			}
		}
		for _, f := range d.Pinnable {
			if _, ok := pinFields[f]; !ok {
				t.Errorf("%s: %q không có trong pinFields", d.Type, f)
			}
		}
	}
	for _, typ := range model.StepTypes {
		if _, err := model.NewStep(typ); err != nil {
			t.Errorf("NewStep(%s): %v", typ, err)
		}
		if !docs[typ] {
			t.Errorf("thiếu mô tả cho loại %s", typ)
		}
	}
}
