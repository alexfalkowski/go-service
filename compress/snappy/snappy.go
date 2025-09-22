package snappy

import "github.com/klauspost/compress/snappy"

// NewCompressor for snappy.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor for snappy.
type Compressor struct{}

func (c *Compressor) Compress(data []byte) []byte {
	return snappy.Encode(nil, data)
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
