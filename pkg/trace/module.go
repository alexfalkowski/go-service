package trace

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Invoke(Register)
