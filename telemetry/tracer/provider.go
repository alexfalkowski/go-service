package tracer

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

// GetProvider returns the global OpenTelemetry tracer provider.
func GetProvider() Provider {
	return otel.GetTracerProvider()
}

// SetProvider installs the global OpenTelemetry tracer provider.
func SetProvider(provider Provider) {
	setProvider(provider, isEnabledProvider(provider))
}

func setProvider(provider Provider, isEnabled bool) {
	otel.SetTracerProvider(provider)
	enabled.Store(isEnabled)
}

func isEnabledProvider(provider Provider) bool {
	if provider == nil {
		return false
	}

	_, ok := provider.(noop.TracerProvider)
	return !ok
}
