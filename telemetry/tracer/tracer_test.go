package tracer_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestIsEnabled(t *testing.T) {
	t.Cleanup(func() {
		require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	})

	tracer.SetProvider(tracer.NewNoopProvider())
	require.False(t, tracer.IsEnabled())

	require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	require.False(t, tracer.IsEnabled())

	require.NoError(t, tracer.Register(tracer.TracerParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    &tracer.Config{},
	}))
	require.False(t, tracer.IsEnabled())

	require.NoError(t, tracer.Register(tracer.TracerParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &tracer.Config{Kind: "otlp", URL: "https://localhost:4318/v1/traces"},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}))
	require.True(t, tracer.IsEnabled())

	require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	require.False(t, tracer.IsEnabled())
}

func TestConfigIsEnabled(t *testing.T) {
	require.False(t, (*tracer.Config)(nil).IsEnabled())
	require.False(t, (&tracer.Config{}).IsEnabled())
	require.True(t, (&tracer.Config{Kind: "otlp"}).IsEnabled())
}

func TestRegisterStopResetsGlobalProvider(t *testing.T) {
	t.Cleanup(func() {
		require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	})

	tracer.SetProvider(tracer.NewNoopProvider())
	lc := fxtest.NewLifecycle(t)
	require.NoError(t, tracer.Register(tracer.TracerParams{
		Lifecycle:   lc,
		Config:      &tracer.Config{Kind: "otlp", URL: "https://localhost:4318/v1/traces"},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}))
	require.True(t, tracer.IsEnabled())

	lc.RequireStart()
	require.NoError(t, lc.Stop(t.Context()))

	require.False(t, tracer.IsEnabled())
}

func TestRegisterInvalidKind(t *testing.T) {
	err := tracer.Register(tracer.TracerParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    &tracer.Config{Kind: "wrong"},
	})

	require.ErrorIs(t, err, tracer.ErrNotFound)
}

func TestRegisterInvalidOTLPEndpoint(t *testing.T) {
	err := tracer.Register(tracer.TracerParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config: &tracer.Config{
			Kind: "otlp",
			URL:  "http://collector.example.com/v1/traces",
			Headers: header.Map{
				"Authorization": "Bearer token",
			},
		},
	})

	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestRegisterOTLPEndpointLoopbackIPs(t *testing.T) {
	headers := header.Map{"Authorization": "Bearer token"}
	tests := []struct {
		wantErr error
		name    string
		url     string
	}{
		{name: "IPv6 loopback", url: "http://[::1]:4318/v1/traces"},
		{name: "private IPv4", url: "http://10.0.0.10:4318/v1/traces", wantErr: otlp.ErrInsecureEndpoint},
		{name: "unique local IPv6", url: "http://[fd00::1]:4318/v1/traces", wantErr: otlp.ErrInsecureEndpoint},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := fxtest.NewLifecycle(t)
			err := tracer.Register(tracer.TracerParams{
				Lifecycle: lc,
				Config: &tracer.Config{
					Kind:    "otlp",
					URL:     tt.url,
					Headers: headers,
				},
				ID:          test.ID,
				Name:        test.Name,
				Version:     test.Version,
				Environment: test.Environment,
			})
			if tt.wantErr == nil {
				require.NoError(t, err)
				lc.RequireStart()
				require.NoError(t, lc.Stop(t.Context()))
				return
			}

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRegisterMissingOTLPEndpointIgnoresEnv(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "http://collector.example.com/v1/traces")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://collector.example.com")

	err := tracer.Register(tracer.TracerParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config: &tracer.Config{
			Kind: "otlp",
			Headers: header.Map{
				"Authorization": "Bearer token",
			},
		},
	})

	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}
