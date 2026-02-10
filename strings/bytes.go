package strings

import "unsafe"

// Bytes converts s to a []byte without allocating.
//
// The returned slice aliases the memory backing s. It must be treated as read-only,
// and it must not be retained beyond the lifetime of s.
func Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
