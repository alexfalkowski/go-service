package tracer

import "go.opentelemetry.io/otel"

// GetProvider returns the global OpenTelemetry tracer provider.
func GetProvider() Provider {
	return otel.GetTracerProvider()
}

// SetProvider installs the global OpenTelemetry tracer provider.
func SetProvider(provider Provider) {
	otel.SetTracerProvider(provider)
}
