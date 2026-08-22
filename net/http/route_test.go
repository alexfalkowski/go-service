package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestParseServiceMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		method  string
		service string
		action  string
	}{
		{name: "service route", method: http.MethodGet, url: "/test/hello", service: "test", action: "hello"},
		{name: "deep service route", method: http.MethodPost, url: "/test/users/123", service: "test", action: "users/123"},
		{name: "root", method: http.MethodGet, url: "/", service: "root", action: "get"},
		{name: "single segment", method: http.MethodPost, url: "/health", service: "health", action: "post"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequestWithContext(t.Context(), test.method, test.url, http.NoBody)
			service, action := http.ParseServiceMethod(req)
			require.Equal(t, test.service, service)
			require.Equal(t, test.action, action)
		})
	}
}
