package health

import (
	"context"

	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"go.uber.org/fx"
)

// RegisterParams for health.
type RegisterParams struct {
	fx.In

	Health    *HealthObserver
	Liveness  *LivenessObserver
	Readiness *ReadinessObserver
}

// Register health for HTTP.
func Register(params RegisterParams) {
	resister("/healthz", params.Health.Observer)
	resister("/livez", params.Liveness.Observer)
	resister("/readyz", params.Readiness.Observer)
}

// Response for health.
type Response struct {
	Meta   meta.Map `yaml:"meta,omitempty" json:"meta,omitempty" toml:"meta,omitempty"`
	Status string   `yaml:"status,omitempty" json:"status,omitempty" toml:"status,omitempty"`
}

func resister(path string, ob *subscriber.Observer) {
	rest.Get(path, func(ctx context.Context) (*Response, error) {
		if err := ob.Error(); err != nil {
			return nil, status.ServiceUnavailableError(err)
		}

		res := &Response{
			Status: "SERVING",
			Meta:   meta.CamelStrings(ctx, ""),
		}

		return res, nil
	})
}
