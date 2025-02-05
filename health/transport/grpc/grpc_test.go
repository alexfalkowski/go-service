package grpc_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-health/subscriber"
	"github.com/alexfalkowski/go-service/health"
	shg "github.com/alexfalkowski/go-service/health/transport/grpc"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func TestUnary(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
		world.Register()

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())

		shg.Register(shg.RegisterParams{Server: world.GRPCServer, Observer: &shg.Observer{Observer: o}})
		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()
			ctx = tm.WithRequestID(ctx, meta.String("test-id"))
			ctx = tm.WithUserAgent(ctx, meta.String("test-user-agent"))

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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		o := observer(world.Lifecycle, test.StatusURL("500"), world.NewHTTP())

		shg.Register(shg.RegisterParams{Server: world.GRPCServer, Observer: &shg.Observer{Observer: o}})
		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
			ctx = metadata.NewOutgoingContext(ctx, md)

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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")))
		world.Register()

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())

		shg.Register(shg.RegisterParams{Server: world.GRPCServer, Observer: &shg.Observer{Observer: o}})
		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

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

func TestStream(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 10)))
		world.Register()

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())

		shg.Register(shg.RegisterParams{Server: world.GRPCServer, Observer: &shg.Observer{Observer: o}})
		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()

		o := observer(world.Lifecycle, test.StatusURL("500"), world.NewHTTP())

		shg.Register(shg.RegisterParams{Server: world.GRPCServer, Observer: &shg.Observer{Observer: o}})
		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")))
		world.Register()

		o := observer(world.Lifecycle, test.StatusURL("200"), world.NewHTTP())

		shg.Register(shg.RegisterParams{Server: world.GRPCServer, Observer: &shg.Observer{Observer: o}})
		world.RequireStart()

		time.Sleep(1 * time.Second)

		Convey("When I query health", func() {
			ctx := context.Background()

			conn := world.NewGRPC()
			defer conn.Close()

			client := grpc_health_v1.NewHealthClient(conn)
			req := &grpc_health_v1.HealthCheckRequest{}

			wc, err := client.Watch(ctx, req)
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
	cc := checker.NewHTTPChecker(url, client.Transport, 5*time.Second)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)
	regs := health.Registrations{hr}
	hs := health.NewServer(lc, regs)

	return hs.Observe("http")
}
