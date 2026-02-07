package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTokenErrorAuthUnary(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}

func TestEmptyAuthUnary(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}

func TestMissingClientAuthUnary(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}

func TestInvalidAuthUnary(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	ctx := t.Context()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-forwarded-for", "127.0.0.1")
	ctx = metadata.AppendToOutgoingContext(ctx, "geolocation", "geo:47,11")

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}

func TestAuthUnaryWithAppend(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()
	world.RequireStart()

	ctx := t.Context()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "What Invalid")

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	world.RequireStop()
}

func TestValidAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		cfg := test.NewToken(kind)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := uuid.NewGenerator()
		tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		conn := world.NewGRPC()
		defer conn.Close()

		client := v1.NewGreeterServiceClient(conn)
		req := &v1.SayHelloRequest{Name: "test"}

		resp, err := client.SayHello(t.Context(), req)
		require.NoError(t, err)
		require.Equal(t, "Hello test", resp.GetMessage())

		world.RequireStop()
	}
}

func TestBreakerAuthUnary(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
		test.WithWorldCompression(),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC(
		breaker.WithSettings(breaker.Settings{}),
		breaker.WithFailureCodes(codes.Unauthenticated),
	)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	var err error
	for range 10 {
		_, err = client.SayHello(t.Context(), req)
	}

	require.Equal(t, codes.Unavailable, status.Code(err))

	world.RequireStop()
}

func TestValidAuthStream(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	resp, err := stream.Recv()
	require.NoError(t, err)

	require.Equal(t, "Hello test", resp.GetMessage())

	world.RequireStop()
}

func TestInvalidAuthStream(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}

func TestEmptyAuthStream(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, nil), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}

func TestMissingClientAuthStream(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(nil, test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}

func TestTokenErrorAuthStream(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(test.NewGenerator(strings.Empty, test.ErrGenerate), test.NewVerifier("test")),
		test.WithWorldGRPC(),
	)
	world.Register()
	world.RequireStart()

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	world.RequireStop()
}
