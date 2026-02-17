package snappy

import "github.com/klauspost/compress/snappy"

// NewCompressor constructs a snappy compressor.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements snappy compression.
type Compressor struct{}

// Compress compresses data with snappy.
func (c *Compressor) Compress(data []byte) []byte {
	return snappy.Encode(nil, data)
}

// Decompress decompresses snappy-encoded data.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
