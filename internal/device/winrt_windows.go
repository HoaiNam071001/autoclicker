package device

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

// Gọi WinRT thẳng qua vtable (combase.dll), không cgo: MinGW không có header Windows.Media.Ocr, máy build CI cũng vậy.
// Thứ tự method theo metadata Windows (đã đối chiếu bằng reflection trên Windows 11). Slot 0–2 IUnknown, 3–5 IInspectable.

var (
	combase                       = syscall.NewLazyDLL("combase.dll")
	procRoInitialize              = combase.NewProc("RoInitialize")
	procRoGetActivationFactory    = combase.NewProc("RoGetActivationFactory")
	procWindowsCreateString       = combase.NewProc("WindowsCreateString")
	procWindowsDeleteString       = combase.NewProc("WindowsDeleteString")
	procWindowsGetStringRawBuffer = combase.NewProc("WindowsGetStringRawBuffer")
	procRtlMoveMemory             = syscall.NewLazyDLL("kernel32.dll").NewProc("RtlMoveMemory")
)

// comObj – con trỏ đối tượng COM: trường đầu là vtable.
type comObj struct{ vtbl *[64]uintptr }

type comGUID struct {
	d1     uint32
	d2, d3 uint16
	d4     [8]byte
}

// parseGUID("5a14bc41-5b76-3140-b680-8825562683ac") – chỉ dùng cho hằng trong code, sai định dạng thì panic lúc khởi động.
func parseGUID(s string) *comGUID {
	h := strings.ReplaceAll(s, "-", "")
	n := func(a, b int) uint64 {
		v, err := strconv.ParseUint(h[a:b], 16, 64)
		if err != nil || len(h) != 32 {
			panic("GUID sai: " + s)
		}
		return v
	}
	g := &comGUID{d1: uint32(n(0, 8)), d2: uint16(n(8, 12)), d3: uint16(n(12, 16))}
	for i := range g.d4 {
		g.d4[i] = byte(n(16+2*i, 18+2*i))
	}
	return g
}

// vcall gọi method ở slot; HRESULT âm → lỗi kèm tên để dò. uintptrescapes: tham số [out] dạng uintptr(unsafe.Pointer(&x))
// viết ngay trong lời gọi call thì x được đưa lên heap, không bị dời khi stack lớn lên (như syscall.Proc.Call).
//
//go:uintptrescapes
func vcall(o *comObj, name string, slot int, args ...uintptr) error {
	r, _, _ := syscall.SyscallN(o.vtbl[slot], append([]uintptr{uintptr(unsafe.Pointer(o))}, args...)...)
	if int32(r) < 0 {
		return fmt.Errorf("%s: HRESULT 0x%08X", name, uint32(r))
	}
	return nil
}

func (o *comObj) release() {
	if o != nil {
		syscall.SyscallN(o.vtbl[2], uintptr(unsafe.Pointer(o)))
	}
}

func (o *comObj) query(iid *comGUID) (*comObj, error) {
	var out *comObj
	err := vcall(o, "QueryInterface", 0, uintptr(unsafe.Pointer(iid)), uintptr(unsafe.Pointer(&out)))
	return out, err
}

// rtInit – luồng gọi WinRT phải ở apartment; MTA, gọi một lần trên luồng đã LockOSThread.
func rtInit() { procRoInitialize.Call(1) }

func newHString(s string) (uintptr, error) {
	u := utf16.Encode([]rune(s))
	var h uintptr
	if r, _, _ := procWindowsCreateString.Call(uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)), uintptr(unsafe.Pointer(&h))); int32(r) < 0 {
		return 0, fmt.Errorf("WindowsCreateString: 0x%08X", uint32(r))
	}
	return h, nil
}

// hstring đọc rồi giải phóng HSTRING. Chép qua RtlMoveMemory để không đổi uintptr → unsafe.Pointer (go vet).
func hstring(h uintptr) string {
	if h == 0 {
		return ""
	}
	defer procWindowsDeleteString.Call(h)
	var n uint32
	p, _, _ := procWindowsGetStringRawBuffer.Call(h, uintptr(unsafe.Pointer(&n)))
	if n == 0 {
		return ""
	}
	buf := make([]uint16, n)
	procRtlMoveMemory.Call(uintptr(unsafe.Pointer(&buf[0])), p, uintptr(n)*2)
	return string(utf16.Decode(buf))
}

// activationFactory – đối tượng static / factory của lớp WinRT theo interface iid.
func activationFactory(class string, iid *comGUID) (*comObj, error) {
	h, err := newHString(class)
	if err != nil {
		return nil, err
	}
	defer procWindowsDeleteString.Call(h)
	var f *comObj
	if r, _, _ := procRoGetActivationFactory.Call(h, uintptr(unsafe.Pointer(iid)), uintptr(unsafe.Pointer(&f))); int32(r) < 0 {
		return nil, fmt.Errorf("RoGetActivationFactory %s: 0x%08X", class, uint32(r))
	}
	return f, nil
}

var iidAsyncInfo = parseGUID("00000036-0000-0000-c000-000000000046")

// awaitOp chờ IAsyncOperation<T> xong (hỏi trạng thái – khỏi phải cài handler COM) rồi lấy kết quả.
func awaitOp(op *comObj, timeout time.Duration) (*comObj, error) {
	info, err := op.query(iidAsyncInfo)
	if err != nil {
		return nil, err
	}
	defer info.release()
	for end := time.Now().Add(timeout); ; time.Sleep(3 * time.Millisecond) {
		var status int32 // 0 đang chạy, 1 xong, 2 huỷ, 3 lỗi
		if err := vcall(info, "IAsyncInfo.Status", 7, uintptr(unsafe.Pointer(&status))); err != nil {
			return nil, err
		}
		if status == 1 {
			break
		}
		if status != 0 {
			var code int32
			vcall(info, "IAsyncInfo.ErrorCode", 8, uintptr(unsafe.Pointer(&code)))
			return nil, fmt.Errorf("tác vụ WinRT lỗi (trạng thái %d, 0x%08X)", status, uint32(code))
		}
		if time.Now().After(end) {
			vcall(info, "IAsyncInfo.Cancel", 9)
			return nil, fmt.Errorf("tác vụ WinRT quá %v", timeout)
		}
	}
	var res *comObj
	err = vcall(op, "IAsyncOperation.GetResults", 8, uintptr(unsafe.Pointer(&res)))
	return res, err
}

// vectorView duyệt IVectorView<T>: GetAt slot 6, Size slot 7. fn nhận phần tử (tự release sau khi fn trả về).
func vectorView(v *comObj, fn func(i int, item *comObj) error) error {
	var n uint32
	if err := vcall(v, "IVectorView.Size", 7, uintptr(unsafe.Pointer(&n))); err != nil {
		return err
	}
	for i := uint32(0); i < n; i++ {
		var item *comObj
		if err := vcall(v, "IVectorView.GetAt", 6, uintptr(i), uintptr(unsafe.Pointer(&item))); err != nil {
			return err
		}
		err := fn(int(i), item)
		item.release()
		if err != nil {
			return err
		}
	}
	return nil
}
