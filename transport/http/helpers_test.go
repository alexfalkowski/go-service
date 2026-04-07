package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	servicehttp "github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/stretchr/testify/require"
)

func requireHTTPClient(tb testing.TB, world *test.World, opts ...breaker.Option) *servicehttp.Client {
	tb.Helper()

	client, err := world.NewHTTP(opts...)
	require.NoError(tb, err)

	return client
}
