// +build !appengine

package fasttemplate

import (
	"unsafe"
)

func unsafeBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func unsafeString2Bytes(s string) []byte {
	temporaryData := (*[2]uintptr)(unsafe.Pointer(&s))
	bytesData := [3]uintptr{temporaryData[0], temporaryData[1], temporaryData[1]}
	return *(*[]byte)(unsafe.Pointer(&bytesData))
}
