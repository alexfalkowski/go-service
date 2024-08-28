package compress

import (
	"github.com/alexfalkowski/go-service/compress/none"
	"github.com/alexfalkowski/go-service/compress/s2"
	"github.com/alexfalkowski/go-service/compress/snappy"
	"github.com/alexfalkowski/go-service/compress/zstd"
)

type compressors map[string]Compressor

// Map of compressor.
type Map struct {
	compressors compressors
}

// NewMap for compressor.
func NewMap() *Map {
	f := &Map{
		compressors: compressors{
			"zstd":   zstd.NewCompressor(),
			"s2":     s2.NewCompressor(),
			"snappy": snappy.NewCompressor(),
			"none":   none.NewCompressor(),
		},
	}

	return f
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
