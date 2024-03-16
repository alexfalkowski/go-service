package marshaller

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewJSON),
	fx.Provide(NewProto),
	fx.Provide(NewTOML),
	fx.Provide(NewYAML),
	fx.Provide(NewFactory),
)
