package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestRouterRegistersHandlerAndPolicy(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)

	router.HandleRoute("POST /events/{tenant}", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusAccepted)
	}), http.WithRouteUnauthenticated())

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events/acme", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, "POST /events/{tenant}", req.Pattern)
	require.True(t, policy.IsUnauthenticated(req))
	require.Equal(t, http.StatusAccepted, res.Code)
}

func TestRouterAppliesRouteOptions(t *testing.T) {
	tests := []struct {
		name              string
		options           []http.RouteOption
		operation         bool
		unauthenticated   bool
		requestStreaming  bool
		responseStreaming bool
	}{
		{name: "operation", options: []http.RouteOption{http.WithRouteOperation()}, operation: true},
		{name: "unauthenticated", options: []http.RouteOption{http.WithRouteUnauthenticated()}, unauthenticated: true},
		{name: "request streaming", options: []http.RouteOption{http.WithRouteRequestStreaming()}, requestStreaming: true},
		{name: "response streaming", options: []http.RouteOption{http.WithRouteResponseStreaming()}, responseStreaming: true},
		{name: "streaming", options: []http.RouteOption{http.WithRouteStreaming()}, requestStreaming: true, responseStreaming: true},
		{name: "combined streaming", options: []http.RouteOption{http.WithRouteRequestStreaming(), http.WithRouteResponseStreaming()}, requestStreaming: true, responseStreaming: true},
		{
			name:              "combined policy and streaming",
			options:           []http.RouteOption{http.WithRouteOperation(), http.WithRouteUnauthenticated(), http.WithRouteResponseStreaming()},
			operation:         true,
			unauthenticated:   true,
			responseStreaming: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			policy := http.NewRoutePolicy()
			router := http.NewRouter(mux, policy)
			pattern := "POST /stream/{id}"
			path := "/stream/acme"
			if tt.operation {
				pattern = "POST /stream"
				path = "/stream"
			}

			router.HandleRoute(pattern, http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
				res.WriteHeader(http.StatusOK)
			}), tt.options...)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, path, http.NoBody)
			res := httptest.NewRecorder()
			mux.ServeHTTP(res, req)

			require.Equal(t, pattern, req.Pattern)
			require.Equal(t, tt.operation, policy.IsOperation(req))
			require.Equal(t, tt.unauthenticated, policy.IsUnauthenticated(req))
			require.Equal(t, tt.requestStreaming, policy.IsRequestStreaming(req))
			require.Equal(t, tt.responseStreaming, policy.IsResponseStreaming(req))
			require.Equal(t, http.StatusOK, res.Code)
		})
	}
}

func TestRouterMatchesParameterizedOperationPath(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)
	router.HandleRoute("GET /operations/{id}", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}), http.WithRouteOperation())

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/operations/ready", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusNoContent, res.Code)
	require.True(t, policy.IsOperation(req))
}

func TestRouterRetainsLiteralOperationFallbackAfterNonOperationMatch(t *testing.T) {
	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	router := http.NewRouter(mux, policy)
	router.HandleRoute("GET /health", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), http.WithRouteOperation())
	router.HandleRoute("POST /{id}", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/health", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, "POST /{id}", req.Pattern)
	require.Equal(t, http.StatusNoContent, res.Code)
	require.True(t, policy.IsOperation(req))
}

func TestRoutePolicyFallsBackForEachRouteProperty(t *testing.T) {
	tests := []struct {
		name     string
		options  []http.RouteOption
		check    func(*http.RoutePolicy, *http.Request) bool
		fallback http.RouteOption
	}{
		{
			name:     "unauthenticated",
			options:  []http.RouteOption{http.WithRouteResponseStreaming()},
			check:    (*http.RoutePolicy).IsUnauthenticated,
			fallback: http.WithRouteUnauthenticated(),
		},
		{
			name:     "response streaming",
			options:  []http.RouteOption{http.WithRouteUnauthenticated()},
			check:    (*http.RoutePolicy).IsResponseStreaming,
			fallback: http.WithRouteResponseStreaming(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := http.NewRoutePolicy()
			router := http.NewRouter(http.NewServeMux(), policy)
			router.HandleRoute("POST /events", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), tt.options...)
			router.HandleRoute("/events", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), tt.fallback)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", http.NoBody)
			require.True(t, tt.check(policy, req))
		})
	}
}
