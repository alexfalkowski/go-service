package bytes

import (
	"bytes"
	"unsafe"
)

// Buffer is an alias bytes.Buffer.
type Buffer = bytes.Buffer

var (
	// NewBuffer is an alias bytes.NewBuffer.
	NewBuffer = bytes.NewBuffer

	// NewBufferString is an alias bytes.NewBufferString.
	NewBufferString = bytes.NewBufferString

	// NewReader is an alias bytes.NewReader.
	NewReader = bytes.NewReader

	// TrimSpace is an alias bytes.TrimSpace.
	TrimSpace = bytes.TrimSpace
)

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
