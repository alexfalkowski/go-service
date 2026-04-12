package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func init() {
	transportgrpc.Register(test.FS)
}

//nolint:funlen
func BenchmarkGRPC(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		b.ReportAllocs()

		l, err := net.Listen(b.Context(), "tcp", "localhost:0")
		require.NoError(b, err)

		server := grpc.NewServer(test.ConfigOptions, test.DefaultTimeout)
		defer server.GracefulStop()

		v1.RegisterGreeterServiceServer(server, test.NewService())

		//nolint:errcheck
		go server.Serve(l)

		conn, err := transportgrpc.NewClientConn(l.Addr().String(), transportgrpc.WithTransportCredentials(transportgrpc.NewInsecureCredentials()))
		require.NoError(b, err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
	})

	b.Run("none", func(b *testing.B) {
		b.ReportAllocs()

		lc := fxtest.NewLifecycle(b)
		cfg := test.NewInsecureTransportConfig()

		g, err := transportgrpc.NewServer(transportgrpc.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC,
			UserAgent:  test.UserAgent, Version: test.Version,
		})
		require.NoError(b, err)
		cfg.GRPC.Address = test.BoundAddress(cfg.GRPC.Address, g.GetService().String())

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		server.Register(lc, []*server.Service{g.GetService()})

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.GRPC.Address)
		cl := &client.Config{Address: addr}

		conn, err := transportgrpc.NewClient(cl.Address)
		require.NoError(b, err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("log", func(b *testing.B) {
		b.ReportAllocs()

		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)
		lc := fxtest.NewLifecycle(b)
		cfg := test.NewInsecureTransportConfig()

		g, err := transportgrpc.NewServer(transportgrpc.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC, Logger: logger,
			UserAgent: test.UserAgent, Version: test.Version,
		})
		require.NoError(b, err)
		cfg.GRPC.Address = test.BoundAddress(cfg.GRPC.Address, g.GetService().String())

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		server.Register(lc, []*server.Service{g.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.GRPC.Address)
		cl := &client.Config{Address: addr}

		conn, err := transportgrpc.NewClient(cl.Address)
		require.NoError(b, err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
		lc.RequireStop()
	})

	b.Run("trace", func(b *testing.B) {
		b.ReportAllocs()

		logger, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)
		lc := fxtest.NewLifecycle(b)
		cfg := test.NewInsecureTransportConfig()

		test.RegisterTracer(lc, nil)

		g, err := transportgrpc.NewServer(transportgrpc.ServerParams{
			Shutdowner: test.NewShutdowner(),
			Config:     cfg.GRPC, Logger: logger,
			UserAgent: test.UserAgent, Version: test.Version,
		})
		require.NoError(b, err)
		cfg.GRPC.Address = test.BoundAddress(cfg.GRPC.Address, g.GetService().String())

		v1.RegisterGreeterServiceServer(g.ServiceRegistrar(), test.NewService())
		server.Register(lc, []*server.Service{g.GetService()})
		errors.Register(errors.NewHandler(logger))

		lc.RequireStart()

		_, addr, _ := net.SplitNetworkAddress(cfg.GRPC.Address)
		cl := &client.Config{Address: addr}

		conn, err := transportgrpc.NewClient(cl.Address)
		require.NoError(b, err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		b.ResetTimer()

		for b.Loop() {
			_, err := client.SayHello(b.Context(), req)
			if err != nil {
				require.NoError(b, err)
			}
		}

		b.StopTimer()
		lc.RequireStop()
	})
}
