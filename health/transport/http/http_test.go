package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/health"
	shc "github.com/alexfalkowski/go-service/health/checker"
	shh "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func init() {
	tracer.Register()
	tm.RegisterKeys()
}

//nolint:funlen
func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig("token", "1s", 0))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := noop.Meter{}
			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
			client := cl.NewHTTP()

			pg.Register(cl.NewTracer(), logger)
			rest.Register(mux, test.Content)

			so, err := observer(lc, "redis", "http://localhost:6000/v1/status/200", client, logger)
			So(err, ShouldBeNil)

			o := so.Observe("http")

			s := &test.Server{
				Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
				Limiter: l, Key: k, Mux: mux,
			}
			s.Register()

			params := shh.RegisterParams{
				Health:   &shh.HealthObserver{Observer: o},
				Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
			}

			shh.Register(params)
			lc.RequireStart()

			Convey("When I query "+check, func() {
				ctx := context.Background()
				ctx = tm.WithRequestID(ctx, meta.String("test-id"))
				ctx = tm.WithUserAgent(ctx, meta.String("test-user-agent"))

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/%s", cfg.HTTP.Address, check), http.NoBody)
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				actual := strings.TrimSpace(string(body))

				lc.RequireStop()

				Convey("Then I should have a healthy response", func() {
					So(actual, ShouldEqual, "{\"status\":\"SERVING\"}")
				})
			})
		})
	}
}

func TestReadinessNoop(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := noop.Meter{}
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()

		pg.Register(cl.NewTracer(), logger)
		rest.Register(mux, test.Content)

		so, err := observer(lc, "redis", "http://localhost:6000/v1/status/500", client, logger)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		params := shh.RegisterParams{
			Health:   &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: so.Observe("noop")},
		}

		shh.Register(params)
		lc.RequireStart()

		Convey("When I query health", func() {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("http://%s/readyz", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Add("Request-Id", "test-id")
			req.Header.Add("User-Agent", "test-user-agent")
			req.Header.Set("Content-Type", "application/json")

			res, err := client.Do(req)
			So(err, ShouldBeNil)

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(actual, ShouldEqual, "{\"status\":\"SERVING\"}")
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get("Content-Type"), ShouldEqual, "application/json")
			})
		})
	})
}

func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := noop.Meter{}
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()

		pg.Register(cl.NewTracer(), logger)
		rest.Register(mux, test.Content)

		so, err := observer(lc, "redis", "http://localhost:6000/v1/status/500", client, logger)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		params := shh.RegisterParams{
			Health:   &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
		}

		shh.Register(params)
		lc.RequireStart()

		Convey("When I query health", func() {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("http://%s/healthz", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			res, err := client.Do(req)
			So(err, ShouldBeNil)

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(actual, ShouldEqual, "{\"errors\":{\"http\":\"invalid status code\"},\"status\":\"NOT_SERVING\"}")
				So(res.StatusCode, ShouldEqual, 503)
				So(res.Header.Get("Content-Type"), ShouldEqual, "application/json")
			})
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
