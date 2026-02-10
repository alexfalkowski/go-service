package test

import (
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
)

// Compressor for tests.
var Compressor = compress.NewMap(compress.MapParams{
	Zstd:   zstd.NewCompressor(),
	S2:     s2.NewCompressor(),
	Snappy: snappy.NewCompressor(),
	None:   none.NewCompressor(),
})

// NewCompressor for test.
func NewCompressor(err error) compress.Compressor {
	return &compressor{err: err}
}

type compressor struct {
	err error
}

// Compress implements compress.Compressor for tests.
func (c *compressor) Compress(_ []byte) []byte {
	return nil
}

// Decompress implements compress.Compressor for tests.
func (c *compressor) Decompress(_ []byte) ([]byte, error) {
	return nil, c.err
}
