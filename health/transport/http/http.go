package http

import (
	"encoding/json"
	"net/http"

	"github.com/alexfalkowski/go-health/subscriber"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	serving    = "SERVING"
	notServing = "NOT_SERVING"
)

// Register health for HTTP.
func Register(server *shttp.Server, hob *HealthObserver, lob *LivenessObserver, rob *ReadinessObserver) error {
	if err := resister("/health", server.Mux, hob.Observer); err != nil {
		return err
	}

	if err := resister("/liveness", server.Mux, lob.Observer); err != nil {
		return err
	}

	if err := resister("/readiness", server.Mux, hob.Observer); err != nil {
		return err
	}

	return nil
}

func resister(path string, mux *runtime.ServeMux, ob *subscriber.Observer) error {
	return mux.HandlePath("GET", path, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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

		data := map[string]string{"status": response}

		json.NewEncoder(w).Encode(data) // nolint:errcheck
	})
}
