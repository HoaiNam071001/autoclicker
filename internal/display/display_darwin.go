package display

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework CoreAudio -framework CoreGraphics -framework Foundation
#import <Foundation/Foundation.h>
#import <CoreAudio/CoreAudio.h>
#import <CoreGraphics/CoreGraphics.h>
#include <dlfcn.h>

// ---- Độ sáng: DisplayServices là framework riêng của Apple (không có header) – nạp lúc chạy, thiếu thì báo không hỗ trợ.
typedef int (*ds_get_fn)(CGDirectDisplayID, float *);
typedef int (*ds_set_fn)(CGDirectDisplayID, float);

static void *ds_handle(void) {
	static void *h; static int tried;
	if (!tried) { tried = 1; h = dlopen("/System/Library/PrivateFrameworks/DisplayServices.framework/DisplayServices", RTLD_LAZY); }
	return h;
}

// Màn hình tích hợp của laptop trước (màn ngoài thường không chỉnh được bằng phần mềm), không có thì màn chính.
static CGDirectDisplayID brightness_display(void) {
	CGDirectDisplayID ids[16]; uint32_t n = 0;
	if (CGGetActiveDisplayList(16, ids, &n) == kCGErrorSuccess)
		for (uint32_t i = 0; i < n; i++) if (CGDisplayIsBuiltin(ids[i])) return ids[i];
	return CGMainDisplayID();
}

// 0 = ok; 1 = không hỗ trợ
static int brightness_get(float *out) {
	void *h = ds_handle(); if (!h) return 1;
	ds_get_fn f = (ds_get_fn)dlsym(h, "DisplayServicesGetBrightness"); if (!f) return 1;
	return f(brightness_display(), out) == 0 ? 0 : 1;
}
static int brightness_set(float v) {
	void *h = ds_handle(); if (!h) return 1;
	ds_set_fn f = (ds_set_fn)dlsym(h, "DisplayServicesSetBrightness"); if (!f) return 1;
	return f(brightness_display(), v) == 0 ? 0 : 1;
}

// ---- Âm thanh
static NSString *audio_uid(AudioDeviceID d) {
	AudioObjectPropertyAddress a = {kAudioDevicePropertyDeviceUID, kAudioObjectPropertyScopeGlobal, kAudioObjectPropertyElementMain};
	CFStringRef s = NULL; UInt32 n = sizeof(s);
	if (AudioObjectGetPropertyData(d, &a, 0, NULL, &n, &s) != 0 || !s) return nil;
	return (__bridge_transfer NSString *)s;
}
static NSString *audio_name(AudioDeviceID d) {
	AudioObjectPropertyAddress a = {kAudioObjectPropertyName, kAudioObjectPropertyScopeGlobal, kAudioObjectPropertyElementMain};
	CFStringRef s = NULL; UInt32 n = sizeof(s);
	if (AudioObjectGetPropertyData(d, &a, 0, NULL, &n, &s) != 0 || !s) return nil;
	return (__bridge_transfer NSString *)s;
}
static int audio_has_streams(AudioDeviceID d, AudioObjectPropertyScope scope) {
	AudioObjectPropertyAddress a = {kAudioDevicePropertyStreams, scope, kAudioObjectPropertyElementMain};
	UInt32 n = 0;
	return AudioObjectGetPropertyDataSize(d, &a, 0, NULL, &n) == 0 && n > 0;
}
static AudioDeviceID audio_default(int input) {
	AudioObjectPropertyAddress a = {input ? kAudioHardwarePropertyDefaultInputDevice : kAudioHardwarePropertyDefaultOutputDevice,
		kAudioObjectPropertyScopeGlobal, kAudioObjectPropertyElementMain};
	AudioDeviceID d = kAudioObjectUnknown; UInt32 n = sizeof(d);
	AudioObjectGetPropertyData(kAudioObjectSystemObject, &a, 0, NULL, &n, &d);
	return d;
}

// Mỗi dòng "uid\tname\tinput(0/1)\tdefault(0/1)"; thiết bị vừa loa vừa micro ra hai dòng.
static char *audio_list(void) {
	AudioObjectPropertyAddress a = {kAudioHardwarePropertyDevices, kAudioObjectPropertyScopeGlobal, kAudioObjectPropertyElementMain};
	UInt32 size = 0;
	if (AudioObjectGetPropertyDataSize(kAudioObjectSystemObject, &a, 0, NULL, &size) != 0) return NULL;
	AudioDeviceID *ids = malloc(size);
	if (AudioObjectGetPropertyData(kAudioObjectSystemObject, &a, 0, NULL, &size, ids) != 0) { free(ids); return NULL; }
	AudioDeviceID defOut = audio_default(0), defIn = audio_default(1);
	NSMutableArray *lines = [NSMutableArray array];
	for (UInt32 i = 0; i < size / sizeof(AudioDeviceID); i++) {
		NSString *uid = audio_uid(ids[i]), *name = audio_name(ids[i]);
		if (!uid || !name) continue;
		if (audio_has_streams(ids[i], kAudioObjectPropertyScopeOutput))
			[lines addObject:[NSString stringWithFormat:@"%@\t%@\t0\t%d", uid, name, ids[i] == defOut]];
		if (audio_has_streams(ids[i], kAudioObjectPropertyScopeInput))
			[lines addObject:[NSString stringWithFormat:@"%@\t%@\t1\t%d", uid, name, ids[i] == defIn]];
	}
	free(ids);
	return strdup([[lines componentsJoinedByString:@"\n"] UTF8String]);
}

// 0 = ok; 1 = không thấy thiết bị; 2 = hệ điều hành từ chối
static int audio_set_default(const char *uid, int input) {
	AudioObjectPropertyAddress a = {kAudioHardwarePropertyDevices, kAudioObjectPropertyScopeGlobal, kAudioObjectPropertyElementMain};
	UInt32 size = 0;
	if (AudioObjectGetPropertyDataSize(kAudioObjectSystemObject, &a, 0, NULL, &size) != 0) return 2;
	AudioDeviceID *ids = malloc(size);
	if (AudioObjectGetPropertyData(kAudioObjectSystemObject, &a, 0, NULL, &size, ids) != 0) { free(ids); return 2; }
	NSString *want = [NSString stringWithUTF8String:uid];
	AudioDeviceID found = kAudioObjectUnknown;
	for (UInt32 i = 0; i < size / sizeof(AudioDeviceID); i++) {
		if ([audio_uid(ids[i]) isEqualToString:want] && audio_has_streams(ids[i], input ? kAudioObjectPropertyScopeInput : kAudioObjectPropertyScopeOutput)) { found = ids[i]; break; }
	}
	free(ids);
	if (found == kAudioObjectUnknown) return 1;
	AudioObjectPropertyAddress s = {input ? kAudioHardwarePropertyDefaultInputDevice : kAudioHardwarePropertyDefaultOutputDevice,
		kAudioObjectPropertyScopeGlobal, kAudioObjectPropertyElementMain};
	return AudioObjectSetPropertyData(kAudioObjectSystemObject, &s, 0, NULL, sizeof(found), &found) == 0 ? 0 : 2;
}
*/
import "C"

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"
)

// System – Backend thật cho macOS.
type System struct{}

func (System) Dark() (bool, error) {
	// AppleInterfaceStyle chỉ có khi đang tối; thiếu key = sáng (lệnh trả lỗi, không phải lỗi thật)
	out, err := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle").Output()
	if err != nil {
		return false, nil
	}
	return strings.EqualFold(strings.TrimSpace(string(out)), "dark"), nil
}

func (System) SetDark(dark bool) error {
	script := fmt.Sprintf(`tell application "System Events" to tell appearance preferences to set dark mode to %t`, dark)
	if out, err := exec.Command("osascript", "-e", script).CombinedOutput(); err != nil {
		return fmt.Errorf("không đổi được chế độ sáng / tối: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func (System) Brightness() (int, error) {
	var v C.float
	if C.brightness_get(&v) != 0 {
		return 0, ErrNoBrightness
	}
	return int(math.Round(float64(v) * 100)), nil
}

func (System) SetBrightness(percent int) error {
	if C.brightness_set(C.float(float64(percent)/100)) != 0 {
		return ErrNoBrightness
	}
	return nil
}

func (System) AudioDevices() ([]Audio, error) {
	c := C.audio_list()
	if c == nil {
		return nil, fmt.Errorf("không đọc được danh sách thiết bị âm thanh")
	}
	defer C.free(unsafe.Pointer(c))
	var out []Audio
	for _, line := range strings.Split(C.GoString(c), "\n") {
		f := strings.Split(line, "\t")
		if len(f) == 4 {
			out = append(out, Audio{ID: f[0], Name: f[1], Input: f[2] == "1", Default: f[3] == "1"})
		}
	}
	return out, nil
}

func (System) SetDefaultAudio(id string, input bool) error {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	in := 0
	if input {
		in = 1
	}
	switch C.audio_set_default(cid, C.int(in)) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("không thấy thiết bị âm thanh %q", id)
	}
	return fmt.Errorf("hệ điều hành từ chối đổi thiết bị âm thanh mặc định")
}

func (System) PlaySound(ctx context.Context, target string, wait bool) error {
	path := target
	if !strings.ContainsAny(target, `/\`) && filepath.Ext(target) == "" { // tên âm hệ thống: Glass → /System/Library/Sounds/Glass.aiff
		path = "/System/Library/Sounds/" + target + ".aiff"
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("không thấy file / âm %q", target)
	}
	cmd := exec.CommandContext(ctx, "afplay", path)
	if !wait {
		if err := cmd.Start(); err != nil {
			return err
		}
		go cmd.Wait()
		return nil
	}
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("không phát được %q: %w", target, err)
	}
	return nil
}
