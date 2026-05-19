package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http/strings"
	"github.com/stretchr/testify/require"
)

func TestIsIgnorablePath(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		match bool
	}{
		{name: "service healthz", path: "/service/healthz", match: true},
		{name: "service metrics", path: "/service/metrics", match: true},
		{name: "favicon ico", path: "/favicon.ico", match: true},
		{name: "nested exact metrics endpoint", path: "/service/api/metrics", match: false},
		{name: "nested health plan", path: "/v1/customer-health-plans", match: false},
		{name: "metrics substring only", path: "/v1/metrics-dashboard", match: false},
		{name: "health substring only", path: "/v1/healthcare", match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, strings.IsIgnorable(tt.path))
		})
	}
}

func TestIsOperationPath(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		match bool
	}{
		{name: "service healthz", path: "/service/healthz", match: true},
		{name: "service livez", path: "/service/livez", match: true},
		{name: "service readyz", path: "/service/readyz", match: true},
		{name: "service metrics", path: "/service/metrics", match: true},
		{name: "wrong service metrics", path: "/admin/metrics", match: false},
		{name: "wrong service healthz", path: "/admin/healthz", match: false},
		{name: "nested metrics", path: "/service/admin/metrics", match: false},
		{name: "bare metrics", path: "/metrics", match: false},
		{name: "favicon ico", path: "/favicon.ico", match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, strings.IsOperationPath(env.Name("service"), tt.path))
		})
	}
}
