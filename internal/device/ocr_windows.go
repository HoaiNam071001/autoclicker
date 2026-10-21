package device

import (
	"errors"
	"image"
	"log"
	"math"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
	"unsafe"

	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// Đọc chữ trên Windows: mặc định Tesseract + gói vie (ocr_tess_windows.go); nạp không được thì dùng Windows.Media.Ocr
// (có sẵn từ Windows 10, không cài thêm). Gói nhận chữ của Windows theo ngôn ngữ hệ thống – không có tiếng Việt: chữ có
// dấu đọc sai dấu → dòng mang Approx (vision/fuzzy.go). Chạy trên luồng COM (onCOM); engine + factory dựng một lần,
// chỉ dùng trên luồng đó.

var (
	iidOcrEngineStatics      = parseGUID("5bffa85a-3384-3540-9940-699120d428a8")
	iidSoftwareBitmapFactory = parseGUID("c99feb69-2d62-4d47-a6b3-4fdb6a07fdf8")
	iidMemoryBuffer          = parseGUID("fbc4dd2a-245b-11e4-af98-689423260cf8")
	iidMemoryBufferByteAcc   = parseGUID("5b0d3235-4dba-4d44-865e-8f1d0e4fd04d")
	ocrEngine, ocrBmpFactory *comObj // chỉ đụng trên luồng COM
)

var tessFailLog sync.Once

func readText(src image.Image) ([]vision.TextLine, error) {
	lines, err := tessRead(src)
	if err == nil {
		return lines, nil
	}
	tessFailLog.Do(func() {
		log.Printf("OCR: không dùng được Tesseract, đọc bằng OCR Windows (sai dấu tiếng Việt): %v", err)
	})
	return readTextWinRT(src)
}

func readTextWinRT(src image.Image) (lines []vision.TextLine, err error) {
	b := src.Bounds()
	if b.Dx() == 0 || b.Dy() == 0 {
		return nil, nil
	}
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			img.Set(x, y, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	onCOM(func() {
		if ocrEngine == nil {
			if ocrEngine, err = newOcrEngine(); err != nil {
				return
			}
		}
		if ocrBmpFactory == nil {
			if ocrBmpFactory, err = activationFactory("Windows.Graphics.Imaging.SoftwareBitmap", iidSoftwareBitmapFactory); err != nil {
				return
			}
		}
		lines, err = recognize(ocrEngine, ocrBmpFactory, img)
	})
	return lines, err
}

// newOcrEngine: theo ngôn ngữ người dùng; không có gói nào hợp thì lấy gói nhận chữ đầu tiên đã cài.
func newOcrEngine() (*comObj, error) {
	st, err := activationFactory("Windows.Media.Ocr.OcrEngine", iidOcrEngineStatics)
	if err != nil {
		return nil, err
	}
	defer st.release()
	var eng *comObj
	if err := vcall(st, "TryCreateFromUserProfileLanguages", 10, uintptr(unsafe.Pointer(&eng))); err != nil {
		return nil, err
	}
	if eng != nil {
		return eng, nil
	}
	var langs *comObj
	if err := vcall(st, "AvailableRecognizerLanguages", 7, uintptr(unsafe.Pointer(&langs))); err != nil {
		return nil, err
	}
	defer langs.release()
	err = vectorView(langs, func(i int, lang *comObj) error {
		if eng != nil {
			return nil
		}
		return vcall(st, "TryCreateFromLanguage", 9, uintptr(unsafe.Pointer(lang)), uintptr(unsafe.Pointer(&eng)))
	})
	if err == nil && eng == nil {
		err = errors.New("Windows chưa có gói nhận chữ (OCR) – Settings → Time & language → Language & region → thêm English (United States)")
	}
	return eng, err
}

// recognize chép ảnh vào SoftwareBitmap (BGRA) rồi đọc chữ; khung theo px ảnh.
func recognize(eng, factory *comObj, img *image.RGBA) ([]vision.TextLine, error) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	var bmp *comObj
	if err := vcall(factory, "SoftwareBitmap.Create", 6, 87 /* Bgra8 */, uintptr(w), uintptr(h), uintptr(unsafe.Pointer(&bmp))); err != nil {
		return nil, err
	}
	defer bmp.release()
	if err := fillBitmap(bmp, img); err != nil {
		return nil, err
	}
	var op *comObj
	if err := vcall(eng, "OcrEngine.RecognizeAsync", 6, uintptr(unsafe.Pointer(bmp)), uintptr(unsafe.Pointer(&op))); err != nil {
		return nil, err
	}
	defer op.release()
	res, err := awaitOp(op, 20*time.Second)
	if err != nil {
		return nil, err
	}
	defer res.release()
	var ls *comObj
	if err := vcall(res, "OcrResult.Lines", 6, uintptr(unsafe.Pointer(&ls))); err != nil {
		return nil, err
	}
	defer ls.release()
	var lines []vision.TextLine
	err = vectorView(ls, func(_ int, l *comObj) error {
		line, err := readLine(l)
		if err == nil && line.Text != "" {
			lines = append(lines, line)
		}
		return err
	})
	return lines, err
}

// fillBitmap chép điểm ảnh RGBA → BGRA qua IMemoryBufferByteAccess (vùng nhớ của bitmap, có stride riêng).
func fillBitmap(bmp *comObj, img *image.RGBA) error {
	var buf *comObj
	if err := vcall(bmp, "SoftwareBitmap.LockBuffer", 15, 2 /* Write */, uintptr(unsafe.Pointer(&buf))); err != nil {
		return err
	}
	defer buf.release()
	var plane struct{ Start, W, H, Stride int32 }
	if err := vcall(buf, "BitmapBuffer.GetPlaneDescription", 7, 0, uintptr(unsafe.Pointer(&plane))); err != nil {
		return err
	}
	mem, err := buf.query(iidMemoryBuffer)
	if err != nil {
		return err
	}
	defer mem.release()
	var ref *comObj
	if err := vcall(mem, "MemoryBuffer.CreateReference", 6, uintptr(unsafe.Pointer(&ref))); err != nil {
		return err
	}
	defer ref.release()
	acc, err := ref.query(iidMemoryBufferByteAcc)
	if err != nil {
		return err
	}
	defer acc.release()
	var p *byte
	var n uint32
	if err := vcall(acc, "IMemoryBufferByteAccess.GetBuffer", 3, uintptr(unsafe.Pointer(&p)), uintptr(unsafe.Pointer(&n))); err != nil {
		return err
	}
	dst := unsafe.Slice(p, n)
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	for y := 0; y < h; y++ {
		row, src := dst[int(plane.Start)+y*int(plane.Stride):], img.Pix[y*img.Stride:]
		for x := 0; x < w; x++ {
			row[4*x], row[4*x+1], row[4*x+2], row[4*x+3] = src[4*x+2], src[4*x+1], src[4*x], 255
		}
	}
	return nil
}

// readLine ghép các từ của một dòng (cách bằng dấu cách) – vị trí từ theo UTF-16 như vision.TextWord; khung dòng = hợp khung từ.
func readLine(l *comObj) (vision.TextLine, error) {
	var ws *comObj
	if err := vcall(l, "OcrLine.Words", 6, uintptr(unsafe.Pointer(&ws))); err != nil {
		return vision.TextLine{}, err
	}
	defer ws.release()
	line := vision.TextLine{Approx: true} // gói nhận chữ không có tiếng Việt – xem vision/fuzzy.go
	var sb strings.Builder
	pos := 0
	err := vectorView(ws, func(i int, w *comObj) error {
		var r struct{ X, Y, W, H float32 }
		if err := vcall(w, "OcrWord.BoundingRect", 6, uintptr(unsafe.Pointer(&r))); err != nil {
			return err
		}
		var hs uintptr
		if err := vcall(w, "OcrWord.Text", 7, uintptr(unsafe.Pointer(&hs))); err != nil {
			return err
		}
		text := hstring(hs)
		if i > 0 {
			sb.WriteByte(' ')
			pos++
		}
		n := len(utf16.Encode([]rune(text)))
		box := model.Rect{X: int(math.Round(float64(r.X))), Y: int(math.Round(float64(r.Y))),
			W: max(int(math.Round(float64(r.W))), 1), H: max(int(math.Round(float64(r.H))), 1)}
		line.Words = append(line.Words, vision.TextWord{Start: pos, Len: n, Box: box})
		line.Box = vision.Union(line.Box, box)
		sb.WriteString(text)
		pos += n
		return nil
	})
	line.Text = sb.String()
	return line, err
}
