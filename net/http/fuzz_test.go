package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

// FuzzParseServiceMethod explores HTTP telemetry service/method derivation for arbitrary paths and methods.
func FuzzParseServiceMethod(f *testing.F) {
	for _, seed := range []struct {
		path   string
		method string
	}{
		{path: "/test/hello", method: http.MethodGet},
		{path: "/test/users/123", method: http.MethodPost},
		{path: "/", method: http.MethodGet},
		{path: "/health", method: http.MethodPost},
		{path: "", method: http.MethodGet},
		{path: "//hello", method: http.MethodPost},
		{path: "/test/", method: http.MethodPut},
	} {
		f.Add(seed.path, seed.method)
	}

	f.Fuzz(func(t *testing.T, path, method string) {
		if len(path) > 1024 || len(method) > 128 {
			t.Skip()
		}

		req := &http.Request{Method: method, URL: &url.URL{Path: path}}
		service, action := http.ParseServiceMethod(req)
		splitService, splitMethod, ok := url.SplitPath(path)
		if ok {
			require.Equal(t, splitService, service)
			require.Equal(t, splitMethod, action)
			return
		}

		require.Equal(t, strings.ToLower(method), action)
		if strings.IsEmpty(path) || path == "/" {
			require.Equal(t, "root", service)
			return
		}
		if path[0] == '/' {
			require.Equal(t, path[1:], service)
		}
	})
}
