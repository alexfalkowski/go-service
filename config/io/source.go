package io

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewSource based on kind.
func NewSource(name env.Name, kind, location string, fs *os.FS) Source {
	switch kind {
	case "file":
		return NewFile(location, fs)
	case "env":
		return NewENV(location)
	default:
		return NewCommon(name, fs)
	}
}

// Source for io.
type Source interface {
	// Reader for source.
	Reader() io.ReadCloser

	// Kind of source.
	Kind() string
}
