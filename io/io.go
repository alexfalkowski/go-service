package io

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/bytes"
)

type (
	// Reader is an alias for io.Reader.
	Reader = io.Reader

	// ReadCloser is an alias for io.ReadCloser.
	ReadCloser = io.ReadCloser
)

// NopCloser is an alias for io.NopCloser.
func NopCloser(r Reader) ReadCloser {
	return io.NopCloser(r)
}

// ReadAll reads all the bytes from the io.Reader and returns the bytes with an io.ReadCloser.
func ReadAll(r io.Reader) ([]byte, io.ReadCloser, error) {
	data, err := io.ReadAll(r)
	return data, NopCloser(bytes.NewReader(data)), err
}
