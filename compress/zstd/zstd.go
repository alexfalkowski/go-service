package zstd

import (
	"github.com/klauspost/compress/zstd"
)

// Compressor for zstd.
type Compressor struct{}

// NewNone for zstd.
func NewCompressor() *Compressor {
	return &Compressor{}
}

func (c *Compressor) Compress(data []byte) []byte {
	e, _ := zstd.NewWriter(nil)

	return e.EncodeAll(data, nil)
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	d, _ := zstd.NewReader(nil)

	return d.DecodeAll(data, nil)
}
