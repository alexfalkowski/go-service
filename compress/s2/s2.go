package s2

import "github.com/klauspost/compress/s2"

// NewCompressor constructs an S2 compressor implementation.
//
// The returned value implements `github.com/alexfalkowski/go-service/v2/compress.Compressor`.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements S2 compression.
//
// It satisfies the `github.com/alexfalkowski/go-service/v2/compress.Compressor` interface.
type Compressor struct{}

// Compress returns the S2-compressed representation of data.
func (c *Compressor) Compress(data []byte) []byte {
	return s2.Encode(nil, data)
}

// Decompress returns the decompressed representation of data.
//
// An error is returned if data is not valid S2-encoded content.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return s2.Decode(nil, data)
}
