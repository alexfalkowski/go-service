package url_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/stretchr/testify/require"
)

func TestSplitPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		path  string
		first string
		rest  string
		ok    bool
	}{
		{name: "pair", path: "/service/method", first: "service", rest: "method", ok: true},
		{name: "deep rest", path: "/service/users/123", first: "service", rest: "users/123", ok: true},
		{name: "missing leading slash", path: "service/method"},
		{name: "missing rest", path: "/service"},
		{name: "empty first", path: "//method"},
		{name: "empty rest", path: "/service/"},
		{name: "root", path: "/"},
		{name: "empty", path: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			first, rest, ok := url.SplitPath(tt.path)
			require.Equal(t, tt.first, first)
			require.Equal(t, tt.rest, rest)
			require.Equal(t, tt.ok, ok)
		})
	}
}
