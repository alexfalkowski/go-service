package health

import (
	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
)

// RegisterParams for health.
type RegisterParams struct {
	di.In

	Health    *HealthObserver
	Liveness  *LivenessObserver
	Readiness *ReadinessObserver
	Name      env.Name
}

// Register health for HTTP.
func Register(params RegisterParams) {
	resister("/healthz", params.Name, params.Health.Observer)
	resister("/livez", params.Name, params.Liveness.Observer)
	resister("/readyz", params.Name, params.Readiness.Observer)
}

// Response for health.
type Response struct {
	Meta   meta.Map `yaml:"meta,omitempty" json:"meta,omitempty" toml:"meta,omitempty"`
	Status string   `yaml:"status,omitempty" json:"status,omitempty" toml:"status,omitempty"`
}

func resister(pattern string, name env.Name, ob *subscriber.Observer) {
	rest.Get(http.Pattern(pattern, name), func(ctx context.Context) (*Response, error) {
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
