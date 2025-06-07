package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type (
	// Carrier for propagation.
	Carrier = propagation.TextMapCarrier

	// HeaderCarrier for propagation.
	HeaderCarrier = propagation.HeaderCarrier
)

// Register for propagation.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// Extract from the carrier.
func Extract(ctx context.Context, carrier Carrier) context.Context {
	prop := otel.GetTextMapPropagator()

	return prop.Extract(ctx, carrier)
}

// Inject the carrier.
func Inject(ctx context.Context, carrier Carrier) {
	prop := otel.GetTextMapPropagator()

	prop.Inject(ctx, carrier)
}
