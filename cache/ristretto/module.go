package ristretto

import (
	"github.com/alexfalkowski/go-service/cache/ristretto/telemetry/metrics"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewCache),
	fx.Invoke(metrics.Register),
)
