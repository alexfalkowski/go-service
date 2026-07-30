package grpc_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/strings"
	telemetryerrors "github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/stretchr/testify/require"
)

func init() {
	transportgrpc.Register(test.FS)
}

// BenchmarkGRPC compares the standard gRPC stack with the supported go-service server stack and telemetry layers.
func BenchmarkGRPC(b *testing.B) {
	b.Run("std", benchmarkStdGRPC)
	benchmarkGRPCLayers(b, benchmarkGRPC)
}

// BenchmarkGRPCStream measures gRPC bidirectional streaming (SayStreamHello) with fixed 1 KiB request payloads.
func BenchmarkGRPCStream(b *testing.B) {
	for _, messageCount := range []int{1, 10, 100} {
		b.Run(fmt.Sprintf("%d-messages", messageCount), func(b *testing.B) {
			b.Helper()

			benchmarkGRPCLayers(b, func(b *testing.B, log *logger.Logger, trace, tlsEnabled bool) {
				b.Helper()

				benchmarkGRPCStream(b, log, trace, tlsEnabled, messageCount)
			})
		})
	}
}

func benchmarkGRPCLayers(b *testing.B, benchmark func(*testing.B, *logger.Logger, bool, bool)) {
	b.Helper()

	b.Run("none", func(b *testing.B) {
		test.ResetTelemetry(b)
		defer test.ResetTelemetry(b)

		benchmark(b, nil, false, false)
	})

	b.Run("log", func(b *testing.B) {
		log, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		benchmark(b, log, false, false)
	})

	b.Run("trace", func(b *testing.B) {
		log, err := logger.NewLogger(logger.LoggerParams{})
		require.NoError(b, err)

		benchmark(b, log, true, false)
	})

	b.Run("tls", func(b *testing.B) {
		benchmark(b, nil, false, true)
	})
}

func benchmarkStdGRPC(b *testing.B) {
	b.ReportAllocs()

	listener, err := net.Listen(b.Context(), "tcp", "localhost:0")
	require.NoError(b, err)

	grpcServer := grpc.NewServer(test.ConfigOptions, test.DefaultTimeout)
	defer grpcServer.GracefulStop()

	v1.RegisterGreeterServiceServer(grpcServer, test.NewService())

	//nolint:errcheck
	go grpcServer.Serve(listener)

	conn, err := transportgrpc.NewClientConn(listener.Addr().String(), transportgrpc.WithTransportCredentials(transportgrpc.NewInsecureCredentials()))
	require.NoError(b, err)

	greeterClient := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	b.ResetTimer()

	for b.Loop() {
		_, err := greeterClient.SayHello(b.Context(), req)
		if err != nil {
			require.NoError(b, err)
		}
	}

	b.StopTimer()
	require.NoError(b, conn.Close())
}

func benchmarkGRPC(b *testing.B, log *logger.Logger, trace, tlsEnabled bool) {
	b.Helper()
	b.ReportAllocs()

	greeterClient, cleanup := startBenchmarkGRPC(b, log, trace, tlsEnabled)
	req := &v1.SayHelloRequest{Name: "test"}

	b.ResetTimer()

	for b.Loop() {
		_, err := greeterClient.SayHello(b.Context(), req)
		if err != nil {
			require.NoError(b, err)
		}
	}

	b.StopTimer()
	cleanup()
}

func benchmarkGRPCStream(b *testing.B, log *logger.Logger, trace, tlsEnabled bool, messageCount int) {
	b.Helper()
	b.ReportAllocs()

	greeterClient, cleanup := startBenchmarkGRPC(b, log, trace, tlsEnabled)
	name := strings.Repeat("b", 1024)

	b.ResetTimer()

	for b.Loop() {
		stream, err := greeterClient.SayStreamHello(b.Context())
		if err != nil {
			require.NoError(b, err)
		}

		for range messageCount {
			if _, err := test.SendStreamHello(b, stream, name); err != nil {
				require.NoError(b, err)
			}
		}

		require.NoError(b, stream.CloseSend())
	}

	b.StopTimer()
	cleanup()
}

func startBenchmarkGRPC(b *testing.B, log *logger.Logger, trace, tlsEnabled bool) (v1.GreeterServiceClient, func()) {
	b.Helper()

	lc := test.QuietLifecycle(b)
	cfg := test.NewInsecureTransportConfig()
	if tlsEnabled {
		cfg = test.NewSecureTransportConfig()
	}

	if trace {
		test.RegisterTracer(lc, nil)
	}

	grpcServer, err := transportgrpc.NewServer(transportgrpc.ServerParams{
		Shutdowner:   test.NewShutdowner(),
		Config:       cfg.GRPC,
		Logger:       log,
		MethodPolicy: transportgrpc.NewMethodPolicy(),
		UserAgent:    test.UserAgent,
		Version:      test.Version,
	})
	require.NoError(b, err)
	cfg.GRPC.Address = test.BoundAddress(cfg.GRPC.Address, grpcServer.GetService().String())

	v1.RegisterGreeterServiceServer(grpcServer.ServiceRegistrar(), &benchmarkGreeterService{Service: test.NewService()})
	if log != nil {
		telemetryerrors.Register(telemetryerrors.NewHandler(nil))
	}

	server.Register(server.RegisterParams{Lifecycle: lc, Drain: server.NewDrain(), Services: []*server.Service{grpcServer.GetService()}})
	lc.RequireStart()

	_, address, _ := net.SplitNetworkAddress(cfg.GRPC.Address)
	clientOpts := []transportgrpc.ClientOption{}
	if tlsEnabled {
		clientOpts = append(clientOpts, transportgrpc.WithClientTLS(test.NewTLSClientConfig()))
	}

	conn, err := transportgrpc.NewClient(address, clientOpts...)
	require.NoError(b, err)

	return v1.NewGreeterServiceClient(conn), func() {
		require.NoError(b, conn.Close())
		lc.RequireStop()
	}
}

type benchmarkGreeterService struct {
	*test.Service
}

func (s *benchmarkGreeterService) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&v1.SayStreamHelloResponse{Message: "Hello " + req.GetName()}); err != nil {
			return err
		}
	}
}
