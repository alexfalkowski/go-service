package compress

import (
	"github.com/alexfalkowski/go-service/compress/none"
	"github.com/alexfalkowski/go-service/compress/s2"
	"github.com/alexfalkowski/go-service/compress/snappy"
	"github.com/alexfalkowski/go-service/compress/zstd"
)

// Map of compressor.
type Map struct {
	compressors map[string]Compressor
}

// NewMap for compressor.
func NewMap() *Map {
	return &Map{
		compressors: map[string]Compressor{
			"zstd":   zstd.NewCompressor(),
			"s2":     s2.NewCompressor(),
			"snappy": snappy.NewCompressor(),
			"none":   none.NewCompressor(),
		},
	}
}

// Register kind and compressor.
func (f *Map) Register(kind string, c Compressor) {
	f.compressors[kind] = c
}

// Get from kind.
func (f *Map) Get(kind string) Compressor {
	c, ok := f.compressors[kind]
	if !ok {
		return f.compressors["none"]
	}

	return c
}
