package env

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewID),
	fx.Provide(NewName),
	fx.Provide(NewVersion),
	fx.Provide(NewUserAgent),
)
