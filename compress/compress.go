package compress

import (
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
)

// NewMap constructs a Map with the default compressors.
//
// The returned map includes these kinds:
//
//   - "zstd"
//   - "s2"
//   - "snappy"
//   - "none"
//
// Callers can add additional implementations or override existing kinds via [Map.Register].
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

// Map is a registry of compressors keyed by kind (for example "zstd" or "snappy").
//
// This type is a thin convenience around a string-keyed map and is commonly used with configuration
// to select a compression algorithm at runtime.
//
// Map is not concurrency-safe. If you mutate it via Register, do so during initialization.
type Map struct {
	compressors map[string]Compressor
}

// Register adds or replaces a compressor for kind.
//
// If kind already exists, the previous compressor is replaced.
func (f *Map) Register(kind string, c Compressor) {
	f.compressors[kind] = c
}

// Get returns the compressor registered for kind.
//
// If no compressor is registered for kind, Get returns nil. Callers typically treat nil as "unknown kind"
// and fall back to a default (for example "none").
func (f *Map) Get(kind string) Compressor {
	return f.compressors[kind]
}
