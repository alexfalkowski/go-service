package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	logger.Module,
	metrics.Module,
	tracer.Module,
	errors.Module,
	fx.Invoke(Register),
)
