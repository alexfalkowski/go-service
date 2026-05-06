package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/stretchr/testify/require"
)

func TestSecure(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client, err := http.NewClient(
		http.WithClientBreaker(
			breaker.WithSettings(breaker.Settings{}),
			breaker.WithFailureStatuses(http.StatusInternalServerError),
		),
		http.WithClientRoundTripper(server.Client().Transport),
	)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, http.NoBody)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
