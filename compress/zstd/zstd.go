package zstd

import "github.com/klauspost/compress/zstd"

// NewCompressor for zstd.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor for zstd.
type Compressor struct{}

func (c *Compressor) Compress(data []byte) []byte {
	e, _ := zstd.NewWriter(nil)
	return e.EncodeAll(data, nil)
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	d, _ := zstd.NewReader(nil)
	return d.DecodeAll(data, nil)
}
