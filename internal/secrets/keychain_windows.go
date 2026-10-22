package secrets

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// Windows Credential Manager (CredWriteW / CredReadW / CredDeleteW), mục "Ghosthand/<tên>", loại generic. Giá trị lưu
// UTF-16 để trình quản lý của Windows đọc được. Viết nhưng chưa chạy thử trên máy thật.
var (
	advapi32     = syscall.NewLazyDLL("advapi32.dll")
	procCredW    = advapi32.NewProc("CredWriteW")
	procCredR    = advapi32.NewProc("CredReadW")
	procCredDel  = advapi32.NewProc("CredDeleteW")
	procCredFree = advapi32.NewProc("CredFree")
)

const (
	credTypeGeneric       = 1
	credPersistLocalMachn = 2
	errNotFound           = 1168 // ERROR_NOT_FOUND
	maxBlob               = 2560 // CRED_MAX_CREDENTIAL_BLOB_SIZE
)

type credential struct {
	Flags              uint32
	Type               uint32
	TargetName         *uint16
	Comment            *uint16
	LastWritten        syscall.Filetime
	CredentialBlobSize uint32
	CredentialBlob     *byte
	Persist            uint32
	AttributeCount     uint32
	Attributes         uintptr
	TargetAlias        *uint16
	UserName           *uint16
}

// System – Backend thật.
type System struct{}

func target(name string) *uint16 {
	p, _ := syscall.UTF16PtrFromString("Ghosthand/" + name)
	return p
}

func (System) Get(_ context.Context, name string) (string, error) {
	var cred *credential
	r, _, err := procCredR.Call(uintptr(unsafe.Pointer(target(name))), credTypeGeneric, 0, uintptr(unsafe.Pointer(&cred)))
	if r == 0 {
		if errors.Is(err, syscall.Errno(errNotFound)) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("Credential Manager đọc lỗi: %w", err)
	}
	defer procCredFree.Call(uintptr(unsafe.Pointer(cred)))
	blob := unsafe.Slice((*uint16)(unsafe.Pointer(cred.CredentialBlob)), cred.CredentialBlobSize/2)
	return string(utf16.Decode(blob)), nil
}

func (System) Set(_ context.Context, name, value string) error {
	blob := utf16.Encode([]rune(value))
	if len(blob)*2 > maxBlob {
		return fmt.Errorf("giá trị bí mật quá dài (tối đa %d ký tự)", maxBlob/2)
	}
	user, _ := syscall.UTF16PtrFromString("Ghosthand")
	c := credential{Type: credTypeGeneric, TargetName: target(name), UserName: user, Persist: credPersistLocalMachn,
		CredentialBlobSize: uint32(len(blob) * 2), CredentialBlob: (*byte)(unsafe.Pointer(&blob[0]))}
	if r, _, err := procCredW.Call(uintptr(unsafe.Pointer(&c)), 0); r == 0 {
		return fmt.Errorf("Credential Manager ghi lỗi: %w", err)
	}
	return nil
}

func (System) Delete(_ context.Context, name string) error {
	if r, _, err := procCredDel.Call(uintptr(unsafe.Pointer(target(name))), credTypeGeneric, 0); r == 0 && !errors.Is(err, syscall.Errno(errNotFound)) {
		return fmt.Errorf("Credential Manager xoá lỗi: %w", err)
	}
	return nil
}
