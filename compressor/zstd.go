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

func (c *Zstd) Compress(src []byte) []byte {
	e, _ := zstd.NewWriter(nil)

	return e.EncodeAll(src, nil)
}

func (c *Zstd) Decompress(src []byte) ([]byte, error) {
	d, _ := zstd.NewReader(nil)

	return d.DecodeAll(src, nil)
}
