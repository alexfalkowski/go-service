package marshaller

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewProto),
	fx.Provide(NewJSON),
	fx.Provide(NewTOML),
	fx.Provide(NewYAML),
	fx.Provide(NewMap),
)
