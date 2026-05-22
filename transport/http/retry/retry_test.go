package retry_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/retry"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
)

func TestRoundTripperRetriesRetryableResponses(t *testing.T) {
	tests := []struct {
		name  string
		codes []int
		calls int
	}{
		{name: "too many requests", codes: []int{http.StatusTooManyRequests, http.StatusOK}, calls: 2},
		{name: "service unavailable", codes: []int{http.StatusServiceUnavailable, http.StatusOK}, calls: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &test.StatusSequenceRoundTripper{Codes: tt.codes}
			retrying := retry.NewRoundTripper(&retry.Config{
				Attempts: 2,
				Timeout:  time.Second,
				Backoff:  time.Millisecond,
			}, rt)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

			res, err := retrying.RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, tt.calls, rt.Calls)
		})
	}
}

func TestRoundTripperDoesNotRetryWhenAttemptsIsOne(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 1,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperUsesIndependentRetryBudgetPerRequest(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{
		http.StatusTooManyRequests, http.StatusOK,
		http.StatusTooManyRequests, http.StatusOK,
	}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 4, rt.Calls)
}

func TestRoundTripperDoesNotRetryUnhandledStatusCode(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusInternalServerError, http.StatusOK}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperDoesNotRetryWhenPolicyDeniesRequest(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt, retry.SafeMethods)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperRetriesWhenPolicyAllowsRequestID(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt, retry.IdempotentRequests)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperReturnsLastRetryableResponseWhenExhausted(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusTooManyRequests}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperReturnsFinalRetryableResponseWhenExhausted(t *testing.T) {
	rt := &test.BodySequenceRoundTripper{Responses: []string{"first failure", "second failure"}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	body, _, readErr := io.ReadAll(res.Body)
	require.NoError(t, readErr)
	require.Equal(t, "second failure", string(body))
	require.NoError(t, res.Body.Close())
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperClosesRetryableResponseBeforeNextAttempt(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		calls++
		if calls == 1 {
			res.WriteHeader(http.StatusServiceUnavailable)
			_, _ = res.Write([]byte("first failure"))
			return
		}

		res.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	transport := http.Transport(nil)
	transport.MaxConnsPerHost = 1
	t.Cleanup(transport.CloseIdleConnections)

	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  100 * time.Millisecond,
		Backoff:  time.Millisecond,
	}, transport)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, http.NoBody)
	require.NoError(t, err)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, res.Body.Close())
	require.Equal(t, 2, calls)
}

func TestRoundTripperReplaysRequestBodyAcrossRetries(t *testing.T) {
	rt := &test.RequestBodyRecorderRoundTripper{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", strings.NewReader("hello"))
	require.NoError(t, err)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, []string{"hello", "hello"}, rt.Bodies)
}

func TestRoundTripperReplaysRequestBodyAfterTransportError(t *testing.T) {
	rt := &test.TransportErrorThenSuccessRoundTripper{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", strings.NewReader("hello"))
	require.NoError(t, err)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, []string{"hello", "hello"}, rt.Bodies)
}

func TestRoundTripperUsesOriginalBodyOnFirstAttempt(t *testing.T) {
	body := &test.TrackedBody{Reader: strings.NewReader("hello")}
	rt := &test.OriginalBodyRoundTripper{Original: body}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", body)
	require.NoError(t, err)
	req.Body = body
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("hello")), nil
	}

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.True(t, rt.FirstUsedOriginal)
	require.True(t, body.Closed)
	require.Equal(t, []string{"hello", "hello"}, rt.Bodies)
}

func TestRoundTripperDoesNotRetryNonReplayableRequestBody(t *testing.T) {
	rt := &test.RequestBodyRecorderRoundTripper{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", &test.NonReplayableReader{Value: "hello"})
	require.NoError(t, err)
	require.Nil(t, req.GetBody)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, []string{"hello"}, rt.Bodies)
}

func TestRoundTripperPreservesRetryableStatusError(t *testing.T) {
	rt := &test.ErrorRoundTripper{Err: status.Errorf(http.StatusTooManyRequests, "limiter: too many requests")}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.Equal(t, http.StatusTooManyRequests, status.Code(err))
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperPreservesRecoverableTransportError(t *testing.T) {
	wantErr := io.ErrUnexpectedEOF
	rt := &test.ErrorRoundTripper{Err: wantErr}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperDoesNotRetryUnhandledTransportError(t *testing.T) {
	rt := &test.ErrorRoundTripper{Err: status.Errorf(http.StatusInternalServerError, "internal")}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, status.Code(err))
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperDoesNotAccumulateAuthorizationHeadersAcrossRetries(t *testing.T) {
	rt := &test.AuthRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	generator := &test.SequenceGenerator{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, token.NewRoundTripper(env.UserID("user-id"), generator, rt))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, []string{"Bearer token-1", "Bearer token-2"}, rt.AuthValues)
	require.Equal(t, []int{1, 1}, rt.AuthCounts)
}

func TestRoundTripperDoesNotSetAttemptTimeoutWhenTimeoutIsZero(t *testing.T) {
	transport := &test.CauseRoundTripper{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 1,
		Timeout:  0,
		Backoff:  time.Millisecond,
	}, transport)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NoError(t, transport.Err)
	require.NoError(t, transport.Cause)
}

func TestRoundTripperSetsAttemptTimeoutCause(t *testing.T) {
	transport := &test.CauseRoundTripper{Wait: true}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 1,
		Timeout:  time.Nanosecond,
		Backoff:  time.Millisecond,
	}, transport)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.ErrorIs(t, transport.Err, context.DeadlineExceeded)
	require.ErrorIs(t, transport.Cause, retry.ErrAttemptTimeout)
	require.ErrorIs(t, transport.Cause, sync.ErrTimeout)
}
