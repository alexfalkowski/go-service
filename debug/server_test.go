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

var debugPaths = []string{
	"debug/statsviz",
	"debug/pprof/",
	"debug/pprof/cmdline",
	"debug/pprof/symbol",
	"debug/pprof/trace",
	"debug/psutil",
}

func TestInsecureDebug(t *testing.T) {
	requireDebugEndpoints(t, "http", test.WithWorldDebug())
}

func TestSecureDebug(t *testing.T) {
	requireDebugEndpoints(t, "https", test.WithWorldSecure(), test.WithWorldDebug())
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

func TestMaxReceiveSize(t *testing.T) {
	mux := debughttp.NewServeMux()
	cfg := &debug.Config{
		Config: &server.Config{
			Address:        test.RandomAddress(),
			Timeout:        5 * time.Second,
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
		Content:   test.Content,
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
	defer res.Body.Close()

	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
}

func requireDebugEndpoints(t *testing.T, scheme string, options ...test.WorldOption) {
	t.Helper()

	options = append([]test.WorldOption{test.WithWorldTelemetry("otlp")}, options...)

	for _, path := range debugPaths {
		t.Run(path, func(t *testing.T) {
			world := test.NewStartedWorld(t, options...)

			header := http.Header{}
			url := world.NamedDebugURL(scheme, path)

			res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}
