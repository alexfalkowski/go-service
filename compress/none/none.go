package none

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress/errors"
)

// NewCompressor constructs a no-op compressor implementation.
//
// The returned value implements `github.com/alexfalkowski/go-service/v2/compress.Compressor` and is
// typically registered under kind "none" to represent "no compression".
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements a no-op compression codec.
//
// Compress returns the input unchanged, and Decompress returns the input unchanged with a nil error.
// This is useful when you want to disable compression while still satisfying the common compression
// interface.
type Compressor struct{}

// Compress returns data unchanged.
//
// An error is returned if data exceeds size.
func (c *Compressor) Compress(data []byte, size bytes.Size) ([]byte, error) {
	if int64(len(data)) > size.Bytes() {
		return nil, errors.ErrTooLarge
	}

	return data, nil
}

// Decompress returns data unchanged.
//
// An error is returned if data exceeds size.
func (c *Compressor) Decompress(data []byte, size bytes.Size) ([]byte, error) {
	if int64(len(data)) > size.Bytes() {
		return nil, errors.ErrTooLarge
	}

	return data, nil
}
