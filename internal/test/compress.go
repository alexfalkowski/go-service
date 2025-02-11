package test

import (
	"github.com/alexfalkowski/go-service/compress"
	"github.com/alexfalkowski/go-service/compress/none"
	"github.com/alexfalkowski/go-service/compress/s2"
	"github.com/alexfalkowski/go-service/compress/snappy"
	"github.com/alexfalkowski/go-service/compress/zstd"
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

func (c *compressor) Compress(_ []byte) []byte {
	return nil
}

func (c *compressor) Decompress(_ []byte) ([]byte, error) {
	return nil, c.err
}
