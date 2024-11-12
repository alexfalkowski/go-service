package test

import (
	"github.com/alexfalkowski/go-service/compress"
)

// Compressor for tests.
var Compressor = compress.NewMap()

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
