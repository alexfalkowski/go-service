package otlp_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/stretchr/testify/require"
)

func TestValidateEndpoint(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	require.NoError(t, otlp.ValidateEndpoint("http", "https://collector.example.com/v1/traces", headers))
	require.NoError(t, otlp.ValidateEndpoint("http", "http://localhost:4318/v1/traces", headers))
	require.NoError(t, otlp.ValidateEndpoint("http", "http://127.0.0.1:4318/v1/traces", headers))
	require.NoError(t, otlp.ValidateEndpoint("http", "http://collector.example.com/v1/traces", nil))
	require.NoError(t, otlp.ValidateEndpoint("grpc", "localhost:4317", headers))
	require.NoError(t, otlp.ValidateEndpoint("grpc", "collector.example.com:4317", nil))

	err := otlp.ValidateEndpoint("http", "http://collector.example.com/v1/traces", headers)
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)

	err = otlp.ValidateEndpoint("grpc", "collector.example.com:4317", headers)
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestValidateEndpointInvalidURL(t *testing.T) {
	err := otlp.ValidateEndpoint("http", "http://%", map[string]string{"Authorization": "Bearer token"})

	require.Error(t, err)
	require.NotErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestValidateEndpointRejectsInvalidEndpoint(t *testing.T) {
	for _, rawURL := range []string{
		"htps://collector.example.com/v1/traces",
		"https:///v1/traces",
	} {
		t.Run(rawURL, func(t *testing.T) {
			err := otlp.ValidateEndpoint("http", rawURL, nil)

			require.ErrorIs(t, err, otlp.ErrInvalidEndpoint)
		})
	}

	for _, endpoint := range []string{
		"https://collector.example.com:4317",
		"collector.example.com",
		"collector.example.com:4317/v1/traces",
	} {
		t.Run(endpoint, func(t *testing.T) {
			err := otlp.ValidateEndpoint("grpc", endpoint, nil)

			require.ErrorIs(t, err, otlp.ErrInvalidEndpoint)
		})
	}
}

func TestValidateEndpointRequiresEndpoint(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	require.ErrorIs(t, otlp.ValidateEndpoint("http", "", headers), otlp.ErrMissingEndpoint)
	require.NoError(t, otlp.ValidateEndpoint("http", "https://collector.example.com/v1/traces", headers))

	err := otlp.ValidateEndpoint("http", "http://collector.example.com/v1/traces", headers)
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestValidateEndpointRejectsInvalidProtocol(t *testing.T) {
	err := otlp.ValidateEndpoint("wrong", "collector.example.com:4317", nil)

	require.ErrorIs(t, err, otlp.ErrInvalidProtocol)
}

func TestValidateEndpointIgnoresEnv(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "https://collector.example.com/v1/traces")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://collector.example.com")

	err := otlp.ValidateEndpoint("http", "", map[string]string{"Authorization": "Bearer token"})
	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}
