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
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsInvalidStatusCodes(t *testing.T) {
	require.NoError(t, test.Validator.Struct(&retry.Config{StatusCodes: []int{http.StatusBadRequest, http.StatusNetworkAuthenticationRequired}}))
	require.Error(t, test.Validator.Struct(&retry.Config{StatusCodes: []int{http.StatusContinue}}))
	require.Error(t, test.Validator.Struct(&retry.Config{StatusCodes: []int{600}}))
}

func TestRoundTripperRetriesRetryableResponses(t *testing.T) {
	tests := []struct {
		name   string
		method string
		codes  []int
		calls  int
	}{
		{name: "get too many requests", method: http.MethodGet, codes: []int{http.StatusTooManyRequests, http.StatusOK}, calls: 2},
		{name: "get service unavailable", method: http.MethodGet, codes: []int{http.StatusServiceUnavailable, http.StatusOK}, calls: 2},
		{name: "head too many requests", method: http.MethodHead, codes: []int{http.StatusTooManyRequests, http.StatusOK}, calls: 2},
		{name: "head service unavailable", method: http.MethodHead, codes: []int{http.StatusServiceUnavailable, http.StatusOK}, calls: 2},
		{name: "options too many requests", method: http.MethodOptions, codes: []int{http.StatusTooManyRequests, http.StatusOK}, calls: 2},
		{name: "options service unavailable", method: http.MethodOptions, codes: []int{http.StatusServiceUnavailable, http.StatusOK}, calls: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &test.StatusSequenceRoundTripper{Codes: tt.codes}
			retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

			req := httptest.NewRequestWithContext(t.Context(), tt.method, "http://example.com", http.NoBody)

			res, err := retrying.RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, tt.calls, rt.Calls)
		})
	}
}

func TestRoundTripperRetriesConfiguredStatusCodes(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusBadGateway, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond, http.StatusBadGateway), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperConfiguredStatusCodesReplaceDefaults(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond, http.StatusBadGateway), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperDoesNotRetryWhenAttemptsIsOne(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(1, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperClampsAttemptsAboveMax(t *testing.T) {
	calls := 0
	rt := test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		calls++
		return test.ResponseWithStatus(http.StatusTooManyRequests), nil
	})
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(config.MaxAttempts+1, time.Nanosecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, int(config.MaxAttempts), calls)
}

func TestRoundTripperDoesNotPanicWithOmittedBackoff(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(1, 0), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	require.NotPanics(t, func() {
		res, err := retrying.RoundTrip(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	})
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperRetriesWithDefaultBackoff(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, 0), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperDoesNotRetryWhenRetryAfterExceedsMinimumBackoff(t *testing.T) {
	retryAfterDate := time.Now().Add(5 * time.Second.Duration()).UTC().Format(http.TimeFormat)
	tests := []struct {
		name       string
		retryAfter string
		code       int
		backoff    time.Duration
	}{
		{name: "too many requests seconds", code: http.StatusTooManyRequests, retryAfter: "2", backoff: time.Second},
		{name: "too many requests equal backoff", code: http.StatusTooManyRequests, retryAfter: "1", backoff: time.Second},
		{name: "service unavailable date", code: http.StatusServiceUnavailable, retryAfter: retryAfterDate, backoff: time.Second},
		{name: "too many requests overflow", code: http.StatusTooManyRequests, retryAfter: "999999999999999999999999999999", backoff: time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			rt := test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
				calls++
				res := test.ResponseWithStatus(tt.code)
				res.Header.Set("Retry-After", tt.retryAfter)

				return res, nil
			})
			retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, tt.backoff), rt)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

			res, err := retrying.RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, tt.code, res.StatusCode)
			require.Equal(t, 1, calls)
		})
	}
}

func TestRoundTripperRetriesWhenRetryAfterDoesNotExceedBackoff(t *testing.T) {
	calls := 0
	rt := test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			res := test.ResponseWithStatus(http.StatusTooManyRequests)
			res.Header.Set("Retry-After", "1")

			return res, nil
		}

		return test.ResponseWithStatus(http.StatusOK), nil
	})
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, 2*time.Second), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 2, calls)
}

func TestRoundTripperRetriesWhenRetryAfterIsInvalid(t *testing.T) {
	calls := 0
	rt := test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			res := test.ResponseWithStatus(http.StatusTooManyRequests)
			res.Header.Set("Retry-After", "invalid")

			return res, nil
		}

		return test.ResponseWithStatus(http.StatusOK), nil
	})
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 2, calls)
}

func TestRoundTripperUsesIndependentRetryBudgetPerRequest(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{
		http.StatusTooManyRequests, http.StatusOK,
		http.StatusTooManyRequests, http.StatusOK,
	}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

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
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperDoesNotRetryWhenPolicyDeniesRequest(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt, retry.SafeMethods)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperDoesNotRetryUnsafeRequestByDefault(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperRetriesWhenPolicyAllowsRequestID(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt, retry.IdempotentRequests)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperReturnsLastRetryableResponseWhenExhausted(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusTooManyRequests}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperReturnsFinalRetryableResponseWhenExhausted(t *testing.T) {
	rt := &test.BodySequenceRoundTripper{Responses: []string{"first failure", "second failure"}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

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

	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), transport)

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
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", strings.NewReader("hello"))
	require.NoError(t, err)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, []string{"hello", "hello"}, rt.Bodies)
}

func TestRoundTripperReplaysRequestBodyAfterTransportError(t *testing.T) {
	rt := &test.TransportErrorThenSuccessRoundTripper{}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", strings.NewReader("hello"))
	require.NoError(t, err)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, []string{"hello", "hello"}, rt.Bodies)
}

func TestRoundTripperUsesOriginalBodyOnFirstAttempt(t *testing.T) {
	body := &test.TrackedBody{Reader: strings.NewReader("hello")}
	rt := &test.OriginalBodyRoundTripper{Original: body}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", body)
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
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", &test.NonReplayableReader{Value: "hello"})
	require.NoError(t, err)
	require.Nil(t, req.GetBody)

	res, roundTripErr := retrying.RoundTrip(req)
	require.NoError(t, roundTripErr)
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, []string{"hello"}, rt.Bodies)
}

func TestRoundTripperReturnsGetBodyError(t *testing.T) {
	wantErr := test.ErrFailed
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", strings.NewReader("hello"))
	require.NoError(t, err)
	req.GetBody = func() (io.ReadCloser, error) {
		return nil, wantErr
	}

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperClosesBodyWhenContextAlreadyCanceled(t *testing.T) {
	rt := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusOK}}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	ctx := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("request-id")))
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	body := &test.TrackedBody{Reader: strings.NewReader("hello")}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", body)
	require.NoError(t, err)
	req.Body = body

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 0, rt.Calls)
	require.True(t, body.Closed)
}

func TestRoundTripperPreservesRetryableStatusError(t *testing.T) {
	rt := &test.ErrorRoundTripper{Err: status.Errorf(http.StatusTooManyRequests, "limiter: too many requests")}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.Equal(t, http.StatusTooManyRequests, status.Code(err))
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperDoesNotRetryLocalStatusError(t *testing.T) {
	rt := &test.ErrorRoundTripper{Err: status.LocalError(status.Errorf(http.StatusTooManyRequests, "limiter: too many requests"))}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, http.StatusTooManyRequests, status.Code(err))
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperDoesNotRetryUseLastResponse(t *testing.T) {
	rt := &test.ErrorRoundTripper{Err: http.ErrUseLastResponse}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, http.ErrUseLastResponse)
	require.Equal(t, 1, rt.Calls)
}

func TestRoundTripperPreservesRecoverableTransportError(t *testing.T) {
	wantErr := io.ErrUnexpectedEOF
	rt := &test.ErrorRoundTripper{Err: wantErr}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 2, rt.Calls)
}

func TestRoundTripperDoesNotRetryUnhandledTransportError(t *testing.T) {
	rt := &test.ErrorRoundTripper{Err: status.Errorf(http.StatusInternalServerError, "internal")}
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), rt)

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
	retrying := retry.NewRoundTripper(test.NewHTTPRetryConfig(2, time.Millisecond), token.NewRoundTripper(env.UserID("user-id"), generator, rt))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)

	res, err := retrying.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, []string{"Bearer token-1", "Bearer token-2"}, rt.AuthValues)
	require.Equal(t, []int{1, 1}, rt.AuthCounts)
}
