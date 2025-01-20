//nolint:varnamelen
package grpc_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"go.uber.org/fx/fxtest"
)

func BenchmarkGRPC(b *testing.B) {
	b.ReportAllocs()

	lc := fxtest.NewLifecycle(b)
	cfg := test.NewInsecureTransportConfig()

	g, err := grpc.NewServer(grpc.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg.GRPC,
		UserAgent:  test.UserAgent, Version: test.Version,
	})
	runtime.Must(err)

	v1.RegisterGreeterServiceServer(g.Server(), test.NewService(false))
	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{g}})

	cl := &client.Config{Address: cfg.GRPC.Address}

	conn, err := grpc.NewClient(cl.Address)
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
