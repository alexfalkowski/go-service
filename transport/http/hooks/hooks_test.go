package hooks_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
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

func TestRoundTripper(t *testing.T) {
	hook := hooks.NewWebhook(nil, nil)
	rt := hooks.NewRoundTripper(hook, roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: make(http.Header)}, nil
	}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", &test.ErrReaderCloser{})

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
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

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
