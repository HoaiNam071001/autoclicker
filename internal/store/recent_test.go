package store

import (
	"reflect"
	"strings"
	"testing"
)

func TestRecentList(t *testing.T) {
	var l []string
	for _, p := range []string{"a", "b", "c", "d", "e", "f", "c"} {
		l = PushRecent(l, p)
	}
	if want := []string{"c", "f", "e", "d", "b"}; !reflect.DeepEqual(l, want) {
		t.Fatalf("pushRecent = %v, muốn %v", l, want)
	}
	if got := RemoveRecent(l, "e"); strings.Join(got, "") != "cfdb" {
		t.Fatalf("removeRecent = %v", got)
	}
}
