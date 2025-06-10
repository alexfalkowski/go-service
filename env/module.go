package env

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewID),
	fx.Provide(NewUserAgent),
	fx.Provide(NewUserID),
)
