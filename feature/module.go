package feature

import (
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(
		fx.Provide(NewClient),
		fx.Invoke(Register),
	)

	// NoopModule for fx.
	NoopModule = fx.Options(
		Module,
		fx.Provide(NoopProvider),
	)
)
