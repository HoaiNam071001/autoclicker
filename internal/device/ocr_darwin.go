//go:build darwin

package device

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework CoreGraphics -framework Foundation -framework Vision
#include <CoreGraphics/CoreGraphics.h>
#import <Foundation/Foundation.h>
#import <Vision/Vision.h>
#include <stdlib.h>
#include <string.h>

typedef struct { char *text; double x, y, w, h; } ocr_line;
typedef struct { int line, start, len; double x, y, w, h; } ocr_word;

// Đọc chữ trên ảnh RGBA (w × h) bằng Vision, tiếng Việt + tiếng Anh. Khung chuẩn hoá 0..1, gốc trên trái.
// Trả 0 nếu được; *lines / *words do C cấp phát (Go giải phóng bằng ocr_free).
static int ocr_rgba(unsigned char *pix, int w, int h, int stride, ocr_line **lines, int *nl, ocr_word **words, int *nw) {
	CGColorSpaceRef cs = CGColorSpaceCreateWithName(kCGColorSpaceSRGB);
	CGDataProviderRef dp = CGDataProviderCreateWithData(NULL, pix, (size_t)stride * h, NULL);
	CGImageRef img = CGImageCreate(w, h, 8, 32, stride, cs, kCGImageAlphaPremultipliedLast | kCGBitmapByteOrder32Big, dp, NULL, false, kCGRenderingIntentDefault);
	CGDataProviderRelease(dp);
	CGColorSpaceRelease(cs);
	if (!img) return 1;
	VNRecognizeTextRequest *req = [[VNRecognizeTextRequest alloc] init];
	req.recognitionLevel = VNRequestTextRecognitionLevelAccurate;
	req.recognitionLanguages = @[@"vi-VT", @"en-US"];
	req.usesLanguageCorrection = YES;
	VNImageRequestHandler *hd = [[VNImageRequestHandler alloc] initWithCGImage:img options:@{}];
	NSError *err = nil;
	BOOL ok = [hd performRequests:@[req] error:&err];
	CGImageRelease(img);
	if (!ok) return 2;
	NSArray<VNRecognizedTextObservation *> *obs = req.results;
	int n = (int)obs.count;
	__block int cap = 64, m = 0;
	ocr_line *ls = calloc(n > 0 ? n : 1, sizeof(ocr_line));
	__block ocr_word *ws = malloc(cap * sizeof(ocr_word));
	int k = 0;
	for (VNRecognizedTextObservation *o in obs) {
		VNRecognizedText *t = [o topCandidates:1].firstObject;
		if (!t) continue;
		CGRect b = o.boundingBox;
		ls[k] = (ocr_line){strdup(t.string.UTF8String), b.origin.x, 1 - b.origin.y - b.size.height, b.size.width, b.size.height};
		NSString *s = t.string;
		int line = k;
		[s enumerateSubstringsInRange:NSMakeRange(0, s.length) options:NSStringEnumerationByWords
			usingBlock:^(NSString *sub, NSRange r, NSRange er, BOOL *stop) {
				VNRectangleObservation *wb = [t boundingBoxForRange:r error:nil];
				if (!wb) return;
				CGRect bb = wb.boundingBox;
				if (m == cap) { cap *= 2; ocr_word *nws = realloc(ws, cap * sizeof(ocr_word)); if (!nws) { *stop = YES; return; } ws = nws; }
				ws[m++] = (ocr_word){line, (int)r.location, (int)r.length, bb.origin.x, 1 - bb.origin.y - bb.size.height, bb.size.width, bb.size.height};
			}];
		k++;
	}
	*lines = ls; *nl = k; *words = ws; *nw = m;
	return 0;
}

static void ocr_free(ocr_line *ls, int n, ocr_word *ws) {
	for (int i = 0; i < n; i++) free(ls[i].text);
	free(ls);
	free(ws);
}
*/
import "C"

import (
	"errors"
	"image"
	"image/draw"
	"math"
	"unsafe"

	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

func readText(src image.Image) ([]vision.TextLine, error) {
	b := src.Bounds()
	img, ok := src.(*image.RGBA)
	// Vision đọc kém ảnh quá dẹt (vùng 1 dòng kiểu thanh menu 1920×40) → chèn thêm nền (màu góc ảnh) ở phải / dưới
	// cho tỷ lệ ≤ 3:1. Gốc ảnh giữ nguyên nên khung chữ không phải đổi.
	pw, ph := max(b.Dx(), b.Dy()/3), max(b.Dy(), b.Dx()/3)
	if !ok || b.Min != (image.Point{}) || pw != b.Dx() || ph != b.Dy() {
		img = image.NewRGBA(image.Rect(0, 0, pw, ph))
		draw.Draw(img, img.Bounds(), image.NewUniform(src.At(b.Min.X, b.Min.Y)), image.Point{}, draw.Src)
		draw.Draw(img, image.Rect(0, 0, b.Dx(), b.Dy()), src, b.Min, draw.Src)
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w == 0 || h == 0 {
		return nil, nil
	}
	pix := C.CBytes(img.Pix) // Vision giữ con trỏ trong lúc đọc – không đưa thẳng bộ nhớ Go sang C
	defer C.free(pix)
	var ls *C.ocr_line
	var ws *C.ocr_word
	var nl, nw C.int
	if C.ocr_rgba((*C.uchar)(pix), C.int(w), C.int(h), C.int(img.Stride), &ls, &nl, &ws, &nw) != 0 {
		return nil, errors.New("không đọc được chữ trên màn hình (Vision)")
	}
	defer C.ocr_free(ls, nl, ws)
	px := func(x, y, bw, bh C.double) model.Rect {
		x0, y0 := int(math.Round(float64(x)*float64(w))), int(math.Round(float64(y)*float64(h)))
		return model.Rect{X: x0, Y: y0, W: max(int(math.Round(float64(bw)*float64(w))), 1), H: max(int(math.Round(float64(bh)*float64(h))), 1)}
	}
	lines := make([]vision.TextLine, int(nl))
	for i, l := range unsafe.Slice(ls, int(nl)) {
		lines[i] = vision.TextLine{Text: C.GoString(l.text), Box: px(l.x, l.y, l.w, l.h)}
	}
	if nw > 0 {
		for _, wd := range unsafe.Slice(ws, int(nw)) {
			li := int(wd.line)
			lines[li].Words = append(lines[li].Words, vision.TextWord{Start: int(wd.start), Len: int(wd.len), Box: px(wd.x, wd.y, wd.w, wd.h)})
		}
	}
	return lines, nil
}
