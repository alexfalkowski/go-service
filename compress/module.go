package compress

import (
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires the default compressor implementations and the compressor registry into Fx.
//
// It provides constructors for:
//   - snappy.NewCompressor (kind "snappy")
//   - s2.NewCompressor (kind "s2")
//   - zstd.NewCompressor (kind "zstd")
//   - none.NewCompressor (kind "none")
//
// And then constructs a `*compress.Map` via `NewMap`, pre-populated with those implementations.
//
// # Extending / overriding
//
// If you want to support additional kinds (or override a default), register them on the `*Map` after
// construction by calling `(*Map).Register` during initialization.
var Module = di.Module(
	di.Constructor(snappy.NewCompressor),
	di.Constructor(s2.NewCompressor),
	di.Constructor(zstd.NewCompressor),
	di.Constructor(none.NewCompressor),
	di.Constructor(NewMap),
)
