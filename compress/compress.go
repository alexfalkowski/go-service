package compress

import (
	"github.com/alexfalkowski/go-service/compress/none"
	"github.com/alexfalkowski/go-service/compress/s2"
	"github.com/alexfalkowski/go-service/compress/snappy"
	"github.com/alexfalkowski/go-service/compress/zstd"
	"go.uber.org/fx"
)

// MapParams for compressor.
type MapParams struct {
	fx.In

	Zstd   *zstd.Compressor
	S2     *s2.Compressor
	Snappy *snappy.Compressor
	None   *none.Compressor
}

// NewMap for compressor.
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

// Map of compressor.
type Map struct {
	compressors map[string]Compressor
}

// Register kind and compressor.
func (f *Map) Register(kind string, c Compressor) {
	f.compressors[kind] = c
}

// Get from kind.
func (f *Map) Get(kind string) Compressor {
	return f.compressors[kind]
}
