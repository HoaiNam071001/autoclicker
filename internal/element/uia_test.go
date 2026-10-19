package element

import "testing"

func TestRoleFromUIA(t *testing.T) {
	cases := []struct {
		ct   int
		aria string
		want string
	}{
		{50000, "", "button"},
		{50004, "", "textbox"},
		{50020, "", RoleText},
		{50020, "heading", "heading"},       // Chrome: tiêu đề là Text + AriaRole
		{50004, " SearchBox ", "searchbox"}, // AriaRole thắng ControlType
		{50031, "", "button"},               // SplitButton
		{50020, "description", RoleText},    // Chrome: khối chữ
		{59999, "", "uia59999"},
	}
	for _, c := range cases {
		if got := RoleFromUIA(c.ct, c.aria); got != c.want {
			t.Errorf("RoleFromUIA(%d, %q) = %q, cần %q", c.ct, c.aria, got, c.want)
		}
	}
}
