package bytes

import (
	"bytes"
	"unsafe"
)

// Buffer is an alias for bytes.Buffer.
type Buffer = bytes.Buffer

// NewBuffer is an alias for bytes.NewBuffer.
func NewBuffer(buf []byte) *bytes.Buffer {
	return bytes.NewBuffer(buf)
}

// NewBufferString is an alias for bytes.NewBufferString.
func NewBufferString(s string) *bytes.Buffer {
	return bytes.NewBufferString(s)
}

// NewReader is an alias for bytes.NewReader.
func NewReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}

// TrimSpace is an alias for bytes.TrimSpace.
func TrimSpace(s []byte) []byte {
	return bytes.TrimSpace(s)
}

// Clone is an alias for bytes.Clone.
func Clone(b []byte) []byte {
	return bytes.Clone(b)
}

// String converts b to a string without copying.
//
// The returned string aliases b and is only valid while b's contents remain unchanged.
func String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
