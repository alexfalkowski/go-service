package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http/strings"
	"github.com/stretchr/testify/require"
)

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
		{name: "no leading slash metrics", path: "service/metrics", match: false},
		{name: "trailing slash metrics", path: "/service/metrics/", match: false},
		{name: "double leading slash metrics", path: "//service/metrics", match: false},
		{name: "double trailing slash metrics", path: "/service/metrics//", match: false},
		{name: "double leading and trailing slash metrics", path: "//service/metrics//", match: false},
		{name: "favicon ico", path: "/favicon.ico", match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, strings.IsOperationPath(env.Name("service"), tt.path))
		})
	}
}
