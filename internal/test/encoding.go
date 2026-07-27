package test

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
)

// Encoder contains the real encoders exercised by config, cache, and transport tests.
var Encoder = encoding.NewMap()

// Content is the shared HTTP content registry backed by Encoder.
var Content = content.NewContent(Encoder, Pool)

// NewEncoder returns an encoder test double whose Encode and Decode methods fail with the supplied error.
func NewEncoder(err error) encoding.Encoder {
	return &enc{err: err}
}

type enc struct {
	err error
}

// Encode implements [encoding.Encoder] and returns the configured error.
func (e *enc) Encode(_ io.Writer, _ any) error {
	return e.err
}

// Decode implements [encoding.Encoder] and returns the configured error.
func (e *enc) Decode(_ io.Reader, _ any) error {
	return e.err
}

// PartialEncoder writes a partial payload and then fails encoding.
type PartialEncoder struct{}

// Encode writes a partial payload and returns ErrFailed.
func (PartialEncoder) Encode(w io.Writer, _ any) error {
	_, _ = io.WriteString(w, "partial")
	return ErrFailed
}

// Decode implements [encoding.Encoder] and always succeeds.
func (PartialEncoder) Decode(io.Reader, any) error {
	return nil
}
