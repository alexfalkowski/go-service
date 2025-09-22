package s2

import "github.com/klauspost/compress/s2"

// NewCompressor for s2.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor for s2.
type Compressor struct{}

func (c *Compressor) Compress(data []byte) []byte {
	return s2.Encode(nil, data)
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return s2.Decode(nil, data)
}
