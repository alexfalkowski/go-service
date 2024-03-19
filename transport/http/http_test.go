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
	"time"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	gt "github.com/alexfalkowski/go-service/transport/grpc/security/token"
	hl "github.com/alexfalkowski/go-service/transport/http/limiter"
	ht "github.com/alexfalkowski/go-service/transport/http/security/token"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"github.com/urfave/negroni/v3"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func init() {
	tracer.Register()
}

func TestUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx := meta.WithAttribute(context.Background(), "error", meta.Error(http.ErrBodyNotAllowed))

		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
		defer cancel()

		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)
		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewDefaultTracerConfig(), cfg, m)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
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

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)
		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := http.DefaultClient

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

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

//nolint:dupl
func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{gt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{gt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := ht.NewRoundTripper(test.NewGenerator("test", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewDefaultTracerConfig(), cfg, transport, m) //nolint:contextcheck

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

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

//nolint:dupl
func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{gt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{gt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := ht.NewRoundTripper(test.NewGenerator("bob", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewDefaultTracerConfig(), cfg, transport, m) //nolint:contextcheck

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a unauthenticated reply", func() {
				So(actual, ShouldContainSubstring, `could not verify token: invalid token`)
			})

			lc.RequireStop()
		})
	})
}

//nolint:dupl
func TestMissingAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{gt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{gt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewDefaultTracerConfig(), cfg, m) //nolint:contextcheck

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a unauthenticated reply", func() {
				So(actual, ShouldContainSubstring, "authorization is invalid")
			})

			lc.RequireStop()
		})
	})
}

func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{gt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{gt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := ht.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewDefaultTracerConfig(), cfg, transport, m) //nolint:contextcheck

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
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
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{gt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{gt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewDefaultTracerConfig(), cfg, m) //nolint:contextcheck

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a unauthenticated reply", func() {
				So(actual, ShouldContainSubstring, "authorization is invalid")
			})

			lc.RequireStop()
		})
	})
}

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		verifier := test.NewVerifier("test")
		cfg := test.NewInsecureTransportConfig()

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, true, m,
			[]grpc.UnaryServerInterceptor{gt.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{gt.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewDefaultTracerConfig(), nil, m)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet that will generate a token error", func() {
			transport := ht.NewRoundTripper(test.NewGenerator("", errors.New("token error")), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewDefaultTracerConfig(), cfg, transport, m) //nolint:contextcheck

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.HTTP.Port), bytes.NewBuffer(message))
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

func TestGet(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, err := limiter.New("100-S")
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		cfg.GRPC.Enabled = false

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, []negroni.Handler{hl.NewHandler(l, tm.UserAgent)})
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		err = hs.Mux.HandlePath("GET", "/hello", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
			w.Write([]byte("hello!"))
		})
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewDefaultTracerConfig(), cfg, m)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a valid greet", func() {
				So(actual, ShouldEqual, "hello!")
			})

			lc.RequireStop()
		})
	})
}

func TestLimiter(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, err := limiter.New("0-S")
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		cfg.GRPC.Enabled = false

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		hs := test.NewHTTPServer(lc, logger, test.NewDefaultTracerConfig(), cfg, m, []negroni.Handler{hl.NewHandler(l, tm.UserAgent)})
		gs := test.NewGRPCServer(lc, logger, test.NewDefaultTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		err = hs.Mux.HandlePath("GET", "/hello", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
			w.Write([]byte("hello!"))
		})
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewDefaultTracerConfig(), cfg, m)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			Convey("Then I should have been rate limited", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusTooManyRequests)
				So(resp.Header.Get("X-Rate-Limit-Limit"), ShouldEqual, "0")
			})

			lc.RequireStop()
		})
	})
}
