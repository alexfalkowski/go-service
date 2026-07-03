package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// Module composes and wires the telemetry subsystem into [go.uber.org/fx].
//
// Including this module in an Fx application enables the go-service telemetry
// wiring for:
//
//   - logging ([github.com/alexfalkowski/go-service/v2/telemetry/logger.Module]),
//   - metrics ([github.com/alexfalkowski/go-service/v2/telemetry/metrics.Module]),
//   - tracing ([github.com/alexfalkowski/go-service/v2/telemetry/tracer.Module]), and
//   - OpenTelemetry internal error handling ([github.com/alexfalkowski/go-service/v2/telemetry/errors.Module]).
//
// In addition, Module constructs [NewPropagator] and registers it with
// [RegisterPropagation], which configures the global OpenTelemetry
// TextMapPropagator. By default this uses W3C Trace Context plus W3C Baggage
// for extraction and injection. This affects context extraction/injection
// performed by instrumentation that relies on the global propagator (for example
// HTTP/gRPC instrumentation).
//
// Note: This module wires providers/exporters and global configuration. It does
// not itself create spans/metrics/log records; instrumentation elsewhere in your
// service emits telemetry using the configured providers.
var Module = di.Module(
	logger.Module,
	metrics.Module,
	tracer.Module,
	errors.Module,
	di.Constructor(NewPropagator),
	di.Register(RegisterPropagation),
)
