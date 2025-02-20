package env

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewUserAgent),
	fx.Provide(NewName),
	fx.Provide(NewVersion),
)
