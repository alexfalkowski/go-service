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
	"github.com/alexfalkowski/go-service/health"
	shc "github.com/alexfalkowski/go-service/health/checker"
	shh "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alicebob/miniredis/v2"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func init() {
	tracer.Register()
}

func TestHealth(t *testing.T) {
	s := miniredis.RunT(t)
	defer s.Close()

	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			cfg := test.NewInsecureTransportConfig()
			tc := test.NewBaselimeTracerConfig()
			m := metrics.NewNoopMeter()
			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
			client := cl.NewHTTP()
			o := observer(lc, s.Addr(), "http://localhost:6000/v1/status/200", client, logger).Observe("http")

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
			s.Register()

			params := shh.RegisterParams{
				Mux: test.Mux, Health: &shh.HealthObserver{Observer: o},
				Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
				Version: test.Version,
			}
			err := shh.Register(params)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("When I query "+check, func() {
				ctx := context.Background()
				ctx = tm.WithRequestID(ctx, meta.String("test-id"))
				ctx = tm.WithUserAgent(ctx, meta.String("test-user-agent"))

				req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%s/%s", cfg.HTTP.Port, check), http.NoBody)
				So(err, ShouldBeNil)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				actual := strings.TrimSpace(string(body))

				lc.RequireStop()

				Convey("Then I should have a healthy response", func() {
					So(actual, ShouldEqual, "{\n    \"status\": \"SERVING\"\n}")
					So(resp.Header.Get("Version"), ShouldEqual, string(test.Version))
				})
			})
		})
	}
}

func TestReadinessNoop(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		r := miniredis.RunT(t)
		defer r.Close()

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := metrics.NewNoopMeter()
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		server := observer(lc, r.Addr(), "http://localhost:6000/v1/status/500", client, logger)
		o := server.Observe("http")

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		params := shh.RegisterParams{
			Mux: test.Mux, Health: &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: server.Observe("noop")},
			Version: test.Version,
		}
		err := shh.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/readyz", cfg.HTTP.Port), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Add("Request-ID", "test-id")
			req.Header.Add("User-Agent", "test-user-agent")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(actual, ShouldEqual, "{\n    \"status\": \"SERVING\"\n}")
				So(resp.Header.Get("Version"), ShouldEqual, string(test.Version))
			})
		})
	})
}

func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		r := miniredis.RunT(t)
		defer r.Close()

		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := metrics.NewNoopMeter()
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		client := cl.NewHTTP()
		o := observer(lc, r.Addr(), "http://localhost:6000/v1/status/500", client, logger).Observe("http")

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
		s.Register()

		params := shh.RegisterParams{
			Mux: test.Mux, Health: &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
			Version: test.Version,
		}
		err := shh.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/healthz", cfg.HTTP.Port), http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(actual, ShouldEqual, "{\n    \"errors\": {\n        \"http\": \"invalid status code\"\n    },\n    \"status\": \"NOT_SERVING\"\n}")
				So(resp.Header.Get("Version"), ShouldEqual, string(test.Version))
			})
		})
	})
}

func observer(lc fx.Lifecycle, host, url string, client *http.Client, logger *zap.Logger) *server.Server {
	c := &test.Cache{Lifecycle: lc, Redis: test.NewRedisConfig(host, "snappy", "proto"), Logger: logger}
	r := c.NewRedisClient()
	rc := shc.NewRedisChecker(r, 1*time.Second)
	rr := server.NewRegistration("redis", 10*time.Millisecond, rc)

	cc := checker.NewHTTPChecker(url, client)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, rr}

	return health.NewServer(lc, regs)
}
