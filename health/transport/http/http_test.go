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
	hchecker "github.com/alexfalkowski/go-service/health/checker"
	hhttp "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func init() {
	tracer.Register()
}

func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			cfg := test.NewInsecureTransportConfig()
			m := test.NewOTLPMeter(lc)
			o := observer(lc, "http://localhost:6000/v1/status/200", test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m), logger).Observe("http")
			hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
			gs := test.NewGRPCServer(lc, logger, test.NewBaselimeTracerConfig(), cfg, false, m, nil, nil)

			test.RegisterTransport(lc, gs, hs)

			params := hhttp.RegisterParams{
				Mux: test.Mux, Health: &hhttp.HealthObserver{Observer: o},
				Liveness: &hhttp.LivenessObserver{Observer: o}, Readiness: &hhttp.ReadinessObserver{Observer: o},
				Version: test.Version,
			}
			err := hhttp.Register(params)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("When I query "+check, func() {
				client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

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
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewOTLPMeter(lc)
		server := observer(lc, "http://localhost:6000/v1/status/500", test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m), logger)
		o := server.Observe("http")
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		params := hhttp.RegisterParams{
			Mux: test.Mux, Health: &hhttp.HealthObserver{Observer: o},
			Liveness: &hhttp.LivenessObserver{Observer: o}, Readiness: &hhttp.ReadinessObserver{Observer: server.Observe("noop")},
			Version: test.Version,
		}
		err := hhttp.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

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
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewOTLPMeter(lc)
		o := observer(lc, "http://localhost:6000/v1/status/500", test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m), logger).Observe("http")
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)

		params := hhttp.RegisterParams{
			Mux: test.Mux, Health: &hhttp.HealthObserver{Observer: o},
			Liveness: &hhttp.LivenessObserver{Observer: o}, Readiness: &hhttp.ReadinessObserver{Observer: o},
			Version: test.Version,
		}
		err := hhttp.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

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

func observer(lc fx.Lifecycle, url string, client *http.Client, logger *zap.Logger) *server.Server {
	r := test.NewRedisClient(lc, test.NewRedisConfig("localhost:6379", "snappy", "proto"), logger)
	rc := hchecker.NewRedisChecker(r, 1*time.Second)
	rr := server.NewRegistration("redis", 10*time.Millisecond, rc)

	cc := checker.NewHTTPChecker(url, client)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, rr}

	return health.NewServer(lc, regs)
}
