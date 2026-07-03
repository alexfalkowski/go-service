package telemetry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestRegisterPropagationUsesDefaults(t *testing.T) {
	propagator, err := telemetry.NewPropagator(nil)
	require.NoError(t, err)
	telemetry.RegisterPropagation(propagator)

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(remoteContext(t), carrier)

	require.Contains(t, carrier, "traceparent")
	require.NotContains(t, carrier, "b3")
	require.ElementsMatch(t, []string{"baggage", "traceparent", "tracestate"}, otel.GetTextMapPropagator().Fields())
}

func TestRegisterPropagationUsesConfiguredB3Formats(t *testing.T) {
	propagator, err := telemetry.NewPropagator(&telemetry.PropagationConfig{
		Formats: []string{"tracecontext", "baggage", "b3"},
	})
	require.NoError(t, err)
	telemetry.RegisterPropagation(propagator)

	ctx := otel.GetTextMapPropagator().Extract(t.Context(), propagation.MapCarrier{
		"b3": "00000000000000000000000000000001-0000000000000002-1",
	})
	sc := trace.SpanContextFromContext(ctx)
	require.True(t, sc.IsValid())
	require.True(t, sc.IsRemote())
	require.True(t, sc.IsSampled())

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(remoteContext(t), carrier)
	require.Contains(t, carrier, "traceparent")
	require.Contains(t, carrier, "b3")
	require.NotContains(t, carrier, "x-b3-traceid")
}

func TestRegisterPropagationUsesB3Multi(t *testing.T) {
	propagator, err := telemetry.NewPropagator(&telemetry.PropagationConfig{
		Formats: []string{"b3multi"},
	})
	require.NoError(t, err)
	telemetry.RegisterPropagation(propagator)

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(remoteContext(t), carrier)

	require.Contains(t, carrier, "x-b3-traceid")
	require.NotContains(t, carrier, "b3")
	require.NotContains(t, carrier, "traceparent")
}

func TestRegisterPropagationSupportsNone(t *testing.T) {
	propagator, err := telemetry.NewPropagator(&telemetry.PropagationConfig{
		Formats: []string{"none"},
	})
	require.NoError(t, err)
	telemetry.RegisterPropagation(propagator)

	ctx := otel.GetTextMapPropagator().Extract(t.Context(), propagation.MapCarrier{
		"traceparent": "00-00000000000000000000000000000001-0000000000000002-01",
	})
	require.False(t, trace.SpanContextFromContext(ctx).IsValid())

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(remoteContext(t), carrier)
	require.Empty(t, carrier)
	require.Empty(t, otel.GetTextMapPropagator().Fields())
}

func TestRegisterPropagationRejectsInvalidConfig(t *testing.T) {
	tests := map[string]*telemetry.PropagationConfig{
		"unknown":    {Formats: []string{"unknown"}},
		"mixed none": {Formats: []string{"none", "tracecontext"}},
	}

	for name, config := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := telemetry.NewPropagator(config)

			require.ErrorIs(t, err, telemetry.ErrInvalidPropagator)
		})
	}
}

func remoteContext(t *testing.T) context.Context {
	t.Helper()

	traceID, err := trace.TraceIDFromHex("00000000000000000000000000000001")
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex("0000000000000002")
	require.NoError(t, err)

	return trace.ContextWithRemoteSpanContext(t.Context(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	}))
}
