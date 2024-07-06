package snappy

import (
	"github.com/klauspost/compress/snappy"
)

// Compressor for snappy.
type Compressor struct{}

// NewNone for snappy.
func NewCompressor() *Compressor {
	return &Compressor{}
}

func (c *Compressor) Compress(data []byte) []byte {
	return snappy.Encode(nil, data)
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
