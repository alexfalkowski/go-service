package test

import (
	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/v2/health"
	hc "github.com/alexfalkowski/go-service/v2/health/checker"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
	hh "github.com/alexfalkowski/go-service/v2/transport/http/health"
)

// RegisterHealth for test.
func RegisterHealth(health, live, ready *subscriber.Observer) {
	params := hh.RegisterParams{
		Health:    &hh.HealthObserver{Observer: health},
		Liveness:  &hh.LivenessObserver{Observer: live},
		Readiness: &hh.ReadinessObserver{Observer: ready},
	}

	hh.Register(params)
}

// HealthServer for test.
func (w *World) HealthServer(url string) *server.Server {
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

	return health.NewServer(w.Lifecycle, regs)
}
