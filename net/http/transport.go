package http

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/v2/time"
)

// Transport constructs a tuned *http.Transport with reasonable defaults and an optional TLS config.
//
// This helper is intended for services and clients that want consistent HTTP transport behavior without
// having to re-specify common timeouts and connection pool limits.
//
// Defaults applied:
//   - Proxy: http.ProxyFromEnvironment
//   - ForceAttemptHTTP2: true (enables HTTP/2 where supported by the server and TLS config)
//   - Dialer: 1m connect timeout, 30s TCP keepalive
//   - Connection pool limits: 100 max total idle, 100 max per host, 100 max conns per host
//   - Timeouts: 90s idle conn timeout, 10s TLS handshake timeout, 1s expect-continue timeout
//   - Protocols: set via Protocols() (go-service HTTP protocol configuration)
//
// TLS behavior:
//   - If cfg is non-nil it is assigned to Transport.TLSClientConfig.
//   - If cfg is nil, the standard library defaults apply.
func Transport(cfg *tls.Config) *http.Transport {
	dialer := &net.Dialer{Timeout: time.Minute, KeepAlive: 30 * time.Second}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		TLSClientConfig:       cfg,
		Protocols:             Protocols(),
	}
}
