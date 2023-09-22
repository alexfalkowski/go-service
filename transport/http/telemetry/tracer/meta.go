package tracer

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func extract(ctx context.Context, r *http.Request) context.Context {
	prop := otel.GetTextMapPropagator()

	return prop.Extract(ctx, propagation.HeaderCarrier(r.Header))
}

func inject(ctx context.Context, r *http.Request) {
	prop := otel.GetTextMapPropagator()

	prop.Inject(ctx, propagation.HeaderCarrier(r.Header))
}
