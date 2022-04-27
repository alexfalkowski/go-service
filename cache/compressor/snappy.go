package compressor

import (
	"github.com/alexfalkowski/go-service/compressor"
)

// NewSnappy for cache.
func NewSnappy() Compressor {
	return compressor.NewSnappy()
}
