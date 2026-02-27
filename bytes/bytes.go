package bytes

import (
	"bytes"
	"unsafe"
)

// Buffer is an alias for bytes.Buffer.
//
// It is provided so go-service code can depend on a consistent import path while still using the
// standard library implementation.
type Buffer = bytes.Buffer

// NewBuffer returns a new Buffer initialized with buf's contents.
//
// This is a thin wrapper around bytes.NewBuffer. The returned buffer uses buf as its initial
// contents; subsequent writes may grow the buffer.
func NewBuffer(buf []byte) *bytes.Buffer {
	return bytes.NewBuffer(buf)
}

// NewBufferString returns a new Buffer initialized with the contents of s.
//
// This is a thin wrapper around bytes.NewBufferString.
func NewBufferString(s string) *bytes.Buffer {
	return bytes.NewBufferString(s)
}

// NewReader returns a new bytes.Reader reading from b.
//
// This is a thin wrapper around bytes.NewReader. The returned reader reads from b without copying.
func NewReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}

// TrimSpace returns a subslice of s with all leading and trailing white space removed.
//
// This is a thin wrapper around bytes.TrimSpace. The returned slice may refer to the same
// underlying array as s.
func TrimSpace(s []byte) []byte {
	return bytes.TrimSpace(s)
}

// Clone returns a copy of b.
//
// This is a thin wrapper around bytes.Clone. Unlike helpers that may return subslices, Clone always
// allocates a new slice (unless b is nil, in which case it returns nil).
func Clone(b []byte) []byte {
	return bytes.Clone(b)
}

// String converts b to a string without copying.
//
// # Safety and lifetime
//
// The returned string aliases the memory backing b. That means:
//
//   - You MUST NOT modify the contents of b after calling String(b).
//   - The returned string is only valid while b remains alive (reachable) and its backing array is not reused.
//
// Violating these constraints can lead to surprising behavior, data races, or memory safety issues.
// Use this helper only when you control the lifecycle of b and need to avoid an allocation.
// If you need an owning copy, use `string(b)` instead.
func String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
