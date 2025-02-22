package errors

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewHandler),
	fx.Invoke(Register),
)
