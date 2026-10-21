// Package humanize – thao tác giống người (bước 34): đường di chuột cong có tăng / giảm tốc, điểm bấm lệch quanh tâm,
// nhịp gõ / cuộn không đều. Hàm thuần – nhận nguồn ngẫu nhiên R để test bằng seed; engine lo thiết bị + chờ.
package humanize

import "ghosthand/internal/model"

// Mức (Script.Human, Step.Human – model.HumanLevels). "" ở kịch bản = tắt; ở node = theo kịch bản.
const (
	Off    = model.HumanOff
	Low    = model.HumanLow
	Medium = model.HumanMedium
	High   = model.HumanHigh
)

// R – nguồn ngẫu nhiên (*rand.Rand của math/rand/v2 thoả).
type R interface {
	Float64() float64
	NormFloat64() float64
	IntN(n int) int
}

// Profile – thông số một mức.
type Profile struct {
	Speed     float64 // nhân thời gian di chuột (lớn = chậm hơn)
	Curve     float64 // độ cong tối đa, tỷ lệ theo quãng đường
	Overshoot float64 // xác suất đi lố rồi quay lại (quãng > 80 px)
	Tremor    float64 // rung tay, px (độ lệch chuẩn) ở các điểm giữa đường
	ClickHold [2]int  // ms giữ nút mỗi lần click
	KeyGap    [2]int  // ms giữa hai ký tự / hai phím
	KeyHold   [2]int  // ms giữ một phím
	ScrollGap [2]int  // ms giữa hai nấc cuộn
	Pause     float64 // xác suất ngừng lâu hơn (sau dấu cách / giữa các nấc cuộn)
	PauseMs   [2]int
	Typo      float64 // xác suất gõ nhầm phím bên cạnh rồi xoá, mỗi chữ cái
}

var profiles = map[string]Profile{
	Low: {Speed: 0.6, Curve: 0.08, Overshoot: 0.05, Tremor: 0.3, ClickHold: [2]int{40, 80}, KeyGap: [2]int{25, 70},
		KeyHold: [2]int{30, 60}, ScrollGap: [2]int{30, 80}, Pause: 0.02, PauseMs: [2]int{150, 400}},
	Medium: {Speed: 1, Curve: 0.15, Overshoot: 0.15, Tremor: 0.6, ClickHold: [2]int{60, 150}, KeyGap: [2]int{50, 160},
		KeyHold: [2]int{40, 90}, ScrollGap: [2]int{50, 140}, Pause: 0.05, PauseMs: [2]int{250, 700}},
	High: {Speed: 1.3, Curve: 0.22, Overshoot: 0.25, Tremor: 0.9, ClickHold: [2]int{70, 180}, KeyGap: [2]int{70, 220},
		KeyHold: [2]int{50, 110}, ScrollGap: [2]int{70, 200}, Pause: 0.08, PauseMs: [2]int{300, 1000}, Typo: 0.03},
}

// ProfileOf – thông số của mức; false = tắt ("" / off / mức lạ).
func ProfileOf(level string) (Profile, bool) {
	p, ok := profiles[level]
	return p, ok
}

// Effective – mức dùng cho node: node ghi đè ("" = theo kịch bản), kịch bản "" = tắt.
func Effective(script, node string) string {
	if node != "" {
		return node
	}
	return script
}

// Between – số ngẫu nhiên đều trong [r[0], r[1]].
func Between(rnd R, r [2]int) int {
	if r[1] <= r[0] {
		return r[0]
	}
	return r[0] + rnd.IntN(r[1]-r[0]+1)
}
