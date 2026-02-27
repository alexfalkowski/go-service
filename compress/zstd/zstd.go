package zstd

import "github.com/klauspost/compress/zstd"

// NewCompressor constructs a Zstandard (zstd) compressor implementation.
//
// The returned value implements `github.com/alexfalkowski/go-service/v2/compress.Compressor`.
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Compressor implements Zstandard (zstd) compression.
//
// It satisfies the `github.com/alexfalkowski/go-service/v2/compress.Compressor` interface.
type Compressor struct{}

// Compress returns the zstd-compressed representation of data.
//
// This method uses the klauspost/compress zstd encoder.
func (c *Compressor) Compress(data []byte) []byte {
	e, _ := zstd.NewWriter(nil)
	defer e.Close()

	return e.EncodeAll(data, nil)
}

// Decompress returns the decompressed representation of data.
//
// An error is returned if data is not valid zstd-encoded content.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	d, _ := zstd.NewReader(nil)
	defer d.Close()

	return d.DecodeAll(data, nil)
}
