package debug_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/debug"
	debughttp "github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

var debugPaths = []struct {
	name string
	path string
}{
	{name: "serves statsviz", path: "debug/statsviz"},
	{name: "serves pprof index", path: "debug/pprof/"},
	{name: "serves pprof command line", path: "debug/pprof/cmdline"},
	{name: "serves pprof symbol lookup", path: "debug/pprof/symbol"},
	{name: "serves pprof trace", path: "debug/pprof/trace"},
}

func TestDebugServerExposesRegisteredEndpoints(t *testing.T) {
	tests := []struct {
		name    string
		scheme  string
		options []test.WorldOption
	}{
		{name: "insecure", scheme: "http", options: []test.WorldOption{test.WithWorldDebug()}},
		{name: "secure", scheme: "https", options: []test.WorldOption{test.WithWorldSecure(), test.WithWorldDebug()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireDebugEndpoints(t, tt.scheme, tt.options...)
		})
	}
}

func TestNewServerRejectsInvalidTLSConfiguration(t *testing.T) {
	cfg := &debug.Config{
		Config: &server.Config{
			TLS: test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
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

func TestDebugServerRejectsBodiesOverMaxReceiveSize(t *testing.T) {
	mux := debughttp.NewServeMux()
	cfg := &debug.Config{
		Config: &server.Config{
			Address:        test.RandomAddress(),
			MaxReceiveSize: 1,
		},
	}
	lc := fxtest.NewLifecycle(t)
	debugServer, err := debug.NewServer(debug.ServerParams{
		Shutdowner: test.NewShutdowner(),
		Mux:        mux,
		Config:     cfg,
		FS:         test.FS,
	})
	require.NoError(t, err)
	require.NoError(t, debug.Register(debug.RegisterParams{
		Config:    cfg,
		Lifecycle: lc,
		Name:      test.Name,
		Mux:       mux,
	}))

	service := debugServer.GetService()
	service.Start()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		require.NoError(t, service.Stop(ctx))
	})

	_, host, ok := net.SplitNetworkAddress(test.BoundAddress(cfg.Address, service.String()))
	require.True(t, ok)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"http://"+host+"/"+test.Name.String()+"/debug/pprof/symbol",
		bytes.NewBufferString("too large"),
	)
	require.NoError(t, err)

	res, err := http.NewClient(http.DefaultTransport, time.DefaultTimeout).Do(req)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, res.Body.Close())
	})

	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
}

func requireDebugEndpoints(t *testing.T, scheme string, opts ...test.WorldOption) {
	t.Helper()

	opts = append([]test.WorldOption{test.WithWorldTelemetry("otlp")}, opts...)

	for _, endpoint := range debugPaths {
		t.Run(endpoint.name, func(t *testing.T) {
			world := test.NewStartedWorld(t, opts...)

			header := http.Header{}
			url := world.NamedDebugURL(scheme, endpoint.path)

			res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}
