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

	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	jgrpc "github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	jhttp "github.com/alexfalkowski/go-service/transport/http/security/jwt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func init() {
	otel.Register()
}

func TestUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, false, nil, nil)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := http.DefaultClient

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("test", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewOTELConfig(), cfg, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("bob", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewOTELConfig(), cfg, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewOTELConfig(), cfg, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTELConfig(), cfg)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
		cfg := test.NewTransportConfig()
		hs := test.NewHTTPServer(lc, logger, test.NewOTELConfig(), cfg)
		gs := test.NewGRPCServer(lc, logger, test.NewOTELConfig(), cfg, true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		test.RegisterTransport(lc, cfg, gs, hs)
		lc.RequireStart()

		ctx := context.Background()
		conn := test.NewGRPCClient(ctx, lc, logger, cfg, test.NewOTELConfig(), nil)
		defer conn.Close()

		err := v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet that will generate a token error", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("", errors.New("token error")), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(lc, logger, test.NewOTELConfig(), cfg, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", cfg.Port), bytes.NewBuffer(message))
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
