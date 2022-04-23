package compressor

import (
	"go.uber.org/fx"
)

// SnappyModule for fx.
// nolint:gochecknoglobals
var SnappyModule = fx.Options(fx.Provide(NewSnappy))
