package config

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Provide(NewConfig)
