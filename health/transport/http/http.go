package http

import (
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/marshaller"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/version"
	"go.uber.org/fx"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

// RegisterParams health for HTTP.
type RegisterParams struct {
	fx.In

	Mux       sh.ServeMux
	Health    *HealthObserver
	Liveness  *LivenessObserver
	Readiness *ReadinessObserver
	JSON      *marshaller.JSON
	Version   version.Version
}

// Register health for HTTP.
func Register(params RegisterParams) error {
	mux := params.Mux

	resister("/healthz", mux, params.Health.Observer, params.Version, params.JSON, true)
	resister("/livez", mux, params.Liveness.Observer, params.Version, params.JSON, false)
	resister("/readyz", mux, params.Readiness.Observer, params.Version, params.JSON, false)

	return nil
}

func resister(path string, mux sh.ServeMux, ob *subscriber.Observer, version version.Version, json *marshaller.JSON, withErrors bool) {
	mux.Handle("GET", path, func(resp http.ResponseWriter, _ *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
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

		b, _ := json.Marshal(data)

		resp.Write(b)
	})
}
