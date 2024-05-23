package http_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewOTLPMeter(lc)
		tc := test.NewOTLPTracerConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: test.GatewayMux}
		s.Register()

		lc.RequireStart()

		ctx := meta.WithAttribute(context.Background(), "error", meta.Error(http.ErrBodyNotAllowed))

		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
		defer cancel()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		conn := cl.NewGRPC()
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, test.RuntimeMux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")
			req.Header.Set("X-Forwarded-For", "test")
			req.Header.Set("Geolocation", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"message":"Hello test"}`)
			})

			lc.RequireStop()
		})
	})
}

func TestDefaultClientUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: test.GatewayMux}
		s.Register()

		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		conn := cl.NewGRPC()
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, test.RuntimeMux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := http.DefaultClient

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"message":"Hello test"}`)
			})

			lc.RequireStop()
		})
	})
}

func TestSecure(t *testing.T) {
	Convey("Given I a secure client", t, func() {
		mux := sh.NewServeMux(sh.StandardMux, test.RuntimeMux, sh.NewStandardServeMux())
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		m := test.NewPrometheusMeter(lc)
		cfg := test.NewSecureTransportConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			TLS: test.NewTLSClientConfig(),
		}

		lc.RequireStart()

		Convey("When I query github", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(context.Background(), "GET", "https://github.com/alexfalkowski", http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			Convey("Then I should have valid response", func() {
				So(resp.StatusCode, ShouldEqual, 200)
			})
		})

		lc.RequireStop()
	})
}
