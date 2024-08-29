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
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
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

			so, err := observer(lc, "redis", "http://localhost:6000/v1/status/200", client, logger)
			So(err, ShouldBeNil)

			o := so.Observe("http")

			s := &test.Server{
				Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
				Limiter: l, Key: k, Mux: mux,
			}
			s.Register()

			params := shh.RegisterParams{
				Mux: mux, Health: &shh.HealthObserver{Observer: o},
				Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
				Encoder: test.Encoder,
			}
			err = shh.Register(params)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("When I query "+check, func() {
				ctx := context.Background()
				ctx = tm.WithRequestID(ctx, meta.String("test-id"))
				ctx = tm.WithUserAgent(ctx, meta.String("test-user-agent"))

				req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/%s", cfg.HTTP.Address, check), http.NoBody)
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

		so, err := observer(lc, "redis", "http://localhost:6000/v1/status/500", client, logger)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		params := shh.RegisterParams{
			Mux: mux, Health: &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: so.Observe("noop")},
			Encoder: test.Encoder,
		}
		err = shh.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://%s/readyz", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Add("Request-ID", "test-id")
			req.Header.Add("User-Agent", "test-user-agent")
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

		so, err := observer(lc, "redis", "http://localhost:6000/v1/status/500", client, logger)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		params := shh.RegisterParams{
			Mux: mux, Health: &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
			Encoder: test.Encoder,
		}
		err = shh.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://%s/healthz", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(actual, ShouldEqual, "{\"errors\":{\"http\":\"invalid status code\"},\"status\":\"NOT_SERVING\"}")
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

	rc := shc.NewRedisChecker(r, 1*time.Second)
	rr := server.NewRegistration("redis", 10*time.Millisecond, rc)

	cc := checker.NewHTTPChecker(url, client.Transport, 5*time.Second)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, rr}

	return health.NewServer(lc, regs), nil
}
