package http

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/v2/time"
)

// Transport returns a tuned *http.Transport with reasonable defaults and optional TLS config.
//
// It enables HTTP/2 where possible and sets timeouts and connection limits suitable for services.
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
