package strings

import "unsafe"

// Bytes from string.
func Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
