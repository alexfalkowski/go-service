package cmd

import (
	"github.com/alexfalkowski/go-service/os"
)

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

// ReaderWriter for cmd.
type ReaderWriter interface {
	// Read bytes.
	Read() (string, error)

	// Write bytes with file's mode.
	Write(data string, mode os.FileMode) error

	// Kind of read writer.
	Kind() string
}
