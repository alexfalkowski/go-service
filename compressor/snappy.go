package compressor

import (
	"github.com/klauspost/compress/snappy"
)

// Snappy for compressor.
type Snappy struct{}

// NewSnappy compressor.
func NewSnappy() *Snappy {
	return &Snappy{}
}

func (c *Snappy) Compress(data []byte) []byte {
	return snappy.Encode(nil, data)
}

func (c *Snappy) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
