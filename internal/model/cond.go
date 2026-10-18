package model

// Cond.Kind
const (
	CondImage   = "image"
	CondText    = "text"
	CondElement = "element" // thấy phần tử màn hình (bước 37) – Selectors, ElementApp
	CondPixel   = "pixel"
	CondValue   = "value" // so 2 vế Left Op Right (đã thay {{biến}})
	CondAll     = "all"
	CondAny     = "any"
)

// Probe – node Chờ ảnh / chữ tương ứng điều kiện (dùng lại vùng tìm, ảnh mẫu, tìm chữ). key: id riêng để nhớ ảnh mẫu.
func (c Cond) Probe(key string) Step {
	s := Step{
		ID: key, Type: StepWaitImage, Template: c.Template, TemplateScale: c.TemplateScale, Threshold: c.Threshold,
		SearchArea: c.SearchArea, FindText: c.FindText, TextExact: c.TextExact, TextNoAccent: c.TextNoAccent, Occurrence: c.Occurrence,
		Selectors: c.Selectors, ElementApp: c.ElementApp,
	}
	switch c.Kind {
	case CondText:
		s.FindBy = FindByText
	case CondElement:
		s.FindBy = FindByElement
	}
	return s
}
