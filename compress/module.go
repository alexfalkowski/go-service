package compress

import (
	"github.com/alexfalkowski/go-service/v2/compress/none"
	"github.com/alexfalkowski/go-service/v2/compress/s2"
	"github.com/alexfalkowski/go-service/v2/compress/snappy"
	"github.com/alexfalkowski/go-service/v2/compress/zstd"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires all default compressors and the compressor map into Fx.
var Module = di.Module(
	di.Constructor(snappy.NewCompressor),
	di.Constructor(s2.NewCompressor),
	di.Constructor(zstd.NewCompressor),
	di.Constructor(none.NewCompressor),
	di.Constructor(NewMap),
)
