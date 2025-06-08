package grpc_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/health/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func TestUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
		world.Register()

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())
		server := grpc.NewServer(grpc.ServerParams{Observer: &grpc.Observer{Observer: o}})
		grpc.Register(grpc.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

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

		o := observer(world.Lifecycle, test.StatusURL("500"), world.NewHTTP())
		server := grpc.NewServer(grpc.ServerParams{Observer: &grpc.Observer{Observer: o}})
		grpc.Register(grpc.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

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

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())
		server := grpc.NewServer(grpc.ServerParams{Observer: &grpc.Observer{Observer: o}})
		grpc.Register(grpc.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

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

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())
		server := grpc.NewServer(grpc.ServerParams{Observer: &grpc.Observer{Observer: o}})
		grpc.Register(grpc.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

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

		o := observer(world.Lifecycle, test.StatusURL("500"), world.NewHTTP())
		server := grpc.NewServer(grpc.ServerParams{Observer: &grpc.Observer{Observer: o}})
		grpc.Register(grpc.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

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

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())
		server := grpc.NewServer(grpc.ServerParams{Observer: &grpc.Observer{Observer: o}})
		grpc.Register(grpc.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

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

func observer(lc fx.Lifecycle, url string, client *http.Client) *subscriber.Observer {
	cc := checker.NewHTTPChecker(url, 5*time.Second, checker.WithRoundTripper(client.Transport))
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)
	regs := health.Registrations{hr}
	hs := health.NewServer(lc, regs)

	return hs.Observe("http")
}
