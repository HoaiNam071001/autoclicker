package inputsource

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Carbon -framework Foundation
#import <Carbon/Carbon.h>

// Text Input Source API chỉ được gọi trên main queue (macOS 14+ assert) – Wails chạy vòng lặp chính nên dispatch_sync an toàn;
// đã ở main thread (chương trình thử) thì gọi thẳng.
static void on_main(void (^block)(void)) {
	if ([NSThread isMainThread]) block();
	else dispatch_sync(dispatch_get_main_queue(), block);
}

static NSString *prop(TISInputSourceRef s, CFStringRef key) {
	return (__bridge NSString *)TISGetInputSourceProperty(s, key);
}

static char *describe(TISInputSourceRef s) {
	NSString *id = prop(s, kTISPropertyInputSourceID) ?: @"";
	NSString *name = prop(s, kTISPropertyLocalizedName) ?: @"";
	NSArray *langs = (__bridge NSArray *)TISGetInputSourceProperty(s, kTISPropertyInputSourceLanguages);
	NSString *lang = langs.count ? langs[0] : @"";
	return strdup([[NSString stringWithFormat:@"%@\t%@\t%@", id, name, lang] UTF8String]);
}

static char *current_source(void) {
	__block char *out = NULL;
	on_main(^{
		TISInputSourceRef s = TISCopyCurrentKeyboardInputSource();
		if (s) { out = describe(s); CFRelease(s); }
	});
	return out;
}

// Mọi bộ gõ bàn phím chọn được, mỗi dòng "id\tname\tlang".
static char *list_sources(void) {
	__block char *out = NULL;
	on_main(^{
		NSDictionary *filter = @{
			(__bridge NSString *)kTISPropertyInputSourceIsSelectCapable: @YES,
			(__bridge NSString *)kTISPropertyInputSourceCategory: (__bridge NSString *)kTISCategoryKeyboardInputSource,
		};
		CFArrayRef all = TISCreateInputSourceList((__bridge CFDictionaryRef)filter, false);
		if (!all) return;
		NSMutableArray *lines = [NSMutableArray array];
		for (CFIndex i = 0; i < CFArrayGetCount(all); i++) {
			char *d = describe((TISInputSourceRef)CFArrayGetValueAtIndex(all, i));
			[lines addObject:[NSString stringWithUTF8String:d]];
			free(d);
		}
		CFRelease(all);
		out = strdup([[lines componentsJoinedByString:@"\n"] UTF8String]);
	});
	return out;
}

static int select_source(const char *id) {
	__block int ok = 0;
	on_main(^{
		NSDictionary *filter = @{(__bridge NSString *)kTISPropertyInputSourceID: [NSString stringWithUTF8String:id]};
		CFArrayRef found = TISCreateInputSourceList((__bridge CFDictionaryRef)filter, false);
		if (!found) return;
		if (CFArrayGetCount(found) > 0) ok = TISSelectInputSource((TISInputSourceRef)CFArrayGetValueAtIndex(found, 0)) == noErr;
		CFRelease(found);
	});
	return ok;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

func parseSource(line string) Source {
	f := strings.SplitN(line, "\t", 3)
	for len(f) < 3 {
		f = append(f, "")
	}
	return Source{ID: f[0], Name: f[1], Lang: f[2]}
}

func (System) Current() (Source, error) {
	c := C.current_source()
	if c == nil {
		return Source{}, errors.New("không đọc được bộ gõ hiện tại")
	}
	defer C.free(unsafe.Pointer(c))
	return parseSource(C.GoString(c)), nil
}

func (System) List() ([]Source, error) {
	c := C.list_sources()
	if c == nil {
		return nil, errors.New("không đọc được danh sách bộ gõ")
	}
	defer C.free(unsafe.Pointer(c))
	var out []Source
	for _, l := range strings.Split(C.GoString(c), "\n") {
		if l != "" {
			out = append(out, parseSource(l))
		}
	}
	return out, nil
}

func (System) Select(id string) error {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	if C.select_source(cid) == 0 {
		return fmt.Errorf("không đổi được sang bộ gõ %s", id)
	}
	return nil
}
