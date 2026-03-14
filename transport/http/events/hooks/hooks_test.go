package hooks_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/events/hooks"
	"github.com/stretchr/testify/require"
)

func TestHandlerWithNilHook(t *testing.T) {
	called := false
	handler := hooks.NewHandler(nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
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
