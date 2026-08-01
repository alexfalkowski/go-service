package test

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Encoder contains the real encoders exercised by config, cache, and transport tests.
var Encoder = encoding.NewMap()

// StreamEncoder contains the real streaming encoders/decoders exercised by content and transport tests.
var StreamEncoder = stream.NewMap()

// Content is the shared unary HTTP content registry backed by Encoder.
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

// Unencodable is a zero-field type whose MarshalJSON always fails with ErrFailed, letting tests force
// a JSON encode failure without a real payload shape.
type Unencodable struct{}

// MarshalJSON returns ErrFailed so encoding/json's Encoder.Encode fails.
func (Unencodable) MarshalJSON() ([]byte, error) {
	return nil, ErrFailed
}

// OpaqueErrorDecoder is a [stream.Decoder] test double that reads R one byte at a time until a Read
// error occurs, then returns a new error carrying only the original error's message — discarding its
// type the way the standard library JSON decoder does under GOEXPERIMENT=jsonv2 when a mid-scan Read
// fails. Used to prove a caller recovers a read-time sentinel error directly, rather than relying on
// the decoder to return that error unwrapped.
type OpaqueErrorDecoder struct {
	R io.Reader
}

// Decode reads from R one byte at a time until a Read error occurs, then returns a new, differently
// typed error carrying the same message.
func (d *OpaqueErrorDecoder) Decode(_ any) error {
	buf := make([]byte, 1)

	for {
		if _, err := d.R.Read(buf); err != nil {
			return errors.New(err.Error())
		}
	}
}

// Close is a no-op.
func (d *OpaqueErrorDecoder) Close() error {
	return nil
}

// TripleReadDecoder is a [stream.Decoder] test double that reads R one byte at a time exactly three
// times per Decode call, regardless of intervening errors, exercising a
// [github.com/alexfalkowski/go-service/v2/net/http/budget.Reader]'s sticky repeat-Read guard: a Read
// call made after the reader has already latched a size-limit error from a previous call.
type TripleReadDecoder struct {
	R io.Reader
}

// Decode reads from R exactly three times, returning the last call's error.
func (d *TripleReadDecoder) Decode(_ any) error {
	buf := make([]byte, 1)

	var lastErr error

	for range 3 {
		_, lastErr = d.R.Read(buf)
	}

	return lastErr
}

// Close is a no-op.
func (d *TripleReadDecoder) Close() error {
	return nil
}

// SingleReadDecoder is a [stream.Decoder] test double that performs exactly one large Read against R
// and always reports a successful decode, regardless of how many bytes that Read actually consumed.
// A real streaming decoder's internal buffering strategy can vary by Go toolchain (for example under
// GOEXPERIMENT=jsonv2), which can make it read incrementally and trip a
// [github.com/alexfalkowski/go-service/v2/net/http/budget.Reader]'s live per-Read cap check before
// Decode ever returns success. SingleReadDecoder bypasses that variability so a
// [github.com/alexfalkowski/go-service/v2/net/http/budget.Reader] post-decode size check can be
// exercised deterministically.
type SingleReadDecoder struct {
	R io.Reader
}

// Decode performs one Read against R with a generous buffer and always reports success.
func (d *SingleReadDecoder) Decode(_ any) error {
	buf := make([]byte, 4096)
	_, _ = d.R.Read(buf)

	return nil
}

// Close is a no-op.
func (d *SingleReadDecoder) Close() error {
	return nil
}

// NoBufferedDecoder is a [stream.Decoder] test double that always decodes successfully and has no
// Buffered method, exercising the no-correction fallback branch of a stream package's bufferedLen
// helper.
type NoBufferedDecoder struct{}

// Decode always succeeds without consuming v.
func (NoBufferedDecoder) Decode(_ any) error {
	return nil
}

// Close is a no-op.
func (NoBufferedDecoder) Close() error {
	return nil
}

// WrongBufferedDecoder is a [stream.Decoder] test double that always decodes successfully and whose
// Buffered method returns a reader that is not a *bytes.Reader, exercising the other no-correction
// fallback branch of a stream package's bufferedLen helper.
type WrongBufferedDecoder struct{}

// Decode always succeeds without consuming v.
func (WrongBufferedDecoder) Decode(_ any) error {
	return nil
}

// Close is a no-op.
func (WrongBufferedDecoder) Close() error {
	return nil
}

// Buffered returns an empty reader that is not a *bytes.Reader.
func (WrongBufferedDecoder) Buffered() io.Reader {
	return strings.NewReader("")
}

// CloseErrEncoder is a [stream.Encoder] test double whose Encode always succeeds and whose Close
// always fails with ErrFailed.
type CloseErrEncoder struct{}

// Encode always succeeds without writing v.
func (CloseErrEncoder) Encode(_ any) error {
	return nil
}

// Close returns ErrFailed.
func (CloseErrEncoder) Close() error {
	return ErrFailed
}
