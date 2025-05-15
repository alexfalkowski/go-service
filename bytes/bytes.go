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
func Copy(bytes []byte) []byte {
	newBytes := make([]byte, len(bytes))
	copy(newBytes, bytes)

	return newBytes
}

// String from the bytes.
func String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
