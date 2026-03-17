package strings_test

import (
	"testing"

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
