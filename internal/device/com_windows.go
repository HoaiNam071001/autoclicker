package device

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// Gọi COM / WinRT thẳng qua vtable, không cgo (MinGW / máy build CI thiếu header OCR; một kiểu cho cả UIA). Thứ tự
// method lấy từ header MinGW (UIA) / metadata Windows (WinRT – đối chiếu bằng reflection). Slot 0–2 IUnknown.
// Mọi lời gọi chạy trên một luồng OS riêng ở MTA (onCOM): apartment gắn với luồng, đối tượng giữ sẵn (engine OCR,
// IUIAutomation) chỉ dùng trên luồng đó.

var (
	procRoInitialize  = syscall.NewLazyDLL("combase.dll").NewProc("RoInitialize")
	procCoCreateInst  = syscall.NewLazyDLL("ole32.dll").NewProc("CoCreateInstance")
	oleaut32          = syscall.NewLazyDLL("oleaut32.dll")
	procSysStringLen  = oleaut32.NewProc("SysStringLen")
	procSysFreeString = oleaut32.NewProc("SysFreeString")
	procVariantClear  = oleaut32.NewProc("VariantClear")
	procRtlMoveMemory = syscall.NewLazyDLL("kernel32.dll").NewProc("RtlMoveMemory")
	comOnce           sync.Once
	comCalls          chan func()
)

// onCOM chạy fn trên luồng COM (lần đầu dựng luồng, RoInitialize MTA – dùng được cho cả COM thường) rồi chờ xong.
func onCOM(fn func()) {
	comOnce.Do(func() {
		comCalls = make(chan func())
		go func() {
			runtime.LockOSThread() // không thả luồng này – apartment gắn với luồng
			procRoInitialize.Call(1)
			for f := range comCalls {
				f()
			}
		}()
	})
	done := make(chan struct{})
	comCalls <- func() { defer close(done); fn() }
	<-done
}

// comObj – con trỏ đối tượng COM: trường đầu là vtable.
type comObj struct{ vtbl *[128]uintptr }

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
// viết ngay trong lời gọi vcall thì x được đưa lên heap, không bị dời khi stack lớn lên (như syscall.Proc.Call).
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

// coCreate – đối tượng COM trong tiến trình (CLSCTX_INPROC_SERVER).
func coCreate(clsid, iid *comGUID) (*comObj, error) {
	var o *comObj
	if r, _, _ := procCoCreateInst.Call(uintptr(unsafe.Pointer(clsid)), 0, 1, uintptr(unsafe.Pointer(iid)), uintptr(unsafe.Pointer(&o))); int32(r) < 0 {
		return nil, fmt.Errorf("CoCreateInstance: 0x%08X", uint32(r))
	}
	return o, nil
}

// utf16At chép n ký tự UTF-16 ở địa chỉ p (bộ nhớ của Windows) – qua RtlMoveMemory để không đổi uintptr → unsafe.Pointer (go vet).
func utf16At(p uintptr, n int) string {
	if p == 0 || n <= 0 {
		return ""
	}
	buf := make([]uint16, n)
	procRtlMoveMemory.Call(uintptr(unsafe.Pointer(&buf[0])), p, uintptr(n)*2)
	return string(utf16.Decode(buf))
}

// bstr đọc BSTR (không giải phóng – nơi gọi tự SysFreeString / VariantClear).
func bstr(p uintptr) string {
	n, _, _ := procSysStringLen.Call(p)
	return utf16At(p, int(n))
}

// variant – VARIANT x64 (24 byte): kiểu ở 2 byte đầu, giá trị từ byte 8.
type variant struct {
	vt  uint16
	_   [3]uint16
	val uintptr
	_   uintptr
}

const (
	vtI4   = 3
	vtBSTR = 8
	vtBool = 11
)

func (v *variant) clear() { procVariantClear.Call(uintptr(unsafe.Pointer(v))) }
