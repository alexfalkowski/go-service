package http

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/maps"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"go.uber.org/fx"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

// RegisterParams health for HTTP.
type RegisterParams struct {
	fx.In

	Health    *HealthObserver
	Liveness  *LivenessObserver
	Readiness *ReadinessObserver
}

// Register health for HTTP.
func Register(params RegisterParams) {
	resister("/healthz", params.Health.Observer, true)
	resister("/livez", params.Liveness.Observer, false)
	resister("/readyz", params.Readiness.Observer, false)
}

func resister(path string, ob *subscriber.Observer, withErrors bool) {
	rest.Get(path, func(ctx context.Context) (*maps.StringAny, error) {
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

		return &data, nil
	})
}
