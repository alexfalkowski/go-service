package meta

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Invoke(RegisterKeys),
)
