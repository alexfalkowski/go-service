package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Register configures the global OpenTelemetry TextMapPropagator.
//
// It sets the process-wide propagator to a composite propagator containing:
//
//   - W3C Trace Context (propagation.TraceContext)
//   - W3C Baggage (propagation.Baggage)
//
// The global propagator is used by OpenTelemetry instrumentation to extract trace
// context from inbound requests and inject context into outbound requests for
// supported transports (for example HTTP and gRPC) when instrumentation relies
// on otel.GetTextMapPropagator.
//
// Register is intended to be called once during service startup (for example via
// telemetry.Module).
func Register() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
}
