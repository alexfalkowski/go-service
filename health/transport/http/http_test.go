//nolint:varnamelen
package http_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/health"
	shc "github.com/alexfalkowski/go-service/health/checker"
	shh "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))

			pg.Register(world.Client.NewTracer(), world.Client.Logger)

			so, err := observer(world.Lifecycle, "redis", "http://localhost:6000/v1/status/200", world.NewHTTP(), world.Client.Logger)
			So(err, ShouldBeNil)

			o := so.Observe("http")

			params := shh.RegisterParams{
				Health:   &shh.HealthObserver{Observer: o},
				Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
			}

			shh.Register(params)

			world.Start()

			Convey("When I query "+check, func() {
				ctx := context.Background()
				ctx = tm.WithRequestID(ctx, meta.String("test-id"))
				ctx = tm.WithUserAgent(ctx, meta.String("test-user-agent"))

				header := http.Header{}
				header.Set("Content-Type", "application/json")

				res, body, err := world.Request(ctx, "http", http.MethodGet, check, header, http.NoBody)
				So(err, ShouldBeNil)

				Convey("Then I should have a healthy response", func() {
					So(res.StatusCode, ShouldEqual, http.StatusOK)
					So(body, ShouldContainSubstring, "SERVING")
				})

				world.Stop()
			})
		})
	}
}

func TestReadinessNoop(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))

		pg.Register(world.Client.NewTracer(), world.Client.Logger)
		rest.Register(world.ServeMux, test.Content)

		so, err := observer(world.Lifecycle, "redis", "http://localhost:6000/v1/status/500", world.NewHTTP(), world.Client.Logger)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		params := shh.RegisterParams{
			Health:   &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: so.Observe("noop")},
		}

		shh.Register(params)
		world.Start()

		Convey("When I query health", func() {
			ctx := context.Background()

			header := http.Header{}
			header.Add("Request-Id", "test-id")
			header.Add("User-Agent", "test-user-agent")
			header.Set("Content-Type", "application/json")

			res, body, err := world.Request(ctx, "http", http.MethodGet, "readyz", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(body, ShouldContainSubstring, "SERVING")
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get("Content-Type"), ShouldEqual, "application/json")
			})

			world.Stop()
		})
	})
}

func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))

		pg.Register(world.Client.NewTracer(), world.Client.Logger)
		rest.Register(world.ServeMux, test.Content)

		so, err := observer(world.Lifecycle, "redis", "http://localhost:6000/v1/status/500", world.NewHTTP(), world.Client.Logger)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		params := shh.RegisterParams{
			Health:   &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
		}

		shh.Register(params)
		world.Start()

		Convey("When I query health", func() {
			ctx := context.Background()
			header := http.Header{}

			res, body, err := world.Request(ctx, "http", http.MethodGet, "healthz", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an unhealthy response", func() {
				So(body, ShouldEqual, "rest: http: invalid status code")
				So(res.StatusCode, ShouldEqual, 503)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/plain; charset=utf-8")
			})

			world.Stop()
		})
	})
}

func observer(lc fx.Lifecycle, secret, url string, client *http.Client, logger *zap.Logger) (*server.Server, error) {
	c := &test.Cache{Lifecycle: lc, Redis: test.NewRedisConfig(secret, "snappy", "proto"), Logger: logger}

	r, err := c.NewRedisClient()
	if err != nil {
		return nil, err
	}

	db, err := pg.Open(pg.OpenParams{Lifecycle: lc, Config: test.NewPGConfig()})
	if err != nil {
		return nil, err
	}

	rc := shc.NewRedisChecker(r, 1*time.Second)
	rr := server.NewRegistration("redis", 10*time.Millisecond, rc)

	dc := shc.NewDBChecker(db, 1*time.Second)
	dr := server.NewRegistration("db", 10*time.Millisecond, dc)

	cc := checker.NewHTTPChecker(url, client.Transport, 5*time.Second)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, rr, dr}

	return health.NewServer(lc, regs), nil
}
