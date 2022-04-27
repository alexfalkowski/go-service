package compressor

import (
	"github.com/klauspost/compress/snappy"
)

type snappyCompressor struct{}

// NewSnappy compressor.
func NewSnappy() Compressor {
	return &snappyCompressor{}
}

func (c *snappyCompressor) Compress(src []byte) []byte {
	return snappy.Encode(nil, src)
}

func (c *snappyCompressor) Decompress(src []byte) ([]byte, error) {
	return snappy.Decode(nil, src)
}
