package health

import (
	health "github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
)

// RegisterParams defines dependencies for registering HTTP health endpoints.
//
// It is an Fx parameter struct (`di.In`) used to wire standard Kubernetes-style health endpoints:
//
//   - `/healthz` (general health)
//   - `/livez` (liveness)
//   - `/readyz` (readiness)
//
// `Register` is a no-op unless a non-nil `Server` is provided.
type RegisterParams struct {
	di.In

	// Server is the underlying health server that stores and exposes health observers.
	//
	// It is expected to be non-nil when health endpoints are wired.
	Server *health.Server

	// Name is the service name used for route prefixing and observer lookup.
	Name env.Name
}

// Register registers the standard HTTP health endpoints.
//
// It registers REST GET handlers for `/healthz`, `/livez`, and `/readyz` using the route prefix
// `http.Pattern(params.Name, <path>)`.
//
// Each handler checks the corresponding observer on the underlying health server (using the check name
// without the leading slash) and returns:
//
//   - HTTP 200 with `{status: "SERVING"}` when the observer reports no error.
//   - HTTP 503 when the observer is missing or reports an error.
//
// The response also includes request-scoped metadata extracted into the context (see transport metadata middleware).
func Register(params RegisterParams) {
	resister(params.Name, "/healthz", params.Server)
	resister(params.Name, "/livez", params.Server)
	resister(params.Name, "/readyz", params.Server)
}

// Response is the response body returned by the health endpoints.
//
// The shape is designed to be human-readable and machine-consumable.
type Response struct {
	// Meta contains request-scoped metadata derived from the context.
	Meta meta.Map `yaml:"meta,omitempty" json:"meta,omitempty" toml:"meta,omitempty"`

	// Status is the serving status string (for example "SERVING").
	Status string `yaml:"status,omitempty" json:"status,omitempty" toml:"status,omitempty"`
}

func resister(name env.Name, pattern string, server *health.Server) {
	rest.Get(http.Pattern(name, pattern), func(ctx context.Context) (*Response, error) {
		observer, err := server.Observer(name.String(), pattern[1:])
		if err != nil {
			return nil, status.ServiceUnavailableError(err)
		}
		if err := observer.Error(); err != nil {
			return nil, status.ServiceUnavailableError(err)
		}

		return &Response{Status: "SERVING", Meta: meta.CamelStrings(ctx, meta.NoPrefix)}, nil
	})
}
