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
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func init() {
	otel.Register()
}

func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			cfg := test.NewTransportConfig()
			o := observer(lc, "https://httpstat.us/200", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg), logger).Observe("http")
			hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
			gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

			test.RegisterTransport(lc, cfg, gs, hs)

			params := hhttp.RegisterParams{
				Server: hs, Health: &hhttp.HealthObserver{Observer: o},
				Liveness: &hhttp.LivenessObserver{Observer: o}, Readiness: &hhttp.ReadinessObserver{Observer: o},
				Version: test.Version,
			}
			err := hhttp.Register(params)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey(fmt.Sprintf("When I query %s", check), func() {
				client := test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg)

				req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/%s", cfg.Port, check), nil)
				So(err, ShouldBeNil)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				actual := strings.TrimSpace(string(body))

				lc.RequireStop()

				Convey("Then I should have a healthy response", func() {
					So(actual, ShouldEqual, `{"status":"SERVING"}`)
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
		cfg := test.NewTransportConfig()
		server := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg), logger)
		o := server.Observe("http")
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)

		params := hhttp.RegisterParams{
			Server: hs, Health: &hhttp.HealthObserver{Observer: o},
			Liveness: &hhttp.LivenessObserver{Observer: o}, Readiness: &hhttp.ReadinessObserver{Observer: server.Observe("noop")},
			Version: test.Version,
		}
		err := hhttp.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/readyz", cfg.Port), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(actual, ShouldEqual, `{"status":"SERVING"}`)
				So(resp.Header.Get("Version"), ShouldEqual, string(test.Version))
			})
		})
	})
}

func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		o := observer(lc, "https://httpstat.us/500", test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg), logger).Observe("http")
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)

		params := hhttp.RegisterParams{
			Server: hs, Health: &hhttp.HealthObserver{Observer: o},
			Liveness: &hhttp.LivenessObserver{Observer: o}, Readiness: &hhttp.ReadinessObserver{Observer: o},
			Version: test.Version,
		}
		err := hhttp.Register(params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/healthz", cfg.Port), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(actual, ShouldEqual, `{"errors":{"http":"invalid status code"},"status":"NOT_SERVING"}`)
				So(resp.Header.Get("Version"), ShouldEqual, string(test.Version))
			})
		})
	})
}

func observer(lc fx.Lifecycle, url string, client *http.Client, logger *zap.Logger) *server.Server {
	r := test.NewRedisClient(lc, "localhost:6379", logger)
	rc := hchecker.NewRedisChecker(r, 1*time.Second)
	rr := server.NewRegistration("redis", 10*time.Millisecond, rc)

	cc := checker.NewHTTPChecker(url, client)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, rr}

	return health.NewServer(lc, regs)
}
