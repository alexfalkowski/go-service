package http

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"go.uber.org/fx"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

type (
	// RegisterParams for health.
	RegisterParams struct {
		fx.In

		Health    *HealthObserver
		Liveness  *LivenessObserver
		Readiness *ReadinessObserver
	}

	// Response for health.
	Response struct {
		Errors map[string]string `yaml:"errors,omitempty" json:"errors,omitempty" toml:"errors,omitempty"`
		Status string            `yaml:"status,omitempty" json:"status,omitempty" toml:"status,omitempty"`
	}
)

// Register health for HTTP.
func Register(params RegisterParams) {
	resister("/healthz", params.Health.Observer, true)
	resister("/livez", params.Liveness.Observer, false)
	resister("/readyz", params.Readiness.Observer, false)
}

func resister(path string, ob *subscriber.Observer, withErrors bool) {
	rest.Get(path, func(ctx context.Context) (*Response, error) {
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

		resp := &Response{
			Status: response,
		}

		if withErrors {
			resp.Errors = make(map[string]string)

			for n, e := range ob.Errors() {
				if e == nil {
					continue
				}

				resp.Errors[n] = e.Error()
			}
		}

		return resp, nil
	})
}
