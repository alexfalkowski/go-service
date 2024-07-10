package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"golang.org/x/net/html"
)

func TestRouteSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		mvc.Register(mux)

		t := mvc.NewView(mvc.ParseTemplate(test.Views, "views/hello.tmpl.html"), mvc.ParseTemplate(test.Views, "views/error.tmpl.html"))
		mvc.Route("GET /hello", t, func(_ context.Context, _ *http.Request, _ http.ResponseWriter) (*test.PageData, error) {
			return &test.Model, nil
		})

		Convey("When I query for hello", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "text/html")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				s := string(body)
				So(s, ShouldNotBeEmpty)

				_, err := html.Parse(strings.NewReader(s))
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestRouteNoController(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		mvc.Register(mux)

		t := mvc.NewView(mvc.ParseTemplate(test.Views, "views/hello.tmpl.html"), mvc.ParseTemplate(test.Views, "views/error.tmpl.html"))
		mvc.Route("GET /hello", t, mvc.NoController)

		Convey("When I query for hello", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "text/html")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				s := string(body)
				So(s, ShouldNotBeEmpty)

				_, err := html.Parse(strings.NewReader(s))
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestRouteError(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		mvc.Register(mux)

		t := mvc.NewView(mvc.ParseTemplate(test.Views, "views/hello.tmpl.html"), mvc.ParseTemplate(test.Views, "views/error.tmpl.html"))
		mvc.Route("GET /hello", t, func(_ context.Context, _ *http.Request, _ http.ResponseWriter) (*test.PageData, error) {
			return nil, status.Error(http.StatusServiceUnavailable, "ohh no")
		})

		Convey("When I query for hello", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "text/html")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(resp.StatusCode, ShouldEqual, 503)

				s := string(body)
				So(s, ShouldNotBeEmpty)

				_, err := html.Parse(strings.NewReader(s))
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
