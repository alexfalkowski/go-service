package debug_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/server"
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldDebug())
		world.Register()
		world.RequireStart()

		header := http.Header{}
		url := world.NamedDebugURL("http", path)

		res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		world.RequireStop()
	}
}

func TestSecureDebug(t *testing.T) {
	for _, path := range paths {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldSecure(), test.WithWorldDebug())
		world.Register()
		world.RequireStart()

		header := http.Header{}
		url := world.NamedDebugURL("https", path)

		res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		world.RequireStop()
	}
}

func TestInvalidServer(t *testing.T) {
	cfg := &debug.Config{
		Config: &server.Config{
			Timeout: "5s",
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
