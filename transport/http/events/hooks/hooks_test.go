package hooks_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v3/hooks"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/events"
	eventhooks "github.com/alexfalkowski/go-service/v3/transport/http/events/hooks"
	httphooks "github.com/alexfalkowski/go-service/v3/transport/http/hooks"
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
	webhook, err := hooks.NewHook(test.FS, test.NewHook())
	require.NoError(t, err)

	called := false
	handler := eventhooks.NewHandler(
		httphooks.NewWebhook(webhook, nil),
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			called = true
		}),
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com", http.NoBody)
	req.Header.Set(events.SpecVersionKey, "1.0")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.False(t, called)
	require.Equal(t, http.StatusBadRequest, res.Code)
}
