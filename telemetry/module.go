package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// Module composes and wires the telemetry subsystem into Fx.
//
// Including this module in an Fx application enables the go-service telemetry
// wiring for:
//
//   - logging (telemetry/logger.Module),
//   - metrics (telemetry/metrics.Module),
//   - tracing (telemetry/tracer.Module), and
//   - OpenTelemetry internal error handling (telemetry/errors.Module).
//
// In addition, Module registers telemetry.Register, which configures the global
// OpenTelemetry TextMapPropagator (W3C Trace Context + W3C Baggage). This affects
// context extraction/injection performed by instrumentation that relies on the
// global propagator (for example HTTP/gRPC instrumentation).
//
// Note: This module wires providers/exporters and global configuration. It does
// not itself create spans/metrics/log records; instrumentation elsewhere in your
// service emits telemetry using the configured providers.
var Module = di.Module(
	logger.Module,
	metrics.Module,
	tracer.Module,
	errors.Module,
	di.Register(Register),
)
