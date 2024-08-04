package telemetry

import (
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	logger.Module,
	metrics.Module,
	tracer.Module,
)
