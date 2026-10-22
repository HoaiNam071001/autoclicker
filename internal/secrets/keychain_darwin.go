package secrets

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Security -framework Foundation
#import <Foundation/Foundation.h>
#import <Security/Security.h>

static NSMutableDictionary *kc_query(const char *name) {
	return [@{
		(__bridge id)kSecClass: (__bridge id)kSecClassGenericPassword,
		(__bridge id)kSecAttrService: @"Ghosthand",
		(__bridge id)kSecAttrAccount: [NSString stringWithUTF8String:name],
	} mutableCopy];
}

static int kc_set(const char *name, const char *value) {
	NSData *data = [NSData dataWithBytes:value length:strlen(value)];
	OSStatus st = SecItemUpdate((__bridge CFDictionaryRef)kc_query(name), (__bridge CFDictionaryRef)@{(__bridge id)kSecValueData: data});
	if (st == errSecItemNotFound) {
		NSMutableDictionary *q = kc_query(name);
		q[(__bridge id)kSecValueData] = data;
		st = SecItemAdd((__bridge CFDictionaryRef)q, NULL);
	}
	return (int)st;
}

// Trả 0 và *out = chuỗi malloc (UTF-8), hoặc mã lỗi Security.
static int kc_get(const char *name, char **out) {
	NSMutableDictionary *q = kc_query(name);
	q[(__bridge id)kSecReturnData] = @YES;
	q[(__bridge id)kSecMatchLimit] = (__bridge id)kSecMatchLimitOne;
	CFTypeRef res = NULL;
	OSStatus st = SecItemCopyMatching((__bridge CFDictionaryRef)q, &res);
	if (st != errSecSuccess) return (int)st;
	NSData *data = (__bridge_transfer NSData *)res;
	*out = malloc(data.length + 1);
	memcpy(*out, data.bytes, data.length);
	(*out)[data.length] = 0;
	return 0;
}

static int kc_delete(const char *name) {
	return (int)SecItemDelete((__bridge CFDictionaryRef)kc_query(name));
}
*/
import "C"

import (
	"context"
	"fmt"
	"unsafe"
)

const errSecItemNotFound = -25300

// System – Backend thật: Keychain đăng nhập, mục "Ghosthand" / <tên>.
type System struct{}

func kcErr(op string, st C.int) error {
	return fmt.Errorf("Keychain %s lỗi (mã %d)", op, int(st))
}

func (System) Get(_ context.Context, name string) (string, error) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	var out *C.char
	switch st := C.kc_get(cn, &out); {
	case st == errSecItemNotFound:
		return "", ErrNotFound
	case st != 0:
		return "", kcErr("đọc", st)
	}
	defer C.free(unsafe.Pointer(out))
	return C.GoString(out), nil
}

func (System) Set(_ context.Context, name, value string) error {
	cn, cv := C.CString(name), C.CString(value)
	defer C.free(unsafe.Pointer(cn))
	defer C.free(unsafe.Pointer(cv))
	if st := C.kc_set(cn, cv); st != 0 {
		return kcErr("ghi", st)
	}
	return nil
}

func (System) Delete(_ context.Context, name string) error {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	if st := C.kc_delete(cn); st != 0 && st != errSecItemNotFound {
		return kcErr("xoá", st)
	}
	return nil
}
