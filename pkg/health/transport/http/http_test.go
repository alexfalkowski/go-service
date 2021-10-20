package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/pkg/checker"
	"github.com/alexfalkowski/go-health/pkg/server"
	"github.com/alexfalkowski/go-service/pkg/health"
	healthHTTP "github.com/alexfalkowski/go-service/pkg/health/transport/http"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

// nolint:dupl
func TestHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/200", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		server, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := server.Observe("http")
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		httpServer := pkgHTTP.NewServer(lc, test.NewShutdowner(), cfg, logger)

		err = healthHTTP.Register(httpServer, &healthHTTP.HealthObserver{Observer: o}, &healthHTTP.LivenessObserver{Observer: o}, &healthHTTP.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := pkgHTTP.NewClient(logger)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/health", cfg.Port), nil)
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
			})
		})
	})
}

// nolint:dupl
func TestLiveness(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/200", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		server, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := server.Observe("http")
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		httpServer := pkgHTTP.NewServer(lc, test.NewShutdowner(), cfg, logger)

		err = healthHTTP.Register(httpServer, &healthHTTP.HealthObserver{Observer: o}, &healthHTTP.LivenessObserver{Observer: o}, &healthHTTP.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := pkgHTTP.NewClient(logger)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/liveness", cfg.Port), nil)
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
			})
		})
	})
}

// nolint:dupl
func TestReadiness(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/200", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		server, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := server.Observe("http")
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		httpServer := pkgHTTP.NewServer(lc, test.NewShutdowner(), cfg, logger)

		err = healthHTTP.Register(httpServer, &healthHTTP.HealthObserver{Observer: o}, &healthHTTP.LivenessObserver{Observer: o}, &healthHTTP.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := pkgHTTP.NewClient(logger)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/readiness", cfg.Port), nil)
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
			})
		})
	})
}

// nolint:dupl
func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/500", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		server, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := server.Observe("http")
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		httpServer := pkgHTTP.NewServer(lc, test.NewShutdowner(), cfg, logger)

		err = healthHTTP.Register(httpServer, &healthHTTP.HealthObserver{Observer: o}, &healthHTTP.LivenessObserver{Observer: o}, &healthHTTP.ReadinessObserver{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query health", func() {
			client := pkgHTTP.NewClient(logger)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/health", cfg.Port), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(actual, ShouldEqual, `{"status":"NOT_SERVING"}`)
			})
		})
	})
}
