package http

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/maps"
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
	Content   *content.Content
}

// Register health for HTTP.
func Register(params RegisterParams) error {
	mux := params.Mux

	resister("/healthz", mux, params.Health.Observer, params.Content, true)
	resister("/livez", mux, params.Liveness.Observer, params.Content, false)
	resister("/readyz", mux, params.Readiness.Observer, params.Content, false)

	return nil
}

func resister(path string, mux *http.ServeMux, ob *subscriber.Observer, ct *content.Content, withErrors bool) {
	h := ct.NewHandler("health", func(ctx context.Context) (any, error) {
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

		data := maps.StringAny{"status": response}

		if withErrors {
			errors := maps.StringAny{}

			for n, e := range ob.Errors() {
				if e == nil {
					continue
				}

				errors[n] = e.Error()
			}

			if !errors.IsEmpty() {
				data["errors"] = errors
			}
		}

		return data, nil
	})

	mux.HandleFunc("GET "+path, h)
}
