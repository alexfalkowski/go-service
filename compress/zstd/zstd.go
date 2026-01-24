package zstd

import "github.com/klauspost/compress/zstd"

// NewCompressor for zstd.
func NewCompressor() *Compressor {
	// Error is never returned.
	e, _ := zstd.NewWriter(nil)
	d, _ := zstd.NewReader(nil)

	return &Compressor{encoder: e, decoder: d}
}

// Compressor for zstd.
type Compressor struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

// Compress using zstd.
func (c *Compressor) Compress(data []byte) []byte {
	return c.encoder.EncodeAll(data, nil)
}

// Decompress using zstd.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	return c.decoder.DecodeAll(data, nil)
}
