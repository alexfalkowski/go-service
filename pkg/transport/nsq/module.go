package nsq

import (
	"go.uber.org/fx"
)

var (
	// ConfigModule for fx.
	ConfigModule = fx.Options(fx.Provide(NewConfig))

	// Module for fx.
	Module = fx.Options(ConfigModule)
)
