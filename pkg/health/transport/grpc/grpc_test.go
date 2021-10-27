package grpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/pkg/checker"
	"github.com/alexfalkowski/go-health/pkg/server"
	"github.com/alexfalkowski/go-service/pkg/health"
	healthGRPC "github.com/alexfalkowski/go-service/pkg/health/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc/security/jwt"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// nolint:dupl
func TestUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		serverParams := pkgGRPC.ServerParams{Config: cfg, Logger: logger}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

// nolint:dupl
func TestInvalidUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/500", test.NewHTTPClient(logger))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		serverParams := pkgGRPC.ServerParams{Config: cfg, Logger: logger}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have an unhealthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			})
		})
	})
}

// nolint:funlen
func TestIgnoreAuthUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("test", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

func TestStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		serverParams := pkgGRPC.ServerParams{Config: cfg, Logger: logger}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

// nolint:funlen
func TestIgnoreAuthStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cc := checker.NewHTTPChecker("https://httpstat.us/200", test.NewHTTPClient(logger))
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		cfg := &pkgGRPC.Config{Port: test.GenerateRandomPort()}
		verifier := test.NewVerifier("test")
		serverParams := pkgGRPC.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), serverParams)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			clientParams := &pkgGRPC.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("test", nil))),
			}

			conn, err := pkgGRPC.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have a healthy response", func() {
				So(resp.Status, ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}
