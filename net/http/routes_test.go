package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestRoutePolicyClassifiesOperationPaths(t *testing.T) {
	policy := http.NewRoutePolicy()
	policy.Operation("GET /service/metrics")

	tests := []struct {
		name  string
		path  string
		match bool
	}{
		{name: "registered operation path", path: "/service/metrics", match: true},
		{name: "nested metrics path", path: "/service/admin/metrics", match: false},
		{name: "application metrics path", path: "/admin/metrics", match: false},
		{name: "trailing slash", path: "/service/metrics/", match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, tt.path, http.NoBody)

			require.Equal(t, tt.match, policy.IsOperation(req))
		})
	}
}

func TestRoutePolicyClassifiesUnauthenticatedRequests(t *testing.T) {
	policy := http.NewRoutePolicy()
	policy.AllowUnauthenticated("POST /events")

	tests := []struct {
		name   string
		method string
		path   string
		match  bool
	}{
		{name: "registered event receiver", method: http.MethodPost, path: "/events", match: true},
		{name: "wrong method", method: http.MethodGet, path: "/events", match: false},
		{name: "application event path", method: http.MethodPost, path: "/admin/events", match: false},
		{name: "nested event path", method: http.MethodPost, path: "/events/foo", match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, http.NoBody)

			require.Equal(t, tt.match, policy.IsUnauthenticated(req))
		})
	}
}

func TestRouterRegistersHandlerAndPolicy(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)

	router.HandleUnauthenticated("POST /events", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.True(t, policy.IsUnauthenticated(req))
	require.Equal(t, http.StatusAccepted, res.Code)
}
