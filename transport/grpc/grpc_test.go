package grpc_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport"
	tg "github.com/alexfalkowski/go-service/transport/grpc"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestInsecureUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Compression: true}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(context.Background(), "test", meta.Ignored("test"))
			ctx = meta.WithAttribute(ctx, "ip", meta.ToRedacted(net.ParseIP("192.168.8.0")))

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var header metadata.MD

			resp, err := client.SayHello(ctx, req, grpc.Header(&header))
			So(err, ShouldBeNil)

			h := header.Get("service-version")

			Convey("Then I should have a valid reply", func() {
				So(h, ShouldNotBeEmpty)
				So(h[0], ShouldEqual, "1.0.0")
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestSecureUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewSecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, TLS: test.NewTLSClientConfig()}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(context.Background(), "ip", meta.ToIgnored(net.ParseIP("192.168.8.0")))

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := meta.WithAttribute(context.Background(), "test", meta.Redacted("test"))

			conn := cl.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
			defer cancel()

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			resp, err := stream.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func BenchmarkGRPC(b *testing.B) {
	b.ReportAllocs()

	lc := fxtest.NewLifecycle(b)
	cfg := test.NewInsecureTransportConfig()

	g, err := tg.NewServer(tg.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg.GRPC,
		UserAgent:  test.UserAgent, Version: test.Version,
	})
	runtime.Must(err)

	v1.RegisterGreeterServiceServer(g.Server(), test.NewService(false))
	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{g}})

	cl := &client.Config{Address: cfg.GRPC.Address}

	conn, err := tg.NewClient(cl.Address)
	runtime.Must(err)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	lc.RequireStart()
	b.ResetTimer()

	b.Run("none", func(b *testing.B) {
		for range b.N {
			_, err := client.SayHello(context.Background(), req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}
