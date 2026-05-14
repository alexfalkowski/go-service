package test

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
)

// Compressor contains the real compressor implementations exercised by tests.
var Compressor = compress.NewMap(compress.MapParams{
	Zstd:   zstd.NewCompressor(),
	S2:     s2.NewCompressor(),
	Snappy: snappy.NewCompressor(),
	None:   none.NewCompressor(),
})

// NewCompressor returns a compressor test double whose Decompress method fails with the supplied error.
func NewCompressor(err error) compress.Compressor {
	return &compressor{err: err}
}

type compressor struct {
	err error
}

// Compress implements compress.Compressor for tests.
func (c *compressor) Compress(_ []byte, _ bytes.Size) ([]byte, error) {
	return nil, c.err
}

// Decompress implements compress.Compressor for tests.
func (c *compressor) Decompress(_ []byte, _ bytes.Size) ([]byte, error) {
	return nil, c.err
}
