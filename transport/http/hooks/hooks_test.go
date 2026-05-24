package hooks_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	webhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {
	hook := hooks.NewWebhook(nil, nil)
	req := &http.Request{Body: &test.ErrReaderCloser{}}

	require.NoError(t, hook.Verify(req))
}

func TestSign(t *testing.T) {
	hook := hooks.NewWebhook(nil, nil)
	req := &http.Request{Body: &test.ErrReaderCloser{}}

	require.NoError(t, hook.Sign(req))
}

func TestSignOverwritesHeaders(t *testing.T) {
	webhook, err := webhooks.NewWebhook("whsec_dGVzdA==")
	require.NoError(t, err)

	hook := hooks.NewWebhook(webhook, &test.IDSequenceGenerator{IDs: []string{"id-1", "id-2"}})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)

	require.NoError(t, hook.Sign(req))
	require.NoError(t, hook.Sign(req))
	require.Equal(t, []string{"id-2"}, req.Header.Values(webhooks.HeaderWebhookID))
	require.Len(t, req.Header.Values(webhooks.HeaderWebhookSignature), 1)
	require.Len(t, req.Header.Values(webhooks.HeaderWebhookTimestamp), 1)
	require.NoError(t, hook.Verify(req))
}

func TestSignAndVerifyNilBody(t *testing.T) {
	webhook, err := webhooks.NewWebhook("whsec_dGVzdA==")
	require.NoError(t, err)

	hook := hooks.NewWebhook(webhook, &test.IDSequenceGenerator{IDs: []string{"id-1"}})
	req := &http.Request{Header: make(http.Header)}

	require.NoError(t, hook.Sign(req))
	require.NotNil(t, req.Body)

	verifyReq := &http.Request{Header: req.Header.Clone()}
	require.NoError(t, hook.Verify(verifyReq))
	require.NotNil(t, verifyReq.Body)
}

func TestSignHandlesNilRequestHeader(t *testing.T) {
	webhook, err := webhooks.NewWebhook("whsec_dGVzdA==")
	require.NoError(t, err)

	hook := hooks.NewWebhook(webhook, &test.IDSequenceGenerator{IDs: []string{"id-1"}})
	req := &http.Request{Body: http.NoBody}

	require.NoError(t, hook.Sign(req))
	require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookID))
	require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookSignature))
	require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookTimestamp))
	require.NoError(t, hook.Verify(req))
}

func TestRoundTripper(t *testing.T) {
	hook := hooks.NewWebhook(nil, nil)
	rt := hooks.NewRoundTripper(hook, test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: make(http.Header)}, nil
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", &test.ErrReaderCloser{})

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRoundTripperDoesNotMutateRequest(t *testing.T) {
	webhook, err := webhooks.NewWebhook("whsec_dGVzdA==")
	require.NoError(t, err)

	hook := hooks.NewWebhook(webhook, &test.IDSequenceGenerator{IDs: []string{"id-1"}})
	rt := hooks.NewRoundTripper(hook, test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookID))
		require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookSignature))
		require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookTimestamp))

		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: make(http.Header)}, nil
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Empty(t, req.Header.Values(webhooks.HeaderWebhookID))
	require.Empty(t, req.Header.Values(webhooks.HeaderWebhookSignature))
	require.Empty(t, req.Header.Values(webhooks.HeaderWebhookTimestamp))
}

func TestRoundTripperHandlesNilRequestHeader(t *testing.T) {
	webhook, err := webhooks.NewWebhook("whsec_dGVzdA==")
	require.NoError(t, err)

	hook := hooks.NewWebhook(webhook, &test.IDSequenceGenerator{IDs: []string{"id-1"}})
	rt := hooks.NewRoundTripper(hook, test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookID))
		require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookSignature))
		require.NotEmpty(t, req.Header.Get(webhooks.HeaderWebhookTimestamp))

		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: make(http.Header)}, nil
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)
	req.Header = nil

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Nil(t, req.Header)
}

func TestRoundTripperDoesNotSignCrossOriginRedirect(t *testing.T) {
	webhook, err := webhooks.NewWebhook("whsec_dGVzdA==")
	require.NoError(t, err)

	hook := hooks.NewWebhook(webhook, &test.IDSequenceGenerator{IDs: []string{"id-1"}})
	called := false
	rt := hooks.NewRoundTripper(hook, test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: make(http.Header)}, nil
	}))

	prev := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "https://example.com/events", http.NoBody)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "https://attacker.example.com/events", http.NoBody)
	req.Response = &http.Response{Request: prev}

	res, err := rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, http.ErrUseLastResponse)
	require.False(t, called)
	require.Empty(t, req.Header.Values(webhooks.HeaderWebhookID))
	require.Empty(t, req.Header.Values(webhooks.HeaderWebhookSignature))
	require.Empty(t, req.Header.Values(webhooks.HeaderWebhookTimestamp))
}

func TestHandler(t *testing.T) {
	handler := hooks.NewHandler(nil)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)
	res := httptest.NewRecorder()
	called := false

	handler.ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {
		called = true
	})

	require.True(t, called)
	require.Equal(t, http.StatusOK, res.Code)
}
