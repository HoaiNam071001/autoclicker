// Package devicetest – thiết bị giả cho test (engine, failsafe): ghi lại thao tác chuột / phím thay vì làm thật.
package devicetest

import (
	"fmt"
	"strings"
	"sync"
)

// FakeInput ghi lại các lệnh dưới dạng chuỗi thay vì điều khiển chuột thật; Lines() trả lại để so.
type FakeInput struct {
	mu  sync.Mutex
	log []string
}

func (f *FakeInput) add(format string, a ...any) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.log = append(f.log, fmt.Sprintf(format, a...))
	return nil
}

func (f *FakeInput) Move(x, y int) error { return f.add("move %d,%d", x, y) }

func (f *FakeInput) Drag(x, y int, b string) error { return f.add("drag %d,%d %s", x, y, b) }

func (f *FakeInput) Toggle(b string, down bool) error { return f.add("toggle %s %v", b, down) }

func (f *FakeInput) Click(b string, double bool) error { return f.add("click %s %v", b, double) }

func (f *FakeInput) TypeStr(s string) error { return f.add("type %s", s) }

func (f *FakeInput) Scroll(d string) error { return f.add("scroll %s", d) }

func (f *FakeInput) KeyToggle(k string, down bool, mods []string) error {
	if !down {
		return f.add("up %s", k)
	}
	return f.add("down %s [%s]", k, strings.Join(mods, "+"))
}

func (f *FakeInput) Lines() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.log...)
}
