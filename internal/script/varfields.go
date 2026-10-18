package script

import (
	"fmt"
	"slices"

	"autoclicker/internal/launch"
	"autoclicker/internal/model"
)

// VarField – ô chữ nhận {{biến}} của node (kể cả trong điều kiện / case).
type VarField struct {
	Name string // tên trường (json) cho lời báo lỗi, vd "cond.conds[1].left"
	Ptr  *string
	URL  bool // ô URL: giá trị biến được mã hoá (launch.EscapeURLValue)
	// Num: ô số – toạ độ (Point.XVar / YVar, F5b) hoặc trường số (NumVars, F5c): biến phải ra số (output của node),
	// cả ô là số + / − (EvalCoord)
	Num bool
}

// VarFields – các ô chữ nhận biến của s (input {{key}} ở resolveInputs, biến node ở resolveNodeVars, kiểm tra ở
// validateFlow). Con trỏ trỏ vào s – muốn sửa thì gọi trên bản sao (cloneStep).
func VarFields(s *model.Step) []VarField {
	var out []VarField
	add := func(name string, p *string, url bool) { out = append(out, VarField{Name: name, Ptr: p, URL: url}) }
	if hasPoints(s.Type) {
		for i := range s.Points {
			if p := &s.Points[i]; p.UsesVar() && !(i == 0 && s.Type == model.StepDrag && s.DragFromPrev) { // A "từ điểm trước" không dùng
				out = append(out, VarField{Name: fmt.Sprintf("points[%d].xVar", i), Ptr: &p.XVar, Num: true},
					VarField{Name: fmt.Sprintf("points[%d].yVar", i), Ptr: &p.YVar, Num: true})
			}
		}
	}
	for i := range s.NumVars {
		out = append(out, VarField{Name: fmt.Sprintf("numVars[%d].expr", i), Ptr: &s.NumVars[i].Expr, Num: true})
	}
	switch s.Type {
	case model.StepText:
		add("text", &s.Text, false)
	case model.StepWaitImage, model.StepClickImage:
		if s.FindBy == model.FindByText {
			add("findText", &s.FindText, false)
		}
	case model.StepWebWait, model.StepWebClick, model.StepWebType, model.StepWebRead, model.StepWebScroll:
		for i := range s.Selectors {
			add(fmt.Sprintf("selectors[%d]", i), &s.Selectors[i], false)
		}
		add("tab", &s.Tab, false)
		if s.Type == model.StepWebType {
			add("text", &s.Text, false)
		}
		if s.Type == model.StepWebScroll && s.ScrollMode == model.ScrollModeUntil {
			for i := range s.ScrollIn {
				add(fmt.Sprintf("scrollIn[%d]", i), &s.ScrollIn[i], false)
			}
		}
	case model.StepClipboard:
		switch s.ClipOp {
		case model.ClipWrite:
			add("text", &s.Text, false)
		case model.ClipImage:
			add("filePath", &s.FilePath, false)
		}
	case model.StepScreenshot:
		add("saveDir", &s.SaveDir, false)
		add("fileName", &s.FileName, false)
	case model.StepOpen:
		add("target", &s.Target, launch.TargetKind(s.Target) == launch.TargetURL)
	case model.StepIf:
		if s.Cond != nil {
			condFields(s.Cond, "cond", add)
		}
	case model.StepSwitch:
		add("value", &s.Value, false)
		for i := range s.Cases {
			add(fmt.Sprintf("cases[%d].match", i), &s.Cases[i].Match, false)
		}
	case model.StepSwitchScreen:
		for i := range s.Cases {
			if c := s.Cases[i].Cond; c != nil {
				condFields(c, fmt.Sprintf("cases[%d].cond", i), add)
			}
		}
	case model.StepEnd:
		add("message", &s.Message, false)
	case model.StepSetVar:
		add("value", &s.Value, false)
	}
	return out
}

// hasPoints – loại node có điểm trên màn hình nhận toạ độ từ biến (điểm trong điều kiện pixel chưa nhận).
func hasPoints(t string) bool {
	switch t {
	case model.StepClick, model.StepDouble, model.StepHold, model.StepDrag, model.StepScroll, model.StepWaitPixel:
		return true
	}
	return false
}

func condFields(c *model.Cond, name string, add func(string, *string, bool)) {
	switch c.Kind {
	case model.CondText:
		add(name+".findText", &c.FindText, false)
	case model.CondValue:
		add(name+".left", &c.Left, false)
		add(name+".right", &c.Right, false)
	case model.CondAll, model.CondAny:
		for i := range c.Conds {
			condFields(&c.Conds[i], fmt.Sprintf("%s.conds[%d]", name, i), add)
		}
	}
}

// CloneStep chép sâu phần lồng nhau sửa được (điểm, điều kiện, case) – thay biến trên bản sao không đụng kịch bản gốc.
func CloneStep(s model.Step) model.Step {
	s.Points = slices.Clone(s.Points)
	s.Selectors = slices.Clone(s.Selectors)
	s.ScrollIn = slices.Clone(s.ScrollIn)
	s.NumVars = slices.Clone(s.NumVars)
	s.Cond = cloneCond(s.Cond)
	if s.Cases != nil {
		cs := make([]model.Case, len(s.Cases))
		for i, c := range s.Cases {
			c.Cond = cloneCond(c.Cond)
			cs[i] = c
		}
		s.Cases = cs
	}
	return s
}
