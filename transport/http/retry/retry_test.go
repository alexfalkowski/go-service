package retry_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/transport/http/retry"
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
				Timeout:  "1s",
				Backoff:  "1ms",
			}, rt)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

			res, err := retrying.RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, tt.calls, rt.calls)
		})
	}
}

func TestRoundTripperReturnsLastRetryableResponseWhenExhausted(t *testing.T) {
	rt := &roundTripper{codes: []int{http.StatusTooManyRequests, http.StatusTooManyRequests}}
	retrying := retry.NewRoundTripper(&retry.Config{
		Attempts: 1,
		Timeout:  "1s",
		Backoff:  "1ms",
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
		Attempts: 1,
		Timeout:  "1s",
		Backoff:  "1ms",
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
		Attempts: 1,
		Timeout:  "1s",
		Backoff:  "1ms",
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
		Attempts: 1,
		Timeout:  "1s",
		Backoff:  "1ms",
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
		Attempts: 1,
		Timeout:  "1s",
		Backoff:  "1ms",
	}, rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.Equal(t, http.StatusTooManyRequests, status.Code(err))
	require.Equal(t, 2, rt.calls)
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
