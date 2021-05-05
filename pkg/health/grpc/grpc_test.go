package grpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/pkg/checker"
	"github.com/alexfalkowski/go-health/pkg/server"
	"github.com/alexfalkowski/go-service/pkg/config"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/grpc"
	"github.com/alexfalkowski/go-service/pkg/health"
	healthGRPC "github.com/alexfalkowski/go-service/pkg/health/grpc"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestGRPC(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/200", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10005"}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		Convey("When I query health", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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

func TestInvalidGRPC(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/500", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10004"}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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

func TestStreamGRPC(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		cc := checker.NewHTTPChecker("https://httpstat.us/200", 1*time.Second)
		hr := server.NewRegistration("http", 10*time.Millisecond, cc)
		regs := health.Registrations{hr}
		lc := fxtest.NewLifecycle(t)

		hs, err := health.NewServer(lc, regs)
		So(err, ShouldBeNil)

		o, err := hs.Observe("http")
		So(err, ShouldBeNil)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10003"}
		gs := pkgGRPC.NewServer(lc, test.NewShutdowner(), cfg, logger)

		healthGRPC.Register(gs, &healthGRPC.Observer{Observer: o})

		lc.RequireStart()

		Convey("When I query health", func() {
			ctx := context.Background()
			opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

			conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
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
