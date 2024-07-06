package http

import (
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/net/http/content"
	"go.uber.org/fx"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

// RegisterParams health for HTTP.
type RegisterParams struct {
	fx.In

	Mux        *http.ServeMux
	Health     *HealthObserver
	Liveness   *LivenessObserver
	Readiness  *ReadinessObserver
	Marshaller *json.Marshaller
	Version    env.Version
}

// Register health for HTTP.
func Register(params RegisterParams) error {
	mux := params.Mux

	resister("/healthz", mux, params.Health.Observer, params.Version, params.Marshaller, true)
	resister("/livez", mux, params.Liveness.Observer, params.Version, params.Marshaller, false)
	resister("/readyz", mux, params.Readiness.Observer, params.Version, params.Marshaller, false)

	return nil
}

func resister(path string, mux *http.ServeMux, ob *subscriber.Observer, version env.Version, mar *json.Marshaller, withErrors bool) {
	mux.HandleFunc("GET "+path, func(resp http.ResponseWriter, _ *http.Request) {
		content.AddJSONHeader(resp.Header())
		resp.Header().Set("Version", string(version))

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

		resp.WriteHeader(status)

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

		b, _ := mar.Marshal(data)

		resp.Write(b)
	})
}
