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
	"github.com/alexfalkowski/go-service/v2/strings"
)

// RegisterParams for health.
type RegisterParams struct {
	di.In
	Server *health.Server
	Name   env.Name
}

// Register health for HTTP.
func Register(params RegisterParams) {
	resister(params.Name, "/healthz", params.Server)
	resister(params.Name, "/livez", params.Server)
	resister(params.Name, "/readyz", params.Server)
}

// Response for health.
type Response struct {
	Meta   meta.Map `yaml:"meta,omitempty" json:"meta,omitempty" toml:"meta,omitempty"`
	Status string   `yaml:"status,omitempty" json:"status,omitempty" toml:"status,omitempty"`
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

		return &Response{Status: "SERVING", Meta: meta.CamelStrings(ctx, strings.Empty)}, nil
	})
}
