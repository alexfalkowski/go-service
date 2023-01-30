package http

import (
	"encoding/json"
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/version"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

// RegisterParams health for HTTP.
type RegisterParams struct {
	fx.In

	Server    *shttp.Server
	Health    *HealthObserver
	Liveness  *LivenessObserver
	Readiness *ReadinessObserver
	Version   version.Version
}

// Register health for HTTP.
func Register(params RegisterParams) error {
	resister("/healthz", params.Server.Mux, params.Health.Observer, params.Version, true)
	resister("/livez", params.Server.Mux, params.Liveness.Observer, params.Version, false)
	resister("/readyz", params.Server.Mux, params.Readiness.Observer, params.Version, false)

	return nil
}

func resister(path string, mux *runtime.ServeMux, ob *subscriber.Observer, version version.Version, withErrors bool) {
	mux.HandlePath("GET", path, func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Version", string(version))

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

		w.WriteHeader(status)

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

		_ = json.NewEncoder(w).Encode(data)
	})
}
