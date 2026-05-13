package test

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/health"
	healthchecker "github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	grpchealth "github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	httphealth "github.com/alexfalkowski/go-service/v2/transport/http/health"
	"github.com/stretchr/testify/require"
)

// RegisterHealth wires the standard HTTP health routes onto mux for tests.
//
// It delegates to transport/http/health.Register using the shared test service name and the
// provided health server so integration-style tests can exercise the same route registration
// path used in production wiring.
func RegisterHealth(mux *http.ServeMux, server *server.Server) {
	params := httphealth.RegisterParams{
		Name:   Name,
		Server: server,
		Mux:    mux,
	}

	httphealth.Register(params)
}

// RegisterHTTPHealth creates and registers the HTTP health routes on the world's mux.
func (w *World) RegisterHTTPHealth(name, url string, observations ...HealthObservation) *server.Server {
	server := w.HealthServer(name, url)
	w.observeHealth(server, name, observations...)
	RegisterHealth(w.Mux, server)
	w.HTTPHealth = server

	return server
}

// RegisterGRPCHealth creates and registers the gRPC health service on the world's server.
func (w *World) RegisterGRPCHealth(name, url string, observations ...HealthObservation) *server.Server {
	server := w.HealthServer(name, url)
	w.observeHealth(server, name, observations...)

	grpcServer := grpchealth.NewServer(grpchealth.ServerParams{Server: server})
	grpchealth.Register(grpchealth.RegisterParams{
		Registrar: w.GRPCServer.ServiceRegistrar(),
		Server:    grpcServer,
	})
	w.GRPCHealth = server

	return server
}

// HealthServer builds a test health server with the default health registrations used by the test world.
//
// The returned server includes:
//   - a database checker backed by the world's database when the world already has one,
//   - an HTTP checker targeting url through the world's HTTP transport, and
//   - a noop checker for readiness scenarios.
//
// All registrations are attached under name and share the world's lifecycle so tests can register
// the server and then hit the generated HTTP health endpoints through RegisterHealth.
func (w *World) HealthServer(name, url string) *server.Server {
	t := w.httpClient.Transport
	cc := checker.NewHTTPChecker(url, (5 * time.Second).Duration(), checker.WithRoundTripper(t))
	hr := server.NewRegistration("http", (10 * time.Millisecond).Duration(), cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", (10 * time.Millisecond).Duration(), no)

	regs := health.Registrations{hr, nr}
	if w.DB != nil {
		dc := healthchecker.NewDBChecker(w.DB, 1*time.Second)
		dr := server.NewRegistration("db", (10 * time.Millisecond).Duration(), dc)
		regs = append(regs, dr)
	}

	server := health.NewServer(w.Lifecycle)
	server.Register(name, regs...)
	return server
}

func (w *World) observeHealth(server *server.Server, name string, observations ...HealthObservation) {
	w.t.Helper()

	for _, observation := range observations {
		if err := server.Observe(name, observation.Kind, observation.Names...); err != nil {
			require.Fail(w.t, "register health observation", "kind=%s: %v", observation.Kind, err)
		}
	}
}
