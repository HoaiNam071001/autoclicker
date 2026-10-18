package model

// Point – toạ độ điểm trong vùng, mỗi trục neo riêng (script/anchor.go): Ax / Ay trống = start (cách cạnh trái / trên).
// Điểm B kéo thả "lệch so với A" (DragRelative) là độ lệch, không dùng neo.
type Point struct {
	X  int    `json:"x"`
	Y  int    `json:"y"`
	Ax string `json:"ax,omitempty"` // start | end | center | ratio
	Ay string `json:"ay,omitempty"`
	// Toạ độ từ biến (bước 20 F5b): có một trong hai → điểm dùng biến, bỏ qua x / y + neo. Mỗi ô là số / biến cộng trừ px,
	// vd "{{n2.output.x}} + 20"; kết quả là toạ độ màn hình (như output chờ ảnh), không tính theo vùng.
	XVar string `json:"xVar,omitempty"`
	YVar string `json:"yVar,omitempty"`
}

// UsesVar – điểm lấy toạ độ từ biến (XVar / YVar) thay cho x / y + neo.
func (p Point) UsesVar() bool { return p.XVar != "" || p.YVar != "" }

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}
