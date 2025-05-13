//nolint:varnamelen
package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/internal/test"
	v1 "github.com/alexfalkowski/go-service/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/net"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/errors"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/transport"
	tg "github.com/alexfalkowski/go-service/transport/grpc"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//nolint:funlen
func BenchmarkGRPC(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		b.ReportAllocs()

		l, err := net.Listen(test.Address())
		runtime.Must(err)

		server := grpc.NewServer()
		defer server.GracefulStop()

		v1.RegisterGreeterServiceServer(server, test.NewService())

		//nolint:errcheck
		go server.Serve(l)

		conn, err := grpc.NewClient(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}

		b.StopTimer()
	})

	b.Run("none", func(b *testing.B) {
		b.ReportAllocs()

		lc := fxtest.NewLifecycle(b)
		cfg := test.NewInsecureTransportConfig()

		g, err := tg.NewServer(tg.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC,
			UserAgent:  test.UserAgent, Version: test.Version,
		})
		runtime.Must(err)

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		transport.Register(lc, []*server.Service{g.GetServer()})

		lc.RequireStart()

		cl := &client.Config{Address: cfg.GRPC.Address}

		conn, err := tg.NewClient(cl.Address)
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("log", func(b *testing.B) {
		b.ReportAllocs()

		logger, _ := logger.NewLogger(logger.Params{})
		lc := fxtest.NewLifecycle(b)
		cfg := test.NewInsecureTransportConfig()

		g, err := tg.NewServer(tg.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC, Logger: logger,
			UserAgent: test.UserAgent, Version: test.Version,
		})
		runtime.Must(err)

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		transport.Register(lc, []*server.Service{g.GetServer()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		cl := &client.Config{Address: cfg.GRPC.Address}

		conn, err := tg.NewClient(cl.Address)
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("trace", func(b *testing.B) {
		b.ReportAllocs()

		logger, _ := logger.NewLogger(logger.Params{})
		lc := fxtest.NewLifecycle(b)
		tracer := test.NewTracer(lc, nil)
		cfg := test.NewInsecureTransportConfig()

		g, err := tg.NewServer(tg.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC, Logger: logger, Tracer: tracer,
			UserAgent: test.UserAgent, Version: test.Version,
		})
		runtime.Must(err)

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		transport.Register(lc, []*server.Service{g.GetServer()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		cl := &client.Config{Address: cfg.GRPC.Address}

		conn, err := tg.NewClient(cl.Address)
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}

		b.StopTimer()
		lc.RequireStop()
	})
}
