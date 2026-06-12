package grpc_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/stretchr/testify/require"
)

func TestTokenErrorAuthUnary(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestEmptyAuthUnary(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestMissingClientAuthUnary(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestInvalidAuthUnary(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "x-forwarded-for", "127.0.0.1")
	ctx = meta.AppendToOutgoingContext(ctx, "geolocation", "geo:47,11")

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthUnaryWithAppend(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())

	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "authorization", "What Invalid")

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestAuthStreamWithAppend(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())

	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "authorization", "What Invalid")

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(ctx)
	require.NoError(t, err)

	_, err = test.SendStreamHello(t, stream, "test")
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestAuthUnaryWithLowercaseBearer(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())

	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "authorization", "bearer test")

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	resp, err := client.SayHello(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestValidAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, gen)

			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn), test.WithWorldGRPC())

			conn := test.RequireGRPCConn(t, world)
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			resp, err := client.SayHello(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, "Hello test", resp.GetMessage())
		})
	}
}

func TestUnknownTokenKindAuthUnary(t *testing.T) {
	cfg := test.NewToken("none")
	tkn := token.NewToken(test.Name, cfg, test.FS, nil)

	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("test", nil), tkn), test.WithWorldGRPC())

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.Contains(t, err.Error(), grpc.StatusText(codes.Unauthenticated))
}

func TestBreakerAuthUnary(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
		test.WithWorldCompression(),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world,
		breaker.WithSettings(breaker.Settings{}),
		breaker.WithFailureCodes(codes.Unauthenticated),
	)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	var err error
	for i := range 10 {
		t.Run("attempt-"+strconv.Itoa(i+1), func(t *testing.T) {
			_, err = client.SayHello(t.Context(), req)
		})
	}

	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestValidAuthStream(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	resp, err := test.SendStreamHello(t, stream, "test")
	require.NoError(t, err)

	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestInvalidAuthStream(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	_, err = test.SendStreamHello(t, stream, "test")
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestEmptyAuthStream(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestMissingClientAuthStream(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(nil, test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	_, err = test.SendStreamHello(t, stream, "test")
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestTokenErrorAuthStream(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)

	conn := test.RequireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}
