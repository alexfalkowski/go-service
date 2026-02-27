package strings

import "unsafe"

// Bytes converts s to a []byte without allocating.
//
// This function uses unsafe to create a byte slice view over the string's
// backing storage. As a result:
//
//   - The returned slice aliases the same memory as s.
//   - The returned slice must be treated as read-only. Writing to it results in
//     undefined behavior.
//   - Do not retain the returned slice beyond the lifetime of s. In particular,
//     do not store it in long-lived structures or return it when s was derived
//     from a transient buffer.
//
// For a safe conversion that does not alias the input string, use []byte(s),
// which allocates a new slice.
func Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
