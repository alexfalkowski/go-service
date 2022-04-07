package http

import (
	"encoding/json"
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/config"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

// Register health for HTTP.
func Register(server *shttp.Server, hob *HealthObserver, lob *LivenessObserver, rob *ReadinessObserver) error {
	resister("/health", server.Mux, hob.Observer, true)
	resister("/liveness", server.Mux, lob.Observer, false)
	resister("/readiness", server.Mux, hob.Observer, false)

	return nil
}

func resister(path string, mux *runtime.ServeMux, ob *subscriber.Observer, withErrors bool) {
	mux.HandlePath("GET", path, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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
		w.Header().Set("Content-Type", "application/json")

		data := config.Map{"status": response}
		if withErrors {
			errors := config.Map{}
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

		json.NewEncoder(w).Encode(data) // nolint:errcheck,errchkjson
	})
}
