package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/io"
	servicegrpc "github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/stretchr/testify/require"
)

func requireGRPCConn(tb testing.TB, world *test.World, opts ...breaker.Option) *servicegrpc.ClientConn {
	tb.Helper()

	conn, err := world.NewGRPC(opts...)
	require.NoError(tb, err)

	return conn
}

func sendStreamHello(tb testing.TB, stream v1.GreeterService_SayStreamHelloClient, name string) (*v1.SayStreamHelloResponse, error) {
	tb.Helper()

	if err := stream.Send(&v1.SayStreamHelloRequest{Name: name}); err != nil {
		if errors.Is(err, io.EOF) {
			return stream.Recv()
		}

		return nil, err
	}

	return stream.Recv()
}
