package device

import "testing"

func TestParseGUID(t *testing.T) {
	g := parseGUID("5a14bc41-5b76-3140-b680-8825562683ac")
	want := comGUID{0x5a14bc41, 0x5b76, 0x3140, [8]byte{0xb6, 0x80, 0x88, 0x25, 0x56, 0x26, 0x83, 0xac}}
	if *g != want {
		t.Fatalf("có %+v, muốn %+v", *g, want)
	}
	defer func() {
		if recover() == nil {
			t.Error("GUID sai định dạng phải panic")
		}
	}()
	parseGUID("5a14bc41-5b76")
}
