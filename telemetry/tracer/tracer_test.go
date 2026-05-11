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
		tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(t)})
	require.False(t, tracer.IsEnabled())

	tracer.Register(tracer.TracerParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &tracer.Config{},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
	require.True(t, tracer.IsEnabled())
}
