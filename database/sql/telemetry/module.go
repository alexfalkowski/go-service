package telemetry

import (
	"github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	metrics.Module,
)
