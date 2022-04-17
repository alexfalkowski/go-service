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

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	jgrpc "github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	jhttp "github.com/alexfalkowski/go-service/transport/http/security/jwt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

// nolint:funlen
func TestUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		hs := shttp.NewServer(shttp.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: httpCfg, Logger: logger, Tracer: tracer})
		gs := tgrpc.NewServer(tgrpc.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: grpcCfg, Logger: logger, Tracer: tracer})

		v1.RegisterGreeterServiceServer(gs, test.NewServer(false))

		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			tgrpc.WithClientLogger(logger), tgrpc.WithClientConfig(grpcCfg), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(logger, tracer)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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

// nolint:dupl,funlen
func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		hparams := shttp.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: httpCfg, Logger: logger, Tracer: tracer}
		httpServer := shttp.NewServer(hparams)

		verifier := test.NewVerifier("test")
		gparams := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: sh,
			Config:     grpcCfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(gparams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			tgrpc.WithClientLogger(logger), tgrpc.WithClientConfig(grpcCfg), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("test", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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

// nolint:dupl,funlen
func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		hparams := shttp.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: httpCfg, Logger: logger, Tracer: tracer}
		httpServer := shttp.NewServer(hparams)

		verifier := test.NewVerifier("test")
		gparams := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: sh,
			Config:     grpcCfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(gparams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			tgrpc.WithClientLogger(logger), tgrpc.WithClientConfig(grpcCfg), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("bob", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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

// nolint:dupl,funlen
func TestMissingAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		hparams := shttp.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: httpCfg, Logger: logger, Tracer: tracer}
		hs := shttp.NewServer(hparams)

		verifier := test.NewVerifier("test")
		gparams := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: sh,
			Config:     grpcCfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(gparams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			tgrpc.WithClientLogger(logger), tgrpc.WithClientConfig(grpcCfg), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(logger, tracer)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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
				So(actual, ShouldContainSubstring, `authorization token is not provided`)
			})

			lc.RequireStop()
		})
	})
}

// nolint:funlen
func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		hparams := shttp.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: httpCfg, Logger: logger, Tracer: tracer}
		hs := shttp.NewServer(hparams)

		verifier := test.NewVerifier("test")
		gparams := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: sh,
			Config:     grpcCfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(gparams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			tgrpc.WithClientLogger(logger), tgrpc.WithClientConfig(grpcCfg), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			_, err = client.Do(req)

			Convey("Then I should have an auth error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "authorization token is not provided")
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl,funlen
func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		hparams := shttp.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: httpCfg, Logger: logger, Tracer: tracer}
		hs := shttp.NewServer(hparams)

		verifier := test.NewVerifier("test")
		gparams := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: sh,
			Config:     grpcCfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(gparams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			tgrpc.WithClientLogger(logger), tgrpc.WithClientConfig(grpcCfg), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(logger, tracer)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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
				So(actual, ShouldContainSubstring, `authorization token is not provided`)
			})

			lc.RequireStop()
		})
	})
}

// nolint:goerr113,funlen
func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		hparams := shttp.ServerParams{Lifecycle: lc, Shutdowner: sh, Config: httpCfg, Logger: logger, Tracer: tracer}
		hs := shttp.NewServer(hparams)

		verifier := test.NewVerifier("test")
		gparams := tgrpc.ServerParams{
			Lifecycle:  lc,
			Shutdowner: sh,
			Config:     grpcCfg,
			Logger:     logger,
			Tracer:     tracer,
			Unary:      []grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			Stream:     []grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(gparams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			tgrpc.WithClientLogger(logger), tgrpc.WithClientConfig(grpcCfg), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet that will generate a token error", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("", errors.New("token error")), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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
