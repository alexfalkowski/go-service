package otlp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPClientUsesConfiguredTimeoutAndTLS(t *testing.T) {
	client, err := otlp.NewHTTPClient(test.FS, &tls.Config{ServerName: "collector.example.com"}, time.Second)
	require.NoError(t, err)
	require.Equal(t, time.Second.Duration(), client.Timeout)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport.TLSClientConfig)
	require.Equal(t, "collector.example.com", transport.TLSClientConfig.ServerName)
}

func TestNewHTTPClientRejectsInvalidTLSConfig(t *testing.T) {
	client, err := otlp.NewHTTPClient(test.FS, &tls.Config{CA: "invalid ca"}, time.Second)
	require.ErrorIs(t, err, tls.ErrInvalidCA)
	require.Nil(t, client)
}

func TestNewHTTPClientReturnsRedirectWithoutFollowingIt(t *testing.T) {
	trustedHeaders := make(chan string, 1)
	attackerRequests := make(chan struct{}, 1)
	attacker := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		attackerRequests <- struct{}{}
	}))
	t.Cleanup(attacker.Close)

	trusted := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		trustedHeaders <- req.Header.Get("X-Api-Key")
		res.Header().Set("Location", attacker.URL+"/v1/traces")
		res.WriteHeader(http.StatusTemporaryRedirect)
	}))
	t.Cleanup(trusted.Close)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, trusted.URL+"/v1/traces", http.NoBody)
	require.NoError(t, err)
	req.Header.Set("X-Api-Key", "token")

	client, err := otlp.NewHTTPClient(test.FS, nil, 0)
	require.NoError(t, err)
	require.Equal(t, (10 * time.Second).Duration(), client.Timeout)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	require.NoError(t, res.Body.Close())

	select {
	case header := <-trustedHeaders:
		require.Equal(t, "token", header)
	default:
		require.Fail(t, "trusted endpoint did not receive a request")
	}

	select {
	case <-attackerRequests:
		require.Fail(t, "redirect target received a request")
	default:
	}
}
