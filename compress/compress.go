package compress

import (
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/di"
)

// MapParams defines dependencies used to construct a Map.
//
// It is intended for dependency injection (Fx/Dig). The default wiring is provided by `compress.Module`.
type MapParams struct {
	di.In

	// Zstd is the Zstandard compressor implementation registered under kind "zstd".
	Zstd *zstd.Compressor

	// S2 is the S2 compressor implementation registered under kind "s2".
	S2 *s2.Compressor

	// Snappy is the Snappy compressor implementation registered under kind "snappy".
	Snappy *snappy.Compressor

	// None is the no-op compressor implementation registered under kind "none".
	None *none.Compressor
}

// NewMap constructs a Map pre-populated with the default compressors keyed by kind.
//
// The returned map includes these kinds:
//
//   - "zstd"
//   - "s2"
//   - "snappy"
//   - "none"
//
// Callers can add additional implementations or override existing kinds via (*Map).Register.
func NewMap(params MapParams) *Map {
	return &Map{
		compressors: map[string]Compressor{
			"zstd":   params.Zstd,
			"s2":     params.S2,
			"snappy": params.Snappy,
			"none":   params.None,
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
