package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestMaxBytesHandler(t *testing.T) {
	t.Parallel()

	handler := http.MaxBytesHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _, err := io.ReadAll(req.Body)
		var maxBytesError *http.MaxBytesError
		require.ErrorAs(t, err, &maxBytesError)

		_, _ = res.Write([]byte("ok"))
	}), 1)

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewBufferString("too large"))

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestProtocolsReturnsSupportedHTTPVersions(t *testing.T) {
	t.Parallel()

	protocols := http.Protocols()

	require.True(t, protocols.HTTP1())
	require.True(t, protocols.HTTP2())
	require.True(t, protocols.UnencryptedHTTP2())
}

func TestParseTimeParsesHTTPHeaderTime(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	value := now.Format(http.TimeFormat)

	parsed, err := http.ParseTime(value)

	require.NoError(t, err)
	require.Equal(t, now.Truncate(time.Second.Duration()), parsed)
}

func TestTransportReturnsConfiguredRoundTripper(t *testing.T) {
	t.Parallel()

	cfg := &tls.Config{}

	transport := http.Transport(cfg)

	require.NotNil(t, transport.Proxy)
	require.NotNil(t, transport.DialContext)
	require.True(t, transport.ForceAttemptHTTP2)
	require.Equal(t, 100, transport.MaxIdleConns)
	require.Equal(t, 100, transport.MaxIdleConnsPerHost)
	require.Equal(t, 100, transport.MaxConnsPerHost)
	require.Equal(t, (90 * time.Second).Duration(), transport.IdleConnTimeout)
	require.Equal(t, (10 * time.Second).Duration(), transport.TLSHandshakeTimeout)
	require.Equal(t, time.Second.Duration(), transport.ExpectContinueTimeout)
	require.Same(t, cfg, transport.TLSClientConfig)
	require.NotNil(t, transport.Protocols)
	require.True(t, transport.Protocols.HTTP1())
	require.True(t, transport.Protocols.HTTP2())
	require.True(t, transport.Protocols.UnencryptedHTTP2())
}
