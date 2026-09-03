package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestProtocolsReturnsSupportedHTTPVersions(t *testing.T) {
	t.Parallel()

	protocols := http.Protocols()

	require.True(t, protocols.HTTP1())
	require.True(t, protocols.HTTP2())
	require.True(t, protocols.UnencryptedHTTP2())
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
