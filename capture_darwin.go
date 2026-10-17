//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework CoreGraphics -framework Foundation -framework ScreenCaptureKit
#include <CoreGraphics/CoreGraphics.h>
#import <Foundation/Foundation.h>
#import <ScreenCaptureKit/ScreenCaptureKit.h>
#include <stdlib.h>
#include <unistd.h>

static int preflight() { return CGPreflightScreenCaptureAccess() ? 1 : 0; }
static int request()   { return CGRequestScreenCaptureAccess() ? 1 : 0; }

// Chụp vùng (x, y, w, h – point, toạ độ toàn cục, nằm trong một màn) bằng ScreenCaptureKit, bỏ mọi cửa sổ của chính app (thanh mini,
// overlay) và không vẽ con trỏ chuột. Trả 0 nếu được; *out: RGBA (malloc, Go giải phóng).
static int capture_without_self(double x, double y, double w, double h, unsigned char **out, int *ow, int *oh) {
	if (@available(macOS 14.0, *)) {
		__block int rc = 1;
		__block CGImageRef img = NULL;
		dispatch_semaphore_t sem = dispatch_semaphore_create(0);
		[SCShareableContent getShareableContentExcludingDesktopWindows:NO onScreenWindowsOnly:YES
			completionHandler:^(SCShareableContent *content, NSError *err) {
				if (err || !content) { rc = 2; dispatch_semaphore_signal(sem); return; }
				// màn chứa tâm vùng; sourceRect tính từ góc màn đó
				SCDisplay *disp = nil;
				CGPoint mid = CGPointMake(x + w / 2, y + h / 2);
				for (SCDisplay *d in content.displays) if (CGRectContainsPoint(CGDisplayBounds(d.displayID), mid)) disp = d;
				if (!disp) { rc = 3; dispatch_semaphore_signal(sem); return; }
				CGRect db = CGDisplayBounds(disp.displayID);
				NSMutableArray *mine = [NSMutableArray array];
				for (SCRunningApplication *a in content.applications) if (a.processID == getpid()) [mine addObject:a];
				SCContentFilter *f = [[SCContentFilter alloc] initWithDisplay:disp excludingApplications:mine exceptingWindows:@[]];
				SCStreamConfiguration *c = [[SCStreamConfiguration alloc] init];
				CGFloat k = f.pointPixelScale;
				c.sourceRect = CGRectMake(x - db.origin.x, y - db.origin.y, w, h);
				c.width = (size_t)(w * k);
				c.height = (size_t)(h * k);
				c.showsCursor = NO;
				c.captureResolution = SCCaptureResolutionBest;
				[SCScreenshotManager captureImageWithFilter:f configuration:c completionHandler:^(CGImageRef im, NSError *e2) {
					if (im && !e2) { img = CGImageRetain(im); rc = 0; } else { rc = 4; }
					dispatch_semaphore_signal(sem);
				}];
			}];
		if (dispatch_semaphore_wait(sem, dispatch_time(DISPATCH_TIME_NOW, 5 * NSEC_PER_SEC)) != 0) return 6;
		if (rc != 0 || !img) return rc ? rc : 5;
		size_t W = CGImageGetWidth(img), H = CGImageGetHeight(img);
		unsigned char *buf = malloc(W * H * 4);
		CGColorSpaceRef cs = CGColorSpaceCreateWithName(kCGColorSpaceSRGB);
		CGContextRef ctx = CGBitmapContextCreate(buf, W, H, 8, W * 4, cs, kCGImageAlphaPremultipliedLast | kCGBitmapByteOrder32Big);
		CGContextDrawImage(ctx, CGRectMake(0, 0, W, H), img);
		CGContextRelease(ctx);
		CGColorSpaceRelease(cs);
		CGImageRelease(img);
		*out = buf; *ow = (int)W; *oh = (int)H;
		return 0;
	}
	return 10; // macOS < 14: không có SCScreenshotManager
}

static int has_sck() {
	if (@available(macOS 14.0, *)) return 1;
	return 0;
}
*/
import "C"

import (
	"errors"
	"image"
	"unsafe"
)

// macOS: chụp màn hình cần quyền Screen Recording (System Settings → Privacy & Security → Screen Recording).
// Thiếu quyền thì ảnh chụp chỉ có hình nền, không có cửa sổ – báo lỗi rõ thay vì tìm mãi không thấy.
func screenCaptureAllowed() bool { return C.preflight() == 1 }

func requestScreenCapture() bool { return C.request() == 1 }

func screenCaptureCheck() error {
	if !screenCaptureAllowed() {
		return errors.New("chưa cấp quyền ghi màn hình (Screen Recording) – System Settings → Privacy & Security → Screen Recording, bật cho app / terminal rồi mở lại")
	}
	return nil
}

// captureNoSelf chụp vùng r không có cửa sổ của app và con trỏ chuột. ok = false: hệ không hỗ trợ (macOS < 14) –
// dùng robotgo (có cả cửa sổ app nên phải ẩn app trước khi chụp).
func captureNoSelf(r Rect) (img image.Image, ok bool, err error) {
	var buf *C.uchar
	var w, h C.int
	rc := C.capture_without_self(C.double(r.X), C.double(r.Y), C.double(r.W), C.double(r.H), &buf, &w, &h)
	if rc == 10 {
		return nil, false, nil
	}
	if rc != 0 {
		return nil, true, errors.New("không chụp được màn hình (ScreenCaptureKit)")
	}
	defer C.free(unsafe.Pointer(buf))
	out := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	copy(out.Pix, unsafe.Slice((*byte)(unsafe.Pointer(buf)), int(w)*int(h)*4))
	return out, true, nil
}

// excludeSelfFromCapture: macOS chụp bằng captureNoSelf nên không cần làm gì thêm.
func excludeSelfFromCapture() {}

// canCaptureWithoutSelf: chụp được mà không dính cửa sổ app (khỏi ẩn app khi khoanh ảnh / thử tìm).
func canCaptureWithoutSelf() bool { return C.has_sck() == 1 }
