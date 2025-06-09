package health_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func TestUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.StatusURL("200"))
		o := so.Observe("http")
		server := health.NewServer(health.ServerParams{Observer: &health.Observer{Observer: o}})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := t.Context()
			ctx = meta.WithRequestID(ctx, meta.String("test-id"))
			ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestInvalidUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.StatusURL("500"))
		o := so.Observe("http")
		server := health.NewServer(health.ServerParams{Observer: &health.Observer{Observer: o}})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
			ctx := metadata.NewOutgoingContext(t.Context(), md)

			resp, err := client.Check(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have an unhealthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestIgnoreAuthUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.StatusURL("200"))
		o := so.Observe("http")
		server := health.NewServer(health.ServerParams{Observer: &health.Observer{Observer: o}})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			resp, err := client.Check(t.Context(), req)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.StatusURL("200"))
		o := so.Observe("http")
		server := health.NewServer(health.ServerParams{Observer: &health.Observer{Observer: o}})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(t.Context(), req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestInvalidStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.StatusURL("500"))
		o := so.Observe("http")
		server := health.NewServer(health.ServerParams{Observer: &health.Observer{Observer: o}})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(t.Context(), req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			})

			world.RequireStop()
		})
	})
}

func TestIgnoreAuthStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
		world.Register()

		so := world.HealthServer(test.StatusURL("200"))
		o := so.Observe("http")
		server := health.NewServer(health.ServerParams{Observer: &health.Observer{Observer: o}})
		health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(t.Context(), req)
			So(err, ShouldBeNil)

			resp, err := wc.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(resp.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})

			world.RequireStop()
		})
	})
}
