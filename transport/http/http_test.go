package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestSecure(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldSecure(), test.WithWorldTelemetry("prometheus"), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	client := world.NewHTTP()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://github.com/alexfalkowski", http.NoBody)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	world.RequireStop()
}
