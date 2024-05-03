package pg

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(Open),
	fx.Invoke(Register),
)
