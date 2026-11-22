package otlp_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v3/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/telemetry/internal/otlp"
	"github.com/stretchr/testify/require"
)

func TestResolveEndpoint(t *testing.T) {
	const zonedAddress = "https://[fe80::1%25é]:4318/v1/traces"
	headers := map[string]string{http.AuthorizationKey: "Bearer token"}

	tests := []struct {
		config  otlp.Config
		wantErr error
		name    string
	}{
		{name: "secure HTTP", config: otlp.Config{Protocol: "http", Endpoint: "HTTPS://collector.example.com/v1/%74races", Headers: headers}},
		{name: "HTTPS with IPv6 zone", config: otlp.Config{Protocol: "http", Endpoint: "https://[fe80::1%25en0]:4318/v1/traces", Headers: headers}},
		{name: "non-ASCII IPv6 zone with headers", config: otlp.Config{Protocol: "http", Endpoint: zonedAddress, Headers: headers}},
		{name: "non-ASCII IPv6 zone without headers", config: otlp.Config{Protocol: "http", Endpoint: zonedAddress}},
		{name: "lowercase localhost HTTP", config: otlp.Config{Protocol: "http", Endpoint: "http://localhost:4318/v1/traces", Headers: headers}},
		{name: "mixed case localhost HTTP", config: otlp.Config{Protocol: "http", Endpoint: "http://Localhost:4318/v1/traces", Headers: headers}},
		{name: "uppercase localhost HTTP", config: otlp.Config{Protocol: "http", Endpoint: "http://LOCALHOST:4318/v1/traces", Headers: headers}},
		{name: "loopback IP HTTP", config: otlp.Config{Protocol: "http", Endpoint: "http://127.0.0.1:4318/v1/traces", Headers: headers}},
		{name: "headerless HTTP", config: otlp.Config{Protocol: "http", Endpoint: "http://collector.example.com/v1/traces"}},
		{name: "headerless HTTPS", config: otlp.Config{Protocol: "http", Endpoint: "https://collector.example.com/v1/traces"}},
		{name: "loopback grpc", config: otlp.Config{Protocol: "grpc", Endpoint: "localhost:4317", Headers: headers}},
		{name: "mixed case localhost grpc", config: otlp.Config{Protocol: "grpc", Endpoint: "Localhost:4317", Headers: headers}},
		{name: "headerless grpc", config: otlp.Config{Protocol: "grpc", Endpoint: "collector.example.com:4317"}},
		{name: "TLS grpc", config: otlp.Config{Protocol: "grpc", Endpoint: "collector.example.com:4317", Headers: headers, TLS: &tls.Config{}}},
		{name: "insecure HTTP", config: otlp.Config{Protocol: "http", Endpoint: "http://collector.example.com/v1/traces", Headers: headers}, wantErr: otlp.ErrInsecureEndpoint},
		{name: "insecure grpc", config: otlp.Config{Protocol: "grpc", Endpoint: "collector.example.com:4317", Headers: headers}, wantErr: otlp.ErrInsecureEndpoint},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := otlp.ResolveEndpoint(tt.config)
			if tt.wantErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.config.Endpoint, address)
				return
			}

			require.ErrorIs(t, err, tt.wantErr)
			require.Empty(t, address)
		})
	}
}

func TestResolveEndpointInvalidURL(t *testing.T) {
	_, err := otlp.ResolveEndpoint(otlp.Config{
		Protocol: "http",
		Endpoint: "http://%",
		Headers:  map[string]string{http.AuthorizationKey: "Bearer token"},
	})

	require.ErrorIs(t, err, otlp.ErrInvalidEndpoint)
}

func TestResolveEndpointRejectsInvalidEndpoint(t *testing.T) {
	httpTests := []struct {
		name    string
		address string
	}{
		{name: "rejects unsupported URL scheme", address: "htps://collector.example.com/v1/traces"},
		{name: "rejects URL without host", address: "https:///v1/traces"},
		{name: "rejects URL with port and no hostname", address: "https://:4318/v1/traces"},
	}

	for _, tt := range httpTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := otlp.ResolveEndpoint(otlp.Config{Protocol: "http", Endpoint: tt.address})

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
			_, err := otlp.ResolveEndpoint(otlp.Config{Protocol: "grpc", Endpoint: tt.address})

			require.ErrorIs(t, err, otlp.ErrInvalidEndpoint)
		})
	}
}

func TestResolveEndpointRequiresEndpoint(t *testing.T) {
	headers := map[string]string{http.AuthorizationKey: "Bearer token"}

	_, err := otlp.ResolveEndpoint(otlp.Config{Protocol: "http", Headers: headers})
	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}

func TestResolveEndpointRejectsInvalidProtocol(t *testing.T) {
	_, err := otlp.ResolveEndpoint(otlp.Config{Protocol: "wrong", Endpoint: "collector.example.com:4317"})

	require.ErrorIs(t, err, otlp.ErrInvalidProtocol)
}

func TestResolveEndpointIgnoresEnv(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "https://collector.example.com/v1/traces")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://collector.example.com")

	_, err := otlp.ResolveEndpoint(otlp.Config{
		Protocol: "http",
		Headers:  map[string]string{http.AuthorizationKey: "Bearer token"},
	})
	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}
