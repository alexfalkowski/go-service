//nolint:varnamelen
package http_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	ht "github.com/alexfalkowski/go-service/transport/http/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestTokenAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "key"} {
		Convey("Given I have a all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			kid, _ := token.NewKID(rand.NewGenerator(rand.NewReader()))
			a, _ := ed25519.NewSigner(test.NewEd25519())
			jwt := token.NewJWT(kid, a)
			pas := token.NewPaseto(a)
			token := token.NewToken(test.NewToken(kind), jwt, pas)
			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{
				Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, VerifyAuth: true,
				Verifier: token, Mux: mux,
			}
			s.Register()

			lc.RequireStart()

			ctx := context.Background()
			cl := &test.Client{
				Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
				Generator: token,
			}

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I query for an authenticated greet", func() {
				client := cl.NewHTTP()

				message := []byte(`{"name":"test"}`)
				req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Request-Id", "test")
				req.Header.Set("X-Forwarded-For", "127.0.0.1")
				req.Header.Set("Geolocation", "geo:47,11")

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				Convey("Then I should have a valid reply", func() {
					So(resp.StatusCode, ShouldEqual, 200)
					So(strings.TrimSpace(string(body)), ShouldNotBeBlank)
				})

				lc.RequireStop()
			})
		})
	}
}

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

		rpc.Register(mux, test.Content, test.Pool)
		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for an authenticated greet", func() {
			client := cl.NewHTTP()

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-Id", "test")
			req.Header.Set("X-Forwarded-For", "127.0.0.1")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.StatusCode, ShouldEqual, 200)
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

		rpc.Register(mux, test.Content, test.Pool)
		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-Id", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(resp.StatusCode, ShouldEqual, 401)
				So(strings.TrimSpace(string(body)), ShouldContainSubstring, `token: invalid match`)
			})

			lc.RequireStop()
		})
	})
}

func TestAuthUnaryWithAppend(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc,
			Transport: cfg, Meter: m, Mux: mux,
		}
		s.Register()

		lc.RequireStart()

		ctx := context.Background()
		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc,
			Transport: cfg, Meter: m,
		}

		rpc.Register(mux, test.Content, test.Pool)
		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-Id", "test")
			req.Header.Set("Authorization", "What Invalid")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			Convey("Then I should have a reply", func() {
				So(resp.StatusCode, ShouldEqual, 200)
			})

			lc.RequireStop()
		})
	})
}

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

		rpc.Register(mux, test.Content, test.Pool)
		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-Id", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(resp.StatusCode, ShouldEqual, 401)
				So(strings.TrimSpace(string(body)), ShouldContainSubstring, "invalid match")
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

		rpc.Register(mux, test.Content, test.Pool)
		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-Id", "test")

			_, err = client.Do(req)

			Convey("Then I should have an auth error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "authorization is invalid")
			})

			lc.RequireStop()
		})
	})
}

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

		rpc.Register(mux, test.Content, test.Pool)
		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-Id", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(strings.TrimSpace(string(body)), ShouldContainSubstring, "invalid match")
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
			Generator: test.NewGenerator("", test.ErrGenerate),
		}

		rpc.Register(mux, test.Content, test.Pool)
		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a greet that will generate a token error", func() {
			client := cl.NewHTTP()
			message := []byte(`{"name":"test"}`)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-Id", "test")

			_, err = client.Do(req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "token error")
			})

			lc.RequireStop()
		})
	})
}
