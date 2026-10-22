package shortcuts

import (
	"reflect"
	"testing"
)

func TestParseList(t *testing.T) {
	got := parseList("Gửi tin\n\n  Đổi tên ảnh \nGửi tin\n")
	if want := []string{"Gửi tin", "Đổi tên ảnh"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
