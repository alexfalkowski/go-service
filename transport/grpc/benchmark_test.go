//nolint:varnamelen
package grpc_test

import (
	"net"
	"testing"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/internal/test"
	v1 "github.com/alexfalkowski/go-service/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/transport"
	tg "github.com/alexfalkowski/go-service/transport/grpc"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func BenchmarkDefaultGRPC(b *testing.B) {
	b.ReportAllocs()

	addr := test.Address()

	l, err := net.Listen("tcp", addr)
	runtime.Must(err)

	server := grpc.NewServer()
	defer server.GracefulStop()

	v1.RegisterGreeterServiceServer(server, test.NewService(false))

	//nolint:errcheck
	go server.Serve(l)

	b.ResetTimer()

	b.Run("std", func(b *testing.B) {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		for range b.N {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
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
	transport.Register(lc, []transport.Server{g})

	lc.RequireStart()
	b.ResetTimer()

	b.Run("none", func(b *testing.B) {
		cl := &client.Config{Address: cfg.GRPC.Address}

		conn, err := tg.NewClient(cl.Address)
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		for range b.N {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkLogGRPC(b *testing.B) {
	b.ReportAllocs()

	logger := zap.NewNop()
	lc := fxtest.NewLifecycle(b)
	cfg := test.NewInsecureTransportConfig()

	g, err := tg.NewServer(tg.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg.GRPC, Logger: logger,
		UserAgent: test.UserAgent, Version: test.Version,
	})
	runtime.Must(err)

	v1.RegisterGreeterServiceServer(g.Server(), test.NewService(false))
	transport.Register(lc, []transport.Server{g})

	lc.RequireStart()
	b.ResetTimer()

	b.Run("log", func(b *testing.B) {
		cl := &client.Config{Address: cfg.GRPC.Address}

		conn, err := tg.NewClient(cl.Address)
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		for range b.N {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}

func BenchmarkTraceGRPC(b *testing.B) {
	b.ReportAllocs()

	logger := zap.NewNop()
	tc := test.NewOTLPTracerConfig()
	lc := fxtest.NewLifecycle(b)
	tracer := test.NewTracer(lc, tc, logger)
	cfg := test.NewInsecureTransportConfig()

	g, err := tg.NewServer(tg.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg.GRPC, Logger: logger, Tracer: tracer,
		UserAgent: test.UserAgent, Version: test.Version,
	})
	runtime.Must(err)

	v1.RegisterGreeterServiceServer(g.Server(), test.NewService(false))
	transport.Register(lc, []transport.Server{g})

	lc.RequireStart()
	b.ResetTimer()

	b.Run("trace", func(b *testing.B) {
		cl := &client.Config{Address: cfg.GRPC.Address}

		conn, err := tg.NewClient(cl.Address)
		runtime.Must(err)

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		for range b.N {
			_, err := client.SayHello(b.Context(), req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
	lc.RequireStop()
}
