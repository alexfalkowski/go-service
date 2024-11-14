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
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
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

		v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
		r := mvc.NewRouter(mux, v)

		r.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
			return mvc.View("hello.tmpl"), &test.Model
		})

		Convey("When I query for hello", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "text/html")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(resp.StatusCode, ShouldEqual, 200)
				So(resp.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")

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

		v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
		r := mvc.NewRouter(mux, v)

		r.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
			return mvc.View("error.tmpl"), status.Error(http.StatusServiceUnavailable, "ohh no")
		})

		Convey("When I query for hello", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "text/html")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(resp.StatusCode, ShouldEqual, 503)
				So(resp.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")

				s := string(body)
				So(s, ShouldNotBeEmpty)

				_, err := html.Parse(strings.NewReader(s))
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestStaticSuccess(t *testing.T) {
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

		v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
		r := mvc.NewRouter(mux, v)

		r.Static("GET /robots.txt", "static/robots.txt")

		Convey("When I query for robots", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/robots.txt", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(resp.StatusCode, ShouldEqual, 200)
				So(resp.Header.Get("Content-Type"), ShouldEqual, "text/plain; charset=utf-8")

				So(string(body), ShouldNotBeEmpty)
			})

			lc.RequireStop()
		})
	})
}

func TestStaticError(t *testing.T) {
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

		v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
		r := mvc.NewRouter(mux, v)

		r.Static("GET /robots.txt", "static/bob.txt")

		Convey("When I query for hello", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/robots.txt", cfg.HTTP.Address), http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			_, err = io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(resp.StatusCode, ShouldEqual, 500)
			})

			lc.RequireStop()
		})
	})
}

func BenchmarkRoute(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)
	logger := zap.NewNop()
	cfg := test.NewInsecureTransportConfig()
	tc := test.NewOTLPTracerConfig()
	m := test.NewOTLPMeter(lc)

	s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
	s.Register()

	cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

	v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
	r := mvc.NewRouter(mux, v)

	r.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
		return mvc.View("hello.tmpl"), &test.Model
	})

	client := cl.NewHTTP()

	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), http.NoBody)
	runtime.Must(err)

	req.Header.Set("Content-Type", "text/html")

	lc.RequireStart()
	b.ResetTimer()

	b.Run("html", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = client.Do(req)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}
