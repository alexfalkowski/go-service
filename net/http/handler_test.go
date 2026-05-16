package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestNewNotFoundHandler(t *testing.T) {
	tests := []notFoundHandlerTest{
		{
			name:   "handles mux not found",
			method: http.MethodGet,
			path:   "/missing",
			handle: true,
			code:   http.StatusNotFound,
			body:   "fallback",
		},
		{
			name:     "flushes mux not found when unhandled",
			method:   http.MethodGet,
			path:     "/missing",
			code:     http.StatusNotFound,
			contains: "404 page not found",
		},
		{
			name:          "does not handle method not allowed",
			method:        http.MethodPost,
			path:          "/hello",
			registerRoute: true,
			handle:        true,
			code:          http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testNotFoundHandler(t, tt)
		})
	}
}

type notFoundHandlerTest struct {
	name          string
	method        string
	path          string
	body          string
	contains      string
	registerRoute bool
	handle        bool
	code          int
}

func testNotFoundHandler(t *testing.T, tt notFoundHandlerTest) {
	t.Helper()

	mux := http.NewServeMux()
	if tt.registerRoute {
		http.HandleFunc(mux, "GET /hello", func(res http.ResponseWriter, _ *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	}

	handler := http.NewNotFoundHandler(mux, test.Pool, func(res http.ResponseWriter, _ *http.Request) bool {
		if !tt.handle {
			return false
		}

		res.WriteHeader(http.StatusNotFound)
		_, _ = res.Write([]byte("fallback"))
		return true
	})

	req := httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, http.NoBody)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, tt.code, res.Code)
	if tt.body != "" {
		require.Equal(t, tt.body, res.Body.String())
	}
	if tt.contains != "" {
		require.Contains(t, res.Body.String(), tt.contains)
	}
}
