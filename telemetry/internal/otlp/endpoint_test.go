package otlp_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/stretchr/testify/require"
)

func TestValidateEndpoint(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}

	require.NoError(t, otlp.ValidateEndpoint("", headers))
	require.NoError(t, otlp.ValidateEndpoint("https://collector.example.com/v1/traces", headers))
	require.NoError(t, otlp.ValidateEndpoint("http://localhost:4318/v1/traces", headers))
	require.NoError(t, otlp.ValidateEndpoint("http://127.0.0.1:4318/v1/traces", headers))
	require.NoError(t, otlp.ValidateEndpoint("http://collector.example.com/v1/traces", nil))

	err := otlp.ValidateEndpoint("http://collector.example.com/v1/traces", headers)
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestValidateEndpointInvalidURL(t *testing.T) {
	err := otlp.ValidateEndpoint("http://%", map[string]string{"Authorization": "Bearer token"})

	require.Error(t, err)
	require.NotErrorIs(t, err, otlp.ErrInsecureEndpoint)
}
