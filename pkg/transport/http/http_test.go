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

	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/transport/grpc"
	tokenGRPC "github.com/alexfalkowski/go-service/pkg/transport/grpc/security/token"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	tokenHTTP "github.com/alexfalkowski/go-service/pkg/transport/http/security/token"
	"github.com/alexfalkowski/go-service/test"
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

		cfg := &pkgGRPC.Config{GRPCPort: "10009"}
		mux := pkgHTTP.NewMux()

		pkgHTTP.Register(lc, sh, mux, &pkgHTTP.Config{HTTPPort: "10010"}, logger)

		serverParams := pkgGRPC.ServerParams{Config: cfg, Logger: logger}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		clientParams := &pkgGRPC.ClientParams{Logger: logger}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := pkgHTTP.NewClient(logger)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10010/v1/greet/hello", bytes.NewBuffer(message))
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

		cfg := &pkgGRPC.Config{GRPCPort: "10011"}
		mux := pkgHTTP.NewMux()

		pkgHTTP.Register(lc, sh, mux, &pkgHTTP.Config{HTTPPort: "10012"}, logger)

		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{Logger: logger}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := tokenHTTP.NewRoundTripper(test.NewGenerator("test", nil), http.DefaultTransport)
			client := pkgHTTP.NewClientWithRoundTripper(logger, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10012/v1/greet/hello", bytes.NewBuffer(message))
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

		cfg := &pkgGRPC.Config{GRPCPort: "10013"}
		mux := pkgHTTP.NewMux()

		pkgHTTP.Register(lc, sh, mux, &pkgHTTP.Config{HTTPPort: "10014"}, logger)

		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{Logger: logger}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := tokenHTTP.NewRoundTripper(test.NewGenerator("bob", nil), http.DefaultTransport)
			client := pkgHTTP.NewClientWithRoundTripper(logger, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10014/v1/greet/hello", bytes.NewBuffer(message))
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

		cfg := &pkgGRPC.Config{GRPCPort: "10013"}
		mux := pkgHTTP.NewMux()

		pkgHTTP.Register(lc, sh, mux, &pkgHTTP.Config{HTTPPort: "10014"}, logger)

		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{Logger: logger}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := pkgHTTP.NewClient(logger)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10014/v1/greet/hello", bytes.NewBuffer(message))
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

		cfg := &pkgGRPC.Config{GRPCPort: "10013"}
		mux := pkgHTTP.NewMux()

		pkgHTTP.Register(lc, sh, mux, &pkgHTTP.Config{HTTPPort: "10014"}, logger)

		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{Logger: logger}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := tokenHTTP.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport)
			client := pkgHTTP.NewClientWithRoundTripper(logger, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10014/v1/greet/hello", bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			_, err = client.Do(req) // nolint:bodyclose

			Convey("Then I should have an auth error", func() {
				So(err, ShouldBeError)
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

		cfg := &pkgGRPC.Config{GRPCPort: "10013"}
		mux := pkgHTTP.NewMux()

		pkgHTTP.Register(lc, sh, mux, &pkgHTTP.Config{HTTPPort: "10014"}, logger)

		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{Logger: logger}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := pkgHTTP.NewClient(logger)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10014/v1/greet/hello", bytes.NewBuffer(message))
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

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{GRPCPort: "10013"}
		mux := pkgHTTP.NewMux()

		pkgHTTP.Register(lc, sh, mux, &pkgHTTP.Config{HTTPPort: "10014"}, logger)

		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{tokenGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{tokenGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{Logger: logger}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet that will generate a token error", func() {
			transport := tokenHTTP.NewRoundTripper(test.NewGenerator("", errors.New("token error")), http.DefaultTransport)
			client := pkgHTTP.NewClientWithRoundTripper(logger, transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10014/v1/greet/hello", bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			Convey("Then I should have an error", func() {
				_, err := client.Do(req) // nolint:bodyclose
				So(err, ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
