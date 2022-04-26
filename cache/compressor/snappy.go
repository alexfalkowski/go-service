package compressor

import (
	"github.com/alexfalkowski/go-service/compressor"
)

// NewSnappy for cache.
// nolint:ireturn
func NewSnappy() Compressor {
	return compressor.NewSnappy()
}
