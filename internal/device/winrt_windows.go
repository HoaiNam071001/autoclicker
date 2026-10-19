package device

import (
	"fmt"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

// Phần riêng WinRT (OCR) trên nền com_windows.go: HSTRING, activation factory, chờ IAsyncOperation, duyệt
// IVectorView. Slot 3–5 của đối tượng WinRT là IInspectable, method riêng bắt đầu từ 6.

var (
	combase                       = syscall.NewLazyDLL("combase.dll")
	procRoGetActivationFactory    = combase.NewProc("RoGetActivationFactory")
	procWindowsCreateString       = combase.NewProc("WindowsCreateString")
	procWindowsDeleteString       = combase.NewProc("WindowsDeleteString")
	procWindowsGetStringRawBuffer = combase.NewProc("WindowsGetStringRawBuffer")
)

func newHString(s string) (uintptr, error) {
	u := utf16.Encode([]rune(s))
	var h uintptr
	if r, _, _ := procWindowsCreateString.Call(uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)), uintptr(unsafe.Pointer(&h))); int32(r) < 0 {
		return 0, fmt.Errorf("WindowsCreateString: 0x%08X", uint32(r))
	}
	return h, nil
}

// hstring đọc rồi giải phóng HSTRING.
func hstring(h uintptr) string {
	if h == 0 {
		return ""
	}
	defer procWindowsDeleteString.Call(h)
	var n uint32
	p, _, _ := procWindowsGetStringRawBuffer.Call(h, uintptr(unsafe.Pointer(&n)))
	return utf16At(p, int(n))
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
