package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// Module wires logger, metrics, tracer, and error helpers into Fx.
var Module = di.Module(
	logger.Module,
	metrics.Module,
	tracer.Module,
	errors.Module,
	di.Register(Register),
)
