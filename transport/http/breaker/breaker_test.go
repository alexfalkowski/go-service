package breaker_test

import (
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	base "github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/stretchr/testify/require"
)

func TestRoundTripperOpensOnTransportError(t *testing.T) {
	transportErr := errors.New("transport unavailable")
	rt := breaker.NewRoundTripper(
		errorRoundTripper{err: transportErr},
		breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts base.Counts) bool {
				return counts.ConsecutiveFailures >= 1
			},
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, transportErr)

	res, err = rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, base.ErrOpenState)
}

type errorRoundTripper struct {
	err error
}

func (r errorRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, r.err
}
