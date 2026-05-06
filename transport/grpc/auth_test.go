package grpc_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/stretchr/testify/require"
)

func TestTokenErrorAuthUnary(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator("bob", test.ErrGenerate)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
	)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestEmptyAuthUnary(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator(strings.Empty, nil)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
	)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestMissingClientAuthUnary(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t, test.WithBufconnVerifier(test.NewVerifier("test")))

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestInvalidAuthUnary(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator("bob", nil)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
	)

	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "x-forwarded-for", "127.0.0.1")
	ctx = meta.AppendToOutgoingContext(ctx, "geolocation", "geo:47,11")

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthUnaryWithAppend(t *testing.T) {
	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "authorization", "What Invalid")

	conn := test.NewBufconnGRPCConn(t)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(ctx, req)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestAuthStreamWithAppend(t *testing.T) {
	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "authorization", "What Invalid")

	conn := test.NewBufconnGRPCConn(t)

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(ctx)
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestAuthUnaryWithLowercaseBearer(t *testing.T) {
	ctx := t.Context()
	ctx = meta.AppendToOutgoingContext(ctx, "authorization", "bearer test")

	conn := test.NewBufconnGRPCConn(t, test.WithBufconnVerifier(test.NewVerifier("test")))

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
			ec := test.NewEd25519()
			signer, _ := ed25519.NewSigner(test.PEM, ec)
			verifier, _ := ed25519.NewVerifier(test.PEM, ec)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

			conn := test.NewBufconnGRPCConn(t, test.WithBufconnGenerator(tkn), test.WithBufconnVerifier(tkn))

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
	tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator("test", nil)),
		test.WithBufconnVerifier(tkn),
	)

	client := v1.NewGreeterServiceClient(conn)
	req := &v1.SayHelloRequest{Name: "test"}

	_, err := client.SayHello(t.Context(), req)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.Contains(t, err.Error(), "token: invalid config")
}

func TestBreakerAuthUnary(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator("bob", nil)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
		test.WithBufconnCompression(),
		test.WithBufconnBreaker(
			breaker.WithSettings(breaker.Settings{}),
			breaker.WithFailureCodes(codes.Unauthenticated),
		),
	)

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
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator("test", nil)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
	)

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
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator("bob", nil)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
	)

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestEmptyAuthStream(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator(strings.Empty, nil)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
	)

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestMissingClientAuthStream(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t, test.WithBufconnVerifier(test.NewVerifier("test")))

	client := v1.NewGreeterServiceClient(conn)

	stream, err := client.SayStreamHello(t.Context())
	require.NoError(t, err)

	err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestTokenErrorAuthStream(t *testing.T) {
	conn := test.NewBufconnGRPCConn(t,
		test.WithBufconnGenerator(test.NewGenerator(strings.Empty, test.ErrGenerate)),
		test.WithBufconnVerifier(test.NewVerifier("test")),
	)

	client := v1.NewGreeterServiceClient(conn)

	_, err := client.SayStreamHello(t.Context())
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}
