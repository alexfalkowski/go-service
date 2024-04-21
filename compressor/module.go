package compressor

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewSnappy),
	fx.Provide(NewFactory),
)
