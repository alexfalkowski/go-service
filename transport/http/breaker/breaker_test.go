package breaker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/stretchr/testify/require"
)

func TestRoundTripperOpensOnTransportError(t *testing.T) {
	transportErr := errors.New("transport unavailable")
	rt := breaker.NewRoundTripper(
		&test.ErrorRoundTripper{Err: transportErr},
		breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts breaker.Counts) bool {
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
	require.ErrorIs(t, err, breaker.ErrOpenState)
}

func TestRoundTripperClosesBodyWhenBreakerIsOpen(t *testing.T) {
	rt := breaker.NewRoundTripper(
		&test.StatusRoundTripper{Status: http.StatusInternalServerError},
		breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 1
			},
		}),
	)
	first, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(first)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", body)
	require.NoError(t, err)

	res, err = rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, breaker.ErrOpenState)
	require.True(t, body.Closed)
}

func TestRoundTripperOpensOnFailureStatus(t *testing.T) {
	rt := breaker.NewRoundTripper(
		&test.StatusRoundTripper{Status: http.StatusInternalServerError},
		breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 1
			},
		}),
		breaker.WithFailureStatuses(http.StatusInternalServerError),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	res, err = rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, breaker.ErrOpenState)
}

func TestRoundTripperIsolatesBreakersByHost(t *testing.T) {
	transportErr := errors.New("transport unavailable")
	calls := make(map[string]int)
	rt := breaker.NewRoundTripper(
		test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			calls[req.URL.Host]++
			if req.URL.Host == "broken.example" {
				return nil, transportErr
			}

			return test.ResponseWithStatus(http.StatusOK), nil
		}),
		breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 1
			},
		}),
	)

	t.Run("opens breaker for failing host", func(t *testing.T) {
		brokenReq, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://broken.example", http.NoBody)
		require.NoError(t, err)
		res, err := rt.RoundTrip(brokenReq)
		require.Nil(t, res)
		require.ErrorIs(t, err, transportErr)
	})

	t.Run("allows different host", func(t *testing.T) {
		healthyReq, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://healthy.example", http.NoBody)
		require.NoError(t, err)
		res, err := rt.RoundTrip(healthyReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, res.Body.Close())
	})

	t.Run("rejects failing host without transport call", func(t *testing.T) {
		openReq, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://broken.example", http.NoBody)
		require.NoError(t, err)
		res, err := rt.RoundTrip(openReq)
		require.Nil(t, res)
		require.ErrorIs(t, err, breaker.ErrOpenState)
		require.Equal(t, 1, calls["broken.example"])
		require.Equal(t, 1, calls["healthy.example"])
	})
}

func TestRoundTripperCountsFailureStatusWithCustomIsSuccessful(t *testing.T) {
	rt := breaker.NewRoundTripper(
		&test.StatusRoundTripper{Status: http.StatusInternalServerError},
		breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 1
			},
			IsSuccessful: func(error) bool {
				return true
			},
		}),
		breaker.WithFailureStatuses(http.StatusInternalServerError),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	res, err = rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, breaker.ErrOpenState)
}

func TestRoundTripperNilFailureStatusFuncUsesDefault(t *testing.T) {
	rt := breaker.NewRoundTripper(
		&test.StatusRoundTripper{Status: http.StatusInternalServerError},
		breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 1
			},
		}),
		breaker.WithFailureStatusFunc(nil),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	res, err = rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, breaker.ErrOpenState)
}
