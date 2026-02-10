package compress

import (
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/di"
)

// MapParams defines dependencies used to construct a Map.
type MapParams struct {
	di.In
	Zstd   *zstd.Compressor
	S2     *s2.Compressor
	Snappy *snappy.Compressor
	None   *none.Compressor
}

// NewMap constructs a Map pre-populated with the default compressors keyed by kind.
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

// Map holds compressors keyed by kind (for example "zstd" or "snappy").
type Map struct {
	compressors map[string]Compressor
}

// Register adds or replaces a compressor for kind.
func (f *Map) Register(kind string, c Compressor) {
	f.compressors[kind] = c
}

// Get returns the compressor registered for kind.
func (f *Map) Get(kind string) Compressor {
	return f.compressors[kind]
}
