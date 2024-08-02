package metrics

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewReader),
	fx.Provide(NewMeterProvider),
	fx.Provide(NewMeter),
)
