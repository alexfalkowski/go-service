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
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	jwtGRPC "github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	jwtHTTP "github.com/alexfalkowski/go-service/transport/http/security/jwt"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		httpServer := shttp.NewServer(lc, test.NewShutdowner(), httpCfg, logger)

		serverParams := tgrpc.ServerParams{Config: grpcCfg, Logger: logger}
		gs := tgrpc.NewServer(lc, sh, serverParams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(false))

		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port), grpcCfg, logger,
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(logger)

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

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		httpServer := shttp.NewServer(lc, test.NewShutdowner(), httpCfg, logger)

		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, sh, serverParams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port), grpcCfg, logger,
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("test", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, transport)

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

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		httpServer := shttp.NewServer(lc, test.NewShutdowner(), httpCfg, logger)

		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, sh, serverParams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port), grpcCfg, logger,
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("bob", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, transport)

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

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		httpServer := shttp.NewServer(lc, test.NewShutdowner(), httpCfg, logger)

		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, sh, serverParams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port), grpcCfg, logger,
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(logger)

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

func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		httpServer := shttp.NewServer(lc, test.NewShutdowner(), httpCfg, logger)

		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, sh, serverParams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port), grpcCfg, logger,
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, transport)

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

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		httpServer := shttp.NewServer(lc, test.NewShutdowner(), httpCfg, logger)

		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, sh, serverParams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port), grpcCfg, logger,
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(logger)

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

// nolint:goerr113
func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		grpcCfg := test.NewGRPCConfig()
		httpCfg := &shttp.Config{Port: test.GenerateRandomPort()}
		httpServer := shttp.NewServer(lc, test.NewShutdowner(), httpCfg, logger)

		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, sh, serverParams)

		v1.RegisterGreeterServiceServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port), grpcCfg, logger,
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, httpServer.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet that will generate a token error", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("", errors.New("token error")), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, transport)

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
