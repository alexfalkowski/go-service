package io

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewReadWriter based on kind.
func NewReadWriter(name env.Name, kind, location string, fs *os.FS) ReaderWriter {
	switch kind {
	case "file":
		return NewFile(location, fs)
	case "env":
		return NewENV(location, fs)
	default:
		return NewCommon(name, fs)
	}
}

// ReaderWriter for io.
type ReaderWriter interface {
	// Read bytes.
	Read() ([]byte, error)

	// Write bytes with file's mode.
	Write(data []byte, mode os.FileMode) error

	// Kind of read writer.
	Kind() string
}
