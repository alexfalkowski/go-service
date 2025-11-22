package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Register for propagation.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
