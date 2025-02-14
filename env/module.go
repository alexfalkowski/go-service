package env

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewVersion),
	fx.Provide(NewUserAgent),
	fx.Provide(NewName),
)
