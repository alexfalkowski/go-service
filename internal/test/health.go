package test

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/health"
	hc "github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
	hh "github.com/alexfalkowski/go-service/v2/transport/http/health"
)

// RegisterHealth for test.
func RegisterHealth(server *server.Server) {
	params := hh.RegisterParams{
		Name:   Name,
		Server: server,
	}

	hh.Register(params)
}

// HealthServer for test.
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
