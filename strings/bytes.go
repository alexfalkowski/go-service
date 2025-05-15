package strings

import "unsafe"

type stringCap struct {
	string
	Cap int
}

// Bytes from string.
func Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&stringCap{s, len(s)}))
}
