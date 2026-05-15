package tracer_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestIsEnabled(t *testing.T) {
	t.Cleanup(func() {
		require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	})

	require.NoError(t, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)}))
	require.False(t, tracer.IsEnabled())

	require.NoError(t, tracer.Register(tracer.TracerParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    &tracer.Config{},
	}))
	require.False(t, tracer.IsEnabled())

	require.NoError(t, tracer.Register(tracer.TracerParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &tracer.Config{Kind: "otlp"},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}))
	require.True(t, tracer.IsEnabled())
}

func TestConfigIsEnabled(t *testing.T) {
	require.False(t, (*tracer.Config)(nil).IsEnabled())
	require.False(t, (&tracer.Config{}).IsEnabled())
	require.True(t, (&tracer.Config{Kind: "otlp"}).IsEnabled())
}

func TestRegisterInvalidKind(t *testing.T) {
	err := tracer.Register(tracer.TracerParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    &tracer.Config{Kind: "wrong"},
	})

	require.ErrorIs(t, err, tracer.ErrNotFound)
}
