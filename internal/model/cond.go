package model

// Case – một nhánh của switch / switchScreen. Key duy nhất trong node, cổng ra là "case:" + Key.
type Case struct {
	Key   string `json:"key"`
	Label string `json:"label,omitempty"` // tên hiện trên cổng ("" = Key)
	Op    string `json:"op,omitempty"`    // switch: eq | contains | regex
	Match string `json:"match,omitempty"` // switch: so Value với chuỗi này; nhận {{biến}}
	Cond  *Cond  `json:"cond,omitempty"`  // switchScreen: điều kiện ảnh / chữ / màu
}

// Cond – điều kiện của node Nếu / case của Switch theo màn hình (cond.go). Kind image / text / pixel dùng cùng trường
// với node Chờ ảnh / chữ / pixel; value so 2 vế; all / any ghép nhiều điều kiện.
type Cond struct {
	Kind string `json:"kind"`          // image | text | element | pixel | value | all | any
	Not  bool   `json:"not,omitempty"` // đảo kết quả

	Template      string   `json:"template,omitempty"`
	TemplateScale float64  `json:"templateScale,omitempty"`
	Threshold     int      `json:"threshold,omitempty"`
	MatchMode     string   `json:"matchMode,omitempty"`
	ColorTol      int      `json:"colorTol,omitempty"`
	SearchArea    string   `json:"searchArea,omitempty"`
	FindText      string   `json:"findText,omitempty"` // nhận {{biến}}
	TextExact     bool     `json:"textExact,omitempty"`
	TextNoAccent  bool     `json:"textNoAccent,omitempty"`
	Occurrence    int      `json:"occurrence,omitempty"`
	Selectors     []string `json:"selectors,omitempty"`  // element: cách tìm phần tử (text= / aria= / role=); nhận {{biến}}
	ElementApp    string   `json:"elementApp,omitempty"` // element: tên app ("" = app có cửa sổ trên cùng); nhận {{biến}}
	Points        []Point  `json:"points,omitempty"`     // pixel: points[0]
	Color         string   `json:"color,omitempty"`
	Tolerance     int      `json:"tolerance,omitempty"`

	Left  string `json:"left,omitempty"`  // value: vế trái, vd {{n2.output.text}}; nhận {{biến}}
	Op    string `json:"op,omitempty"`    // value: eq | ne | contains | gt | lt | empty
	Right string `json:"right,omitempty"` // value: vế phải (không dùng với empty); nhận {{biến}}

	Conds []Cond `json:"conds,omitempty"` // all (và) / any (hoặc)
}

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
		MatchMode: c.MatchMode, ColorTol: c.ColorTol, SearchArea: c.SearchArea, FindText: c.FindText, TextExact: c.TextExact, TextNoAccent: c.TextNoAccent, Occurrence: c.Occurrence,
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
