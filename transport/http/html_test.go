package http_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	nh "github.com/alexfalkowski/go-service/net/http/html"
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

		nh.Register(mux)

		t := nh.NewView(nh.ParseTemplate(test.HTML, "html/hello.tmpl"), nh.ParseTemplate(test.HTML, "html/error.tmpl"))
		nh.Route("GET /hello", t, func(_ context.Context, _ *http.Request, _ http.ResponseWriter) (*test.PageData, error) {
			return &test.HTMLData, nil
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
				_, err := html.Parse(strings.NewReader(string(body)))
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

		nh.Register(mux)

		t := nh.NewView(nh.ParseTemplate(test.HTML, "html/hello.tmpl"), nh.ParseTemplate(test.HTML, "html/error.tmpl"))
		nh.Route("GET /hello", t, func(_ context.Context, _ *http.Request, _ http.ResponseWriter) (*test.PageData, error) {
			return nil, errors.New("ohh no")
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
				_, err := html.Parse(strings.NewReader(string(body)))
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
