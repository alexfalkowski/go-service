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
	jwtGRPC "github.com/alexfalkowski/go-service/pkg/transport/grpc/security/jwt"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	jwtHTTP "github.com/alexfalkowski/go-service/pkg/transport/http/security/jwt"
	"github.com/alexfalkowski/go-service/test"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

		grpcCfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		httpCfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		server := pkgHTTP.NewServer(lc, sh, httpCfg, logger)
		mux := server.Handler.(*runtime.ServeMux)
		serverParams := pkgGRPC.ServerParams{Config: grpcCfg, Logger: logger}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(false))

		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		clientParams := &pkgGRPC.ClientParams{
			Host:   fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			Logger: logger,
		}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := pkgHTTP.NewClient(&pkgHTTP.ClientParams{Logger: logger})

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

		grpcCfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		httpCfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		server := pkgHTTP.NewServer(lc, sh, httpCfg, logger)
		mux := server.Handler.(*runtime.ServeMux)
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{
			Host:   fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			Logger: logger,
		}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("test", nil), http.DefaultTransport)
			httpClientParams := &pkgHTTP.ClientParams{
				Logger:       logger,
				RoundTripper: transport,
			}
			client := pkgHTTP.NewClient(httpClientParams)

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

		grpcCfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		httpCfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		server := pkgHTTP.NewServer(lc, sh, httpCfg, logger)
		mux := server.Handler.(*runtime.ServeMux)
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{
			Host:   fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			Logger: logger,
		}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("bob", nil), http.DefaultTransport)
			httpClientParams := &pkgHTTP.ClientParams{
				Logger:       logger,
				RoundTripper: transport,
			}
			client := pkgHTTP.NewClient(httpClientParams)

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

		grpcCfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		httpCfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		server := pkgHTTP.NewServer(lc, sh, httpCfg, logger)
		mux := server.Handler.(*runtime.ServeMux)
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{
			Host:   fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			Logger: logger,
		}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := pkgHTTP.NewClient(&pkgHTTP.ClientParams{Logger: logger})

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

		grpcCfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		httpCfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		server := pkgHTTP.NewServer(lc, sh, httpCfg, logger)
		mux := server.Handler.(*runtime.ServeMux)
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{
			Host:   fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			Logger: logger,
		}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport)
			httpClientParams := &pkgHTTP.ClientParams{
				Logger:       logger,
				RoundTripper: transport,
			}
			client := pkgHTTP.NewClient(httpClientParams)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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

		grpcCfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		httpCfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		server := pkgHTTP.NewServer(lc, sh, httpCfg, logger)
		mux := server.Handler.(*runtime.ServeMux)
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{
			Host:   fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			Logger: logger,
		}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := pkgHTTP.NewClient(&pkgHTTP.ClientParams{Logger: logger})

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
func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		grpcCfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		httpCfg := &pkgHTTP.Config{Port: test.GenerateRandomPort()}
		server := pkgHTTP.NewServer(lc, sh, httpCfg, logger)
		mux := server.Handler.(*runtime.ServeMux)
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: grpcCfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwtGRPC.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwtGRPC.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, sh, serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		ctx := context.Background()
		clientParams := &pkgGRPC.ClientParams{
			Host:   fmt.Sprintf("127.0.0.1:%s", grpcCfg.Port),
			Logger: logger,
		}
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet that will generate a token error", func() {
			transport := jwtHTTP.NewRoundTripper(test.NewGenerator("", errors.New("token error")), http.DefaultTransport)
			httpClientParams := &pkgHTTP.ClientParams{
				Logger:       logger,
				RoundTripper: transport,
			}
			client := pkgHTTP.NewClient(httpClientParams)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", httpCfg.Port), bytes.NewBuffer(message))
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
