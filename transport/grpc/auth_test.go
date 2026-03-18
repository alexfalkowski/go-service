package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	grpcbreaker "github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTokenErrorAuthUnary(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", test.ErrGenerate), test.NewVerifier("test")),
	)

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestEmptyAuthUnary(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")),
	)

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestMissingClientAuthUnary(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil, test.WithWorldToken(nil, test.NewVerifier("test")))

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestInvalidAuthUnary(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
	)

	ctx := t.Context()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-forwarded-for", "127.0.0.1")
	ctx = metadata.AppendToOutgoingContext(ctx, "geolocation", "geo:47,11")

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthUnaryWithAppend(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil, test.WithWorldTelemetry("otlp"))

	ctx := t.Context()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "What Invalid")

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestValidAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			ec := test.NewEd25519()
			signer, _ := ed25519.NewSigner(test.PEM, ec)
			verifier, _ := ed25519.NewVerifier(test.PEM, ec)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

			world := test.NewStartedGRPCWorld(t, nil, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn))

			conn := test.OpenGRPCConn(t, world)

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			resp, err := client.SayHello(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, "Hello test", resp.GetMessage())
		})
	}
}

func TestBreakerAuthUnary(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, func(world *test.World) {
		world.TransportConfig.GRPC.Retry = nil
	},
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
		test.WithWorldCompression(),
	)

	conn := test.OpenGRPCConn(t, world,
		grpcbreaker.WithSettings(breaker.Settings{
			MaxRequests: 1,
			Interval:    0,
			Timeout:     time.Minute,
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 2
			},
		}),
		grpcbreaker.WithFailureCodes(codes.Unauthenticated),
	)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	_, err = client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	_, err = client.SayHello(t.Context(), req)
	require.Error(t, err)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestValidAuthStream(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")),
	)

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	resp, err := stream.Recv()
	require.NoError(t, err)

	require.Equal(t, "Hello test", resp.GetMessage())
}

func TestInvalidAuthStream(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
	)

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestEmptyAuthStream(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")),
	)

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestMissingClientAuthStream(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(nil, test.NewVerifier("test")),
	)

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestTokenErrorAuthStream(t *testing.T) {
	world := test.NewStartedGRPCWorld(t, nil,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
	)

	conn := test.OpenGRPCConn(t, world)

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}
