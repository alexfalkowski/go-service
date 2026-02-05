package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport"
	tg "github.com/alexfalkowski/go-service/v2/transport/grpc"
	"go.uber.org/fx/fxtest"
)

func init() {
	tg.Register(test.FS)
}

//nolint:funlen
func BenchmarkGRPC(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		b.ReportAllocs()

		n, a, _ := net.SplitNetworkAddress(test.RandomAddress())

		l, err := net.Listen(b.Context(), n, a)
		runtime.Must(err)

		server := grpc.NewServer(test.ConfigOptions, test.DefaultTimeout)
		defer server.GracefulStop()

		v1.RegisterGreeterServiceServer(server, test.NewService())

		//nolint:errcheck
		go server.Serve(l)

		conn, err := grpc.NewClient(l.Addr().String(), grpc.WithTransportCredentials(grpc.NewInsecureCredentials()))
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
		transport.Register(lc, []*server.Service{g.GetService()})

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.GRPC.Address)
		cl := &client.Config{Address: addr}

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

		logger, _ := logger.NewLogger(logger.LoggerParams{})
		lc := fxtest.NewLifecycle(b)
		cfg := test.NewInsecureTransportConfig()

		g, err := tg.NewServer(tg.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC, Logger: logger,
			UserAgent: test.UserAgent, Version: test.Version,
		})
		runtime.Must(err)

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		transport.Register(lc, []*server.Service{g.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.GRPC.Address)
		cl := &client.Config{Address: addr}

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

		logger, _ := logger.NewLogger(logger.LoggerParams{})
		lc := fxtest.NewLifecycle(b)
		_ = test.NewTracer(lc, nil)
		cfg := test.NewInsecureTransportConfig()

		g, err := tg.NewServer(tg.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC, Logger: logger,
			UserAgent: test.UserAgent, Version: test.Version,
		})
		runtime.Must(err)

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		transport.Register(lc, []*server.Service{g.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.GRPC.Address)
		cl := &client.Config{Address: addr}

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
