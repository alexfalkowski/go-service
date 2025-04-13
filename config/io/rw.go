package io

import (
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
)

// NewReadWriter based on kind.
func NewReadWriter(name env.Name, kind, location string, fs os.FileSystem) ReaderWriter {
	var rw ReaderWriter

	switch kind {
	case "file":
		rw = NewFile(location, fs)
	case "env":
		rw = NewENV(location, fs)
	}

	return NewCommon(name, fs, rw)
}

// ReaderWriter for io.
type ReaderWriter interface {
	// Valid defines whether the read writer has a place to read or write.
	Valid() bool

	// Read bytes.
	Read() ([]byte, error)

	// Write bytes with file's mode.
	Write(data []byte, mode os.FileMode) error

	// Kind of read writer.
	Kind() string
}
