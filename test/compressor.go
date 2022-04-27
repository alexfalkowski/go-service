package test

import (
	"github.com/alexfalkowski/go-service/compressor"
)

// NewCompressor for test.
// nolint:ireturn
func NewCompressor(err error) compressor.Compressor {
	return &cmp{err: err}
}

type cmp struct {
	err error
}

func (c *cmp) Compress(src []byte) []byte {
	return nil
}

func (c *cmp) Decompress(src []byte) ([]byte, error) {
	return nil, c.err
}
