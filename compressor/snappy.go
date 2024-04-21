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

func (c *Snappy) Compress(src []byte) []byte {
	return snappy.Encode(nil, src)
}

func (c *Snappy) Decompress(src []byte) ([]byte, error) {
	return snappy.Decode(nil, src)
}
