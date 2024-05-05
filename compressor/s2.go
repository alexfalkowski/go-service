package compressor

import (
	"github.com/klauspost/compress/s2"
)

// S2 for compressor.
type S2 struct{}

// NewS2 compressor.
func NewS2() *S2 {
	return &S2{}
}

func (c *S2) Compress(data []byte) []byte {
	return s2.Encode(nil, data)
}

func (c *S2) Decompress(data []byte) ([]byte, error) {
	return s2.Decode(nil, data)
}
