package model

// Cổng ra của node điều khiển (ngoài success / fail ở flow.go).
const (
	PortTrue    = "true"
	PortFalse   = "false"
	PortDefault = "default" // switch: không case nào khớp
	PortNone    = "none"    // switchScreen: hết giờ không case nào hiện
	PortNext    = "next"    // counter: còn lượt
	PortDone    = "done"    // counter: hết lượt
	CasePrefix  = "case:"
)

// Step.EndStatus / Step.VarOp
const (
	EndSuccess = "success"
	EndFail    = "fail"
	VarSet     = "set"
	VarAdd     = "add"
)

// casePorts – cổng case:<key> theo thứ tự case, rồi cổng cuối (default / none).
func casePorts(cases []Case, last string) []string {
	out := make([]string, 0, len(cases)+1)
	for _, c := range cases {
		out = append(out, CasePrefix+c.Key)
	}
	return append(out, last)
}
