package url_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/stretchr/testify/require"
)

func TestSplitPath(t *testing.T) {
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			first, rest, ok := url.SplitPath(test.path)
			require.Equal(t, test.first, first)
			require.Equal(t, test.rest, rest)
			require.Equal(t, test.ok, ok)
		})
	}
}
