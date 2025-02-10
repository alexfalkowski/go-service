package s2

import (
	"github.com/klauspost/compress/s2"
)

// Compressor for s2.
type Compressor struct{}

// NewNone for s2.
func NewCompressor() Compressor {
	return Compressor{}
}

func (c Compressor) Compress(data []byte) []byte {
	return s2.Encode(nil, data)
}

func (c Compressor) Decompress(data []byte) ([]byte, error) {
	return s2.Decode(nil, data)
}
