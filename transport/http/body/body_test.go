package body_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/http/body"
	"github.com/stretchr/testify/require"
)

func TestHandlerBuffersNonStreamingRoutes(t *testing.T) {
	routePolicy := http.NewRoutePolicy()
	routePolicy.Streaming("POST /stream")
	handler := body.NewHandler(routePolicy, 1024)

	t.Run("skips empty body buffering", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
		require.NoError(t, err)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req, func(res http.ResponseWriter, req *http.Request) {
			require.Equal(t, http.NoBody, req.Body)
			res.WriteHeader(http.StatusOK)
		})

		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("writes bad request when body read fails", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/test", &test.ErrReaderCloser{})
		require.NoError(t, err)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {
			require.Fail(t, "next handler should not be called")
		})

		require.Equal(t, http.StatusBadRequest, res.Code)
		test.RequireTrimmedResponseBody(t, res, "http: bad request")
	})
}

func TestHandlerDoesNotLimitStreamingRoutesMidStream(t *testing.T) {
	routePolicy := http.NewRoutePolicy()
	routePolicy.Streaming("POST /stream")
	handler := body.NewHandler(routePolicy, 64)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/stream", &test.UnknownLengthReader{Reader: strings.NewReader(strings.Repeat("a", 100))})
	require.NoError(t, err)
	res := httptest.NewRecorder()

	var data []byte
	var readErr error
	handler.ServeHTTP(res, req, func(_ http.ResponseWriter, req *http.Request) {
		data, _, readErr = io.ReadAll(req.Body)
	})

	require.NoError(t, readErr)
	require.Len(t, data, 100)
	require.Equal(t, http.StatusOK, res.Code)
}

func TestHandlerRejectsStreamingRouteContentLengthOverLimit(t *testing.T) {
	routePolicy := http.NewRoutePolicy()
	routePolicy.Streaming("POST /stream")
	handler := body.NewHandler(routePolicy, 64)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/stream", strings.NewReader(strings.Repeat("a", 100)))
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {
		require.Fail(t, "next handler should not be called")
	})

	require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)
	test.RequireTrimmedResponseBody(t, res, "http: request entity too large")
}
