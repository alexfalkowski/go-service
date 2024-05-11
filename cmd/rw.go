package cmd

import "io/fs"

// NewReadWriter for cmd.
func NewReadWriter(kind, location string) ReaderWriter {
	switch kind {
	case "file":
		return NewFile(location)
	case "env":
		return NewENV(location)
	default:
		return NewNone()
	}
}

// ReaderWriter for cmd.
type ReaderWriter interface {
	// Read bytes.
	Read() ([]byte, error)

	// Write bytes with files's mode.
	Write(data []byte, mode fs.FileMode) error

	// Kind of read writer.
	Kind() string
}
