package cmd

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewCommand),
	fx.Provide(NewInputConfig),
	fx.Provide(NewOutputConfig),
)
