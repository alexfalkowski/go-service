package debug_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

var paths = []string{
	"debug/statsviz",
	"debug/pprof/",
	"debug/pprof/cmdline",
	"debug/pprof/symbol",
	"debug/pprof/trace",
	"debug/psutil",
}

func TestInsecureDebug(t *testing.T) {
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldDebug())

			header := http.Header{}
			url := world.NamedDebugURL("http", path)

			res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}

func TestSecureDebug(t *testing.T) {
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldSecure(), test.WithWorldDebug())

			header := http.Header{}
			url := world.NamedDebugURL("https", path)

			res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}

func TestInvalidServer(t *testing.T) {
	cfg := &debug.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
		},
	}
	params := debug.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Config:     cfg,
		FS:         test.FS,
	}

	_, err := debug.NewServer(params)
	require.Error(t, err)
}
