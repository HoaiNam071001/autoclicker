package main

import (
	"encoding/json"
	"testing"
)

// Mọi loại bước NewStep nhận đều có mô tả; ví dụ đọc được qua checkStep và đúng loại.
func TestStepTypeDocs(t *testing.T) {
	docs := map[string]bool{}
	for _, d := range stepTypeDocs {
		docs[d.Type] = true
		var s Step
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
	for _, typ := range []string{StepClick, StepDouble, StepHold, StepDrag, StepWait, StepText, StepKey, StepScroll, StepOpen, StepWaitImage, StepClickImage, StepWaitPixel, StepIf, StepGroup, StepBreak, StepStop} {
		if _, err := (&App{}).NewStep(typ); err != nil {
			t.Errorf("NewStep(%s): %v", typ, err)
		}
		if !docs[typ] {
			t.Errorf("thiếu mô tả cho loại %s", typ)
		}
	}
}
