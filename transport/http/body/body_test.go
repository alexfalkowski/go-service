package body_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/body"
	"github.com/stretchr/testify/require"
)

func TestHandlerSkipsEmptyBodyBuffering(t *testing.T) {
	handler := body.NewHandler(1024)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.NoBody, req.Body)
		res.WriteHeader(http.StatusOK)
	})

	require.Equal(t, http.StatusOK, res.Code)
}

func TestHandlerWritesBadRequestWhenBodyReadFails(t *testing.T) {
	handler := body.NewHandler(1024)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/test", &test.ErrReaderCloser{})
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {
		require.Fail(t, "next handler should not be called")
	})

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, test.ErrFailed.Error()+"\n", res.Body.String())
}
