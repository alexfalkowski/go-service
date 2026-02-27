package snappy

import "github.com/klauspost/compress/snappy"

// NewCompressor constructs a Snappy compressor implementation.
//
// The returned value implements `github.com/alexfalkowski/go-service/v2/compress.Compressor`.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements Snappy compression.
//
// It satisfies the `github.com/alexfalkowski/go-service/v2/compress.Compressor` interface.
type Compressor struct{}

// Compress returns the Snappy-compressed representation of data.
func (c *Compressor) Compress(data []byte) []byte {
	return snappy.Encode(nil, data)
}

// Decompress returns the decompressed representation of data.
//
// An error is returned if data is not valid Snappy-encoded content.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
