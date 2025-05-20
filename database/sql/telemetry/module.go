package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry/metrics"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	metrics.Module,
)
