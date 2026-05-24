package hooks_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	eventhooks "github.com/alexfalkowski/go-service/v2/transport/http/events/hooks"
	httphooks "github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	webhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"github.com/stretchr/testify/require"
)

func TestHandlerWithNilHook(t *testing.T) {
	called := false
	handler := eventhooks.NewHandler(nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)
	res := httptest.NewRecorder()

	require.NotPanics(t, func() {
		handler.ServeHTTP(res, req)
	})

	require.True(t, called)
	require.Equal(t, http.StatusOK, res.Code)
}

func TestHandlerRejectsBinaryCloudEventsWithWebhook(t *testing.T) {
	webhook, err := webhooks.NewWebhook("whsec_dGVzdA==")
	require.NoError(t, err)

	called := false
	handler := eventhooks.NewHandler(
		httphooks.NewWebhook(webhook, nil),
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			called = true
		}),
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)
	req.Header.Set("Ce-Specversion", "1.0")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.False(t, called)
	require.Equal(t, http.StatusBadRequest, res.Code)
}
