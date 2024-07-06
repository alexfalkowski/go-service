package compress

import (
	"github.com/alexfalkowski/go-service/compress/none"
	"github.com/alexfalkowski/go-service/compress/s2"
	"github.com/alexfalkowski/go-service/compress/snappy"
	"github.com/alexfalkowski/go-service/compress/zstd"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(snappy.NewCompressor),
	fx.Provide(s2.NewCompressor),
	fx.Provide(zstd.NewCompressor),
	fx.Provide(none.NewCompressor),
	fx.Provide(NewMap),
)
