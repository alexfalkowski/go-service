package test

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/health"
	hc "github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
	hh "github.com/alexfalkowski/go-service/v2/transport/http/health"
)

// RegisterHealth wires the standard HTTP health routes onto mux for tests.
//
// It delegates to transport/http/health.Register using the shared test service name and the
// provided health server so integration-style tests can exercise the same route registration
// path used in production wiring.
func RegisterHealth(mux *http.ServeMux, server *server.Server) {
	params := hh.RegisterParams{
		Name:   Name,
		Server: server,
		Mux:    mux,
	}

	hh.Register(params)
}

// HealthServer builds a test health server with the default health registrations used by the test world.
//
// The returned server includes:
//   - a database checker backed by the world's configured database,
//   - an HTTP checker targeting url through the world's HTTP transport, and
//   - a noop checker for readiness scenarios.
//
// All registrations are attached under name and share the world's lifecycle so tests can register
// the server and then hit the generated HTTP health endpoints through RegisterHealth.
func (w *World) HealthServer(name, url string) *server.Server {
	db, err := w.OpenDatabase()
	runtime.Must(err)

	dc := hc.NewDBChecker(db, 1*time.Second)
	dr := server.NewRegistration("db", 10*time.Millisecond, dc)

	t := w.NewHTTP().Transport
	cc := checker.NewHTTPChecker(url, 5*time.Second, checker.WithRoundTripper(t))
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, dr}

	server := health.NewServer(w.Lifecycle)
	server.Register(name, regs...)
	return server
}
