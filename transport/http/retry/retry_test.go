package retry_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/retry"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
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
			rt := &roundTripper{codes: tt.codes}
			retrying := retry.NewRoundTripper(&retry.Config{
				Attempts: 2,
				Timeout:  time.Second,
				Backoff:  time.Millisecond,
			}, rt)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

			res, err := retrying.RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, tt.calls, rt.calls)
		})
	}
}

func TestRoundTripperDoesNotRetryWhenAttemptsIsOne(t *testing.T) {
	rt := &roundTripper{codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 1,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, 1, rt.calls)
}

func TestRoundTripperReturnsLastRetryableResponseWhenExhausted(t *testing.T) {
	rt := &roundTripper{codes: []int{http.StatusTooManyRequests, http.StatusTooManyRequests}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, 2, rt.calls)
}

func TestRoundTripperReturnsFirstRetryableResponseWhenExhausted(t *testing.T) {
	rt := &bodyRoundTripper{responses: []string{"first failure", "second failure"}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	body, readErr := io.ReadAll(res.Body)
	require.NoError(t, readErr)
	require.Equal(t, "first failure", string(body))
	require.NoError(t, res.Body.Close())
	require.Equal(t, 2, rt.calls)
}

func TestRoundTripperReplaysRequestBodyAcrossRetries(t *testing.T) {
	rt := &requestRoundTripper{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", strings.NewReader("hello"))
	require.NoError(t, err)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, []string{"hello", "hello"}, rt.bodies)
}

func TestRoundTripperDoesNotRetryNonReplayableRequestBody(t *testing.T) {
	rt := &requestRoundTripper{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, rt)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", &nonReplayableReader{value: "hello"})
	require.NoError(t, err)
	require.Nil(t, req.GetBody)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, []string{"hello"}, rt.bodies)
}

func TestRoundTripperPreservesRetryableTransportError(t *testing.T) {
	rt := &errorRoundTripper{err: status.Errorf(http.StatusTooManyRequests, "limiter: too many requests")}
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
	require.Equal(t, 2, rt.calls)
}

func TestRoundTripperDoesNotAccumulateAuthorizationHeadersAcrossRetries(t *testing.T) {
	rt := &authRoundTripper{codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	generator := &tokenGenerator{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 2,
		Timeout:  time.Second,
		Backoff:  time.Millisecond,
	}, token.NewRoundTripper(env.UserID("user-id"), generator, rt))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, []string{"Bearer token-1", "Bearer token-2"}, rt.authValues)
	require.Equal(t, []int{1, 1}, rt.authCounts)
}

func TestRoundTripperSetsAttemptTimeoutCause(t *testing.T) {
	transport := &causeRoundTripper{}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 1,
		Timeout:  0,
		Backoff:  time.Millisecond,
	}, transport)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.ErrorIs(t, transport.err, context.DeadlineExceeded)
	require.ErrorIs(t, transport.cause, retry.ErrAttemptTimeout)
}

type roundTripper struct {
	codes []int
	calls int
}

func (r *roundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	code := r.codes[r.calls]
	r.calls++

	return &http.Response{
		StatusCode: code,
		Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
		Body:       http.NoBody,
		Header:     make(http.Header),
	}, nil
}

type bodyRoundTripper struct {
	responses []string
	calls     int
}

func (r *bodyRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	body := r.responses[r.calls]
	r.calls++

	return &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Status:     fmt.Sprintf("%d %s", http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}, nil
}

type requestRoundTripper struct {
	bodies []string
}

func (r *requestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	r.bodies = append(r.bodies, string(body))

	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)),
		Body:       http.NoBody,
		Header:     make(http.Header),
	}, nil
}

type errorRoundTripper struct {
	err   error
	calls int
}

func (r *errorRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	r.calls++
	return nil, r.err
}

type authRoundTripper struct {
	authValues []string
	authCounts []int
	codes      []int
	calls      int
}

func (r *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.authValues = append(r.authValues, req.Header.Get("Authorization"))
	r.authCounts = append(r.authCounts, len(req.Header.Values("Authorization")))
	code := r.codes[r.calls]
	r.calls++

	return &http.Response{
		StatusCode: code,
		Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
		Body:       http.NoBody,
		Header:     make(http.Header),
	}, nil
}

type causeRoundTripper struct {
	cause error
	err   error
}

func (r *causeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.cause = context.Cause(req.Context())
	r.err = req.Context().Err()

	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Body:       http.NoBody,
		Header:     make(http.Header),
	}, nil
}

type tokenGenerator struct {
	calls int
}

func (g *tokenGenerator) Generate(_, _ string) ([]byte, error) {
	g.calls++
	return fmt.Appendf(nil, "token-%d", g.calls), nil
}

type nonReplayableReader struct {
	value string
	read  bool
}

func (r *nonReplayableReader) Read(p []byte) (int, error) {
	if r.read {
		return 0, io.EOF
	}

	r.read = true
	copy(p, r.value)
	return len(r.value), io.EOF
}
