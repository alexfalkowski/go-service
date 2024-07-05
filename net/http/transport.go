package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// Transport for HTTP.
func Transport(cfg *tls.Config) http.RoundTripper {
	d := &net.Dialer{Timeout: time.Minute, KeepAlive: 30 * time.Second}
	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           d.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		TLSClientConfig:       cfg,
	}

	return t
}
