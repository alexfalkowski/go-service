package compressor

import (
	"github.com/klauspost/compress/zstd"
)

// Zstd for compressor.
type Zstd struct{}

// NewZstd compressor.
func NewZstd() *Zstd {
	return &Zstd{}
}

func (c *Zstd) Compress(data []byte) []byte {
	e, _ := zstd.NewWriter(nil)

	return e.EncodeAll(data, nil)
}

func (c *Zstd) Decompress(data []byte) ([]byte, error) {
	d, _ := zstd.NewReader(nil)

	return d.DecodeAll(data, nil)
}
