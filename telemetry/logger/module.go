package logger

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewLogger),
)
