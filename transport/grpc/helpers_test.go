package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
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
