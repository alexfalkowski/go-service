package test

import (
	"io"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/net/http/content"
)

// Encoder for tests.
var Encoder = encoding.NewMap()

// Content for tests.
var Content = content.NewContent(Encoder)

// NewEncoder for test.
func NewEncoder(err error) encoding.Encoder {
	return &enc{err: err}
}

type enc struct {
	err error
}

func (e *enc) Encode(_ io.Writer, _ any) error {
	return e.err
}

func (e *enc) Decode(_ io.Reader, _ any) error {
	return e.err
}
