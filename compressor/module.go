package compressor

import (
	"go.uber.org/fx"
)

var (
	// SnappyModule for fx.
	// nolint:gochecknoglobals
	SnappyModule = fx.Options(fx.Provide(NewSnappy))
)
