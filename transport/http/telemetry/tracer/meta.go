package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry"
)

func extract(ctx context.Context, r *http.Request) context.Context {
	return telemetry.Extract(ctx, telemetry.HeaderCarrier(r.Header))
}

func inject(ctx context.Context, r *http.Request) {
	telemetry.Inject(ctx, telemetry.HeaderCarrier(r.Header))
}
