package otlp_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/stretchr/testify/require"
)

func TestValidateEndpoint(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: "https://collector.example.com/v1/traces", Headers: headers}))
	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: "http://localhost:4318/v1/traces", Headers: headers}))
	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: "http://127.0.0.1:4318/v1/traces", Headers: headers}))
	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: "http://collector.example.com/v1/traces"}))
	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "grpc", Address: "localhost:4317", Headers: headers}))
	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "grpc", Address: "collector.example.com:4317"}))
	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "grpc", Address: "collector.example.com:4317", Headers: headers, TLS: &tls.Config{}}))

	err := otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: "http://collector.example.com/v1/traces", Headers: headers})
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)

	err = otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "grpc", Address: "collector.example.com:4317", Headers: headers})
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestValidateEndpointInvalidURL(t *testing.T) {
	err := otlp.ValidateEndpoint(otlp.Endpoint{
		Protocol: "http",
		Address:  "http://%",
		Headers:  map[string]string{"Authorization": "Bearer token"},
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestValidateEndpointRejectsInvalidEndpoint(t *testing.T) {
	for _, rawURL := range []string{
		"htps://collector.example.com/v1/traces",
		"https:///v1/traces",
	} {
		t.Run(rawURL, func(t *testing.T) {
			err := otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: rawURL})

			require.ErrorIs(t, err, otlp.ErrInvalidEndpoint)
		})
	}

	for _, address := range []string{
		"https://collector.example.com:4317",
		"collector.example.com",
		"collector.example.com:4317/v1/traces",
	} {
		t.Run(address, func(t *testing.T) {
			err := otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "grpc", Address: address})

			require.ErrorIs(t, err, otlp.ErrInvalidEndpoint)
		})
	}
}

func TestValidateEndpointRequiresEndpoint(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	require.ErrorIs(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Headers: headers}), otlp.ErrMissingEndpoint)
	require.NoError(t, otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: "https://collector.example.com/v1/traces", Headers: headers}))

	err := otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: "http://collector.example.com/v1/traces", Headers: headers})
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestValidateEndpointRejectsInvalidProtocol(t *testing.T) {
	err := otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "wrong", Address: "collector.example.com:4317"})

	require.ErrorIs(t, err, otlp.ErrInvalidProtocol)
}

func TestValidateEndpointIgnoresEnv(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "https://collector.example.com/v1/traces")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://collector.example.com")

	err := otlp.ValidateEndpoint(otlp.Endpoint{
		Protocol: "http",
		Headers:  map[string]string{"Authorization": "Bearer token"},
	})
	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}
