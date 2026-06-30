package url_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/stretchr/testify/require"
)

// FuzzSplitPath explores slash-prefixed service/method path parsing used by HTTP and gRPC telemetry.
func FuzzSplitPath(f *testing.F) {
	for _, path := range []string{
		"/service/method",
		"/service/users/123",
		"service/method",
		"/service",
		"//method",
		"/service/",
		"/",
		"",
	} {
		f.Add(path)
	}

	f.Fuzz(func(t *testing.T, path string) {
		first, rest, ok := url.SplitPath(path)
		if !ok {
			require.Empty(t, first)
			require.Empty(t, rest)
			return
		}

		require.NotEmpty(t, first)
		require.NotEmpty(t, rest)
		require.Equal(t, "/"+first+"/"+rest, path)
	})
}
