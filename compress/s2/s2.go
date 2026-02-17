package s2

import "github.com/klauspost/compress/s2"

// NewCompressor constructs an s2 compressor.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements s2 compression.
type Compressor struct{}

// Compress compresses data with s2.
func (c *Compressor) Compress(data []byte) []byte {
	return s2.Encode(nil, data)
}

// Decompress decompresses s2-encoded data.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return s2.Decode(nil, data)
}
