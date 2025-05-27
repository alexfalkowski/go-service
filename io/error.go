package io

import "io"

// NewErrReadCloser for io.
func NewErrReadCloser(err error) io.ReadCloser {
	return io.NopCloser(NewErrReader(err))
}

// NewErrReader for io.
func NewErrReader(err error) io.Reader {
	return &ErrReader{err: err}
}

// ErrReader for io.
type ErrReader struct {
	err error
}

// Read returns the error provided.
func (r *ErrReader) Read(_ []byte) (int, error) {
	return 0, r.err
}
