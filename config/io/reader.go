package io

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewReader based on kind.
func NewReader(name env.Name, kind, location string, fs *os.FS) Reader {
	switch kind {
	case "file":
		return NewFile(location, fs)
	case "env":
		return NewENV(location)
	default:
		return NewCommon(name, fs)
	}
}

// Reader for io.
type Reader interface {
	// Read bytes.
	Read() ([]byte, error)

	// Kind of read writer.
	Kind() string
}
