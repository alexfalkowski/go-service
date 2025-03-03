package cmd

import (
	"github.com/alexfalkowski/go-service/os"
)

// ReaderWriter for cmd.
type ReaderWriter interface {
	// Read bytes.
	Read() ([]byte, error)

	// Write bytes with file's mode.
	Write(data []byte, mode os.FileMode) error

	// Kind of read writer.
	Kind() string
}

// NewReadWriter for cmd.
func NewReadWriter(kind, location string, fs os.FileSystem) ReaderWriter {
	switch kind {
	case "file":
		return NewFile(location, fs)
	case "env":
		return NewENV(location, fs)
	default:
		return NewNone()
	}
}
