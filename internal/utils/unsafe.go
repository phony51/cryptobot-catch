package utils

import (
	"unsafe"
)

func UnsafeStringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		int
	}{s, len(s)}))
}
