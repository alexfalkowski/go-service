package grpc_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	grpclimiter "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestServerLimiterUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}
	header := meta.Map{}

	_, err := client.SayHello(t.Context(), req, grpc.Header(&header))
	require.NoError(t, err)
	require.NotEmpty(t, header.Get("ratelimit"))

	rejectedHeader := meta.Map{}
	_, err = client.SayHello(t.Context(), req, grpc.Header(&rejectedHeader))
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.NotEmpty(t, rejectedHeader.Get("ratelimit"))
}

func TestServerLimiterStream(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	require.NoError(t, sayStreamHello(t, client))
	err := sayStreamHello(t, client)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestServerLimiterStreamHeader(t *testing.T) {
	limiter, err := grpclimiter.NewServerLimiter(fxtest.NewLifecycle(t), limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 0))
	require.NoError(t, err)
	ctx := meta.WithAttributes(t.Context(), meta.WithUserAgent(meta.String("test-agent")))
	interceptor := grpclimiter.StreamServerInterceptor(limiter)

	allowed := &test.MetaServerStream{Ctx: ctx}
	err = interceptor(nil, allowed, &grpc.StreamServerInfo{FullMethod: "/greet.v1.GreeterService/SayStreamHello"}, func(any, grpc.ServerStream) error {
		return nil
	})
	require.NoError(t, err)
	require.NotEmpty(t, allowed.Header.Get("ratelimit"))

	rejected := &test.MetaServerStream{Ctx: ctx}
	err = interceptor(nil, rejected, &grpc.StreamServerInfo{FullMethod: "/greet.v1.GreeterService/SayStreamHello"}, func(any, grpc.ServerStream) error {
		return nil
	})
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.NotEmpty(t, rejected.Header.Get("ratelimit"))
}

func TestClientLimiterUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, _ = client.SayHello(t.Context(), req)
	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestClientLimiterStream(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 0)), test.WithWorldGRPC())

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	_ = sayStreamHello(t, client)
	err := sayStreamHello(t, client)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestServerLimiterUsesVerifiedUserIDUnary(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-id", "1s", 0)),
		test.WithWorldToken(&test.SequenceGenerator{}, test.AcceptingVerifier{}),
		test.WithWorldGRPC(),
	)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.NoError(t, err)

	_, err = client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestServerLimiterUsesVerifiedUserIDStream(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-id", "1s", 0)),
		test.WithWorldToken(&test.SequenceGenerator{}, test.AcceptingVerifier{}),
		test.WithWorldGRPC(),
	)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	require.NoError(t, sayStreamHello(t, client))
	err := sayStreamHello(t, client)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestLimiterUnlimitedUnary(t *testing.T) {
	cfg := test.NewLimiterConfig("user-agent", "1s", 10)
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldClientLimiter(cfg),
		test.WithWorldServerLimiter(cfg),
		test.WithWorldGRPC(),
	)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.NoError(t, err)
}

func TestLimiterAuthUnary(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("bob")),
		test.WithWorldGRPC(),
	)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	var err error
	for i := range 10 {
		t.Run("attempt-"+strconv.Itoa(i+1), func(t *testing.T) {
			_, err = client.SayHello(t.Context(), req)
		})
	}
	require.NoError(t, err)
}

func TestServerClosedLimiterUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())

	require.NoError(t, world.Server.GRPCLimiter.Close(t.Context()))

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestClientClosedLimiterUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	require.NoError(t, world.Client.GRPCLimiter.Close(t.Context()))

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func sayStreamHello(t *testing.T, client v1.GreeterServiceClient) error {
	t.Helper()

	stream, err := client.SayStreamHello(t.Context())
	if err != nil {
		return err
	}

	if err := stream.Send(&v1.SayStreamHelloRequest{Name: "test"}); err != nil {
		if errors.Is(err, io.EOF) {
			_, err = stream.Recv()
		}

		return err
	}

	_, err = stream.Recv()
	return err
}
