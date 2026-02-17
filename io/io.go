package io

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/bytes"
)

// Reader is an alias for io.Reader.
type Reader = io.Reader

// ReadCloser is an alias for io.ReadCloser.
type ReadCloser = io.ReadCloser

// NopCloser is an alias for io.NopCloser.
func NopCloser(r Reader) ReadCloser {
	return io.NopCloser(r)
}

// ReadAll reads all bytes from r and returns the data along with a fresh ReadCloser over that data.
//
// The returned ReadCloser allows the caller to re-read the captured bytes without re-reading from r.
func ReadAll(r io.Reader) ([]byte, io.ReadCloser, error) {
	data, err := io.ReadAll(r)
	return data, NopCloser(bytes.NewReader(data)), err
}
