package zstd

import "github.com/klauspost/compress/zstd"

// NewCompressor for zstd.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor for zstd.
type Compressor struct{}

// Compress using zstd.
func (c *Compressor) Compress(data []byte) []byte {
	e, _ := zstd.NewWriter(nil)
	defer e.Close()

	return e.EncodeAll(data, nil)
}

// Decompress using zstd.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	d, _ := zstd.NewReader(nil)
	defer d.Close()

	return d.DecodeAll(data, nil)
}
