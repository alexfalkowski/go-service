package otlp_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/stretchr/testify/require"
)

func TestValidateEndpoint(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	tests := []struct {
		endpoint otlp.Endpoint
		wantErr  error
		name     string
	}{
		{name: "secure HTTP", endpoint: otlp.Endpoint{Protocol: "http", Address: "https://collector.example.com/v1/traces", Headers: headers}},
		{name: "lowercase localhost HTTP", endpoint: otlp.Endpoint{Protocol: "http", Address: "http://localhost:4318/v1/traces", Headers: headers}},
		{name: "mixed case localhost HTTP", endpoint: otlp.Endpoint{Protocol: "http", Address: "http://Localhost:4318/v1/traces", Headers: headers}},
		{name: "uppercase localhost HTTP", endpoint: otlp.Endpoint{Protocol: "http", Address: "http://LOCALHOST:4318/v1/traces", Headers: headers}},
		{name: "loopback IP HTTP", endpoint: otlp.Endpoint{Protocol: "http", Address: "http://127.0.0.1:4318/v1/traces", Headers: headers}},
		{name: "headerless HTTP", endpoint: otlp.Endpoint{Protocol: "http", Address: "http://collector.example.com/v1/traces"}},
		{name: "loopback gRPC", endpoint: otlp.Endpoint{Protocol: "grpc", Address: "localhost:4317", Headers: headers}},
		{name: "mixed case localhost gRPC", endpoint: otlp.Endpoint{Protocol: "grpc", Address: "Localhost:4317", Headers: headers}},
		{name: "headerless gRPC", endpoint: otlp.Endpoint{Protocol: "grpc", Address: "collector.example.com:4317"}},
		{name: "TLS gRPC", endpoint: otlp.Endpoint{Protocol: "grpc", Address: "collector.example.com:4317", Headers: headers, TLS: &tls.Config{}}},
		{name: "insecure HTTP", endpoint: otlp.Endpoint{Protocol: "http", Address: "http://collector.example.com/v1/traces", Headers: headers}, wantErr: otlp.ErrInsecureEndpoint},
		{name: "insecure gRPC", endpoint: otlp.Endpoint{Protocol: "grpc", Address: "collector.example.com:4317", Headers: headers}, wantErr: otlp.ErrInsecureEndpoint},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := otlp.ValidateEndpoint(tt.endpoint)
			if tt.wantErr == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
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
	httpTests := []struct {
		name    string
		address string
	}{
		{name: "rejects unsupported URL scheme", address: "htps://collector.example.com/v1/traces"},
		{name: "rejects URL without host", address: "https:///v1/traces"},
	}

	for _, tt := range httpTests {
		t.Run(tt.name, func(t *testing.T) {
			err := otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "http", Address: tt.address})

			require.ErrorIs(t, err, otlp.ErrInvalidEndpoint)
		})
	}

	grpcTests := []struct {
		name    string
		address string
	}{
		{name: "rejects URL address", address: "https://collector.example.com:4317"},
		{name: "rejects missing port", address: "collector.example.com"},
		{name: "rejects empty port", address: "collector.example.com:"},
		{name: "rejects zero port", address: "collector.example.com:0"},
		{name: "rejects path", address: "collector.example.com:4317/v1/traces"},
	}

	for _, tt := range grpcTests {
		t.Run(tt.name, func(t *testing.T) {
			err := otlp.ValidateEndpoint(otlp.Endpoint{Protocol: "grpc", Address: tt.address})

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
