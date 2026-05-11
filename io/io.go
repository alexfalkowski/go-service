package io

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/bytes"
)

// EOF aliases io.EOF.
var EOF = io.EOF

// WriterTo is an alias for io.WriterTo.
//
// It is provided so go-service code can depend on a consistent import path while
// preserving standard library semantics.
type WriterTo = io.WriterTo

// Writer is an alias for io.Writer.
//
// It is provided so go-service code can depend on a consistent import path while
// preserving standard library semantics.
type Writer = io.Writer

// ReaderFrom is an alias for io.ReaderFrom.
//
// It is provided so go-service code can depend on a consistent import path while
// preserving standard library semantics.
type ReaderFrom = io.ReaderFrom

// Reader is an alias for io.Reader.
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Reader = io.Reader

// ReadCloser is an alias for io.ReadCloser.
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type ReadCloser = io.ReadCloser

// Resetter is implemented by values that can reset themselves to an empty state.
//
// It is useful for generic code that wants to repopulate reusable buffers or
// readers without appending to their previous contents.
type Resetter interface {
	Reset()
}

// NopCloser returns a ReadCloser with a no-op Close method wrapping r.
//
// This is a thin wrapper around io.NopCloser.
func NopCloser(r Reader) ReadCloser {
	return io.NopCloser(r)
}

// LimitReader returns a Reader that reads from r but stops with EOF after n bytes.
//
// This is a thin wrapper around io.LimitReader.
func LimitReader(r Reader, n int64) Reader {
	return io.LimitReader(r, n)
}

// ReadAll reads all remaining bytes from r and returns:
//
//   - the captured bytes, and
//   - a fresh ReadCloser that reads from those captured bytes.
//
// This is useful when you need to consume a stream but also need to re-read the same content later
// (for example for logging, retries, signatures, or decoding twice).
//
// Memory note: ReadAll loads the entire stream into memory. It should only be used when the input
// size is bounded or otherwise acceptable for buffering.
//
// The returned ReadCloser is independent of the original reader; closing it does not affect r.
func ReadAll(r Reader) ([]byte, ReadCloser, error) {
	data, err := io.ReadAll(r)
	return data, NopCloser(bytes.NewReader(data)), err
}
