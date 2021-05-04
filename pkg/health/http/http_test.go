package http_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/pkg/checker"
	"github.com/alexfalkowski/go-health/pkg/server"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/health"
	healthHTTP "github.com/alexfalkowski/go-service/pkg/health/http"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/http"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

// nolint:dupl
func TestHTTP(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/200", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		server, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := server.Observe("http")
		So(err, ShouldBeNil)

		mux := pkgHTTP.NewMux()
		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{HTTPPort: "10000"}

		pkgHTTP.Register(lc, test.NewShutdowner(), mux, cfg, logger)

		err = healthHTTP.Register(mux, &healthHTTP.Observer{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		time.Sleep(2 * time.Second)

		Convey("When I query health", func() {
			client := &http.Client{Transport: pkgHTTP.NewRoundTripper(logger)}

			req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:10000/health", nil)
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
func TestInvalidHTTP(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/500", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		server, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := server.Observe("http")
		So(err, ShouldBeNil)

		mux := pkgHTTP.NewMux()
		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{HTTPPort: "10001"}

		pkgHTTP.Register(lc, test.NewShutdowner(), mux, cfg, logger)

		err = healthHTTP.Register(mux, &healthHTTP.Observer{Observer: o})
		So(err, ShouldBeNil)

		lc.RequireStart()

		time.Sleep(2 * time.Second)

		Convey("When I query health", func() {
			client := &http.Client{Transport: pkgHTTP.NewRoundTripper(logger)}

			req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:10001/health", nil)
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
