package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module for fx.
var Module = di.Module(
	metrics.Module,
)
