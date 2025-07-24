package bytes

import (
	"bytes"
	"unsafe"
)

// Buffer is an alias bytes.Buffer.
type Buffer = bytes.Buffer

// NewBuffer is an alias bytes.NewBuffer.
func NewBuffer(buf []byte) *bytes.Buffer {
	return bytes.NewBuffer(buf)
}

// NewBufferString is an alias bytes.NewBufferString.
func NewBufferString(s string) *bytes.Buffer {
	return bytes.NewBufferString(s)
}

// NewReader is an alias bytes.NewReader.
func NewReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}

// TrimSpace is an alias bytes.TrimSpace.
func TrimSpace(s []byte) []byte {
	return bytes.TrimSpace(s)
}

// Copy bytes to a new slice.
func Copy(b []byte) []byte {
	bytes := make([]byte, len(b))
	copy(bytes, b)
	return bytes
}

// String from the bytes.
func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
