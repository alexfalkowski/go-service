package test

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/compress"
)

// Compressor contains the real compressor implementations exercised by tests.
var Compressor = compress.NewMap()

// NewCompressor returns a compressor test double whose Compress and Decompress methods fail with err.
func NewCompressor(err error) compress.Compressor {
	return &compressor{err: err}
}

type compressor struct {
	err error
}

// Compress implements [compress.Compressor] for tests.
func (c *compressor) Compress(_ []byte, _ bytes.Size) ([]byte, error) {
	return nil, c.err
}

// Decompress implements [compress.Compressor] for tests.
func (c *compressor) Decompress(_ []byte, _ bytes.Size) ([]byte, error) {
	return nil, c.err
}
