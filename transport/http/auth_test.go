package http_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/test"
	ht "github.com/alexfalkowski/go-service/transport/http/security/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("test", nil),
		}

		conn := cl.NewGRPC()
		defer conn.Close()

		rpc.Register(mux, test.Marshaller)
		rpc.Handle("/hello", &test.SuccessHandler{})

		Convey("When I query for an authenticated greet", func() {
			client := cl.NewHTTP()

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")
			req.Header.Set("X-Forwarded-For", "127.0.0.1")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(strings.TrimSpace(string(body)), ShouldNotBeBlank)
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("bob", nil),
		}

		conn := cl.NewGRPC()
		defer conn.Close()

		rpc.Register(mux, test.Marshaller)
		rpc.Handle("/hello", &test.SuccessHandler{})

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")
			req.Header.Set("Authorization", "What Invalid")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(strings.TrimSpace(string(body)), ShouldContainSubstring, `verify token: invalid token`)
			})

			lc.RequireStop()
		})
	})
}

//nolint:dupl
func TestMissingAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		conn := cl.NewGRPC()
		defer conn.Close()

		rpc.Register(mux, test.Marshaller)
		rpc.Handle("/hello", &test.SuccessHandler{})

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(strings.TrimSpace(string(body)), ShouldContainSubstring, "invalid token")
			})

			lc.RequireStop()
		})
	})
}

func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			RoundTripper: ht.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport),
		}

		conn := cl.NewGRPC()
		defer conn.Close()

		rpc.Register(mux, test.Marshaller)
		rpc.Handle("/hello", &test.SuccessHandler{})

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			_, err = client.Do(req)

			Convey("Then I should have an auth error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "authorization is invalid")
			})

			lc.RequireStop()
		})
	})
}

//nolint:dupl
func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		conn := cl.NewGRPC()
		defer conn.Close()

		rpc.Register(mux, test.Marshaller)
		rpc.Handle("/hello", &test.SuccessHandler{})

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(strings.TrimSpace(string(body)), ShouldContainSubstring, "invalid token")
			})

			lc.RequireStop()
		})
	})
}

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
			Verifier: verifier, Mux: mux,
		}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Generator: test.NewGenerator("", errors.New("token error")),
		}

		conn := cl.NewGRPC()
		defer conn.Close()

		rpc.Register(mux, test.Marshaller)
		rpc.Handle("/hello", &test.SuccessHandler{})

		Convey("When I query for a greet that will generate a token error", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			_, err = client.Do(req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "token error")
			})

			lc.RequireStop()
		})
	})
}
