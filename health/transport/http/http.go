package http

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/net/http/content"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"go.uber.org/fx"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

// RegisterParams health for HTTP.
type RegisterParams struct {
	fx.In

	Mux       *http.ServeMux
	Health    *HealthObserver
	Liveness  *LivenessObserver
	Readiness *ReadinessObserver
	Encoder   *encoding.Map
}

// Register health for HTTP.
func Register(params RegisterParams) error {
	mux := params.Mux

	resister("/healthz", mux, params.Health.Observer, params.Encoder, true)
	resister("/livez", mux, params.Liveness.Observer, params.Encoder, false)
	resister("/readyz", mux, params.Readiness.Observer, params.Encoder, false)

	return nil
}

func resister(path string, mux *http.ServeMux, ob *subscriber.Observer, enc *encoding.Map, withErrors bool) {
	h := content.NewHandler("health", enc, func(ctx context.Context) any {
		var (
			status   int
			response string
		)

		if err := ob.Error(); err != nil {
			status = http.StatusServiceUnavailable
			response = notServing
		} else {
			status = http.StatusOK
			response = serving
		}

		res := hc.Response(ctx)
		res.WriteHeader(status)

		data := map[string]any{"status": response}

		if withErrors {
			errors := map[string]any{}

			for n, e := range ob.Errors() {
				if e == nil {
					continue
				}

				errors[n] = e.Error()
			}

			if len(errors) > 0 {
				data["errors"] = errors
			}
		}

		return data
	})

	mux.HandleFunc("GET "+path, h)
}
