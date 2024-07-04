package h2c

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// Transport for H2C.
func Transport() http.RoundTripper {
	d := &net.Dialer{Timeout: time.Minute, KeepAlive: 30 * time.Second}
	t := &http2.Transport{
		AllowHTTP:        true,
		IdleConnTimeout:  90 * time.Second,
		ReadIdleTimeout:  90 * time.Second,
		PingTimeout:      90 * time.Second,
		WriteByteTimeout: 90 * time.Second,
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return d.DialContext(ctx, network, addr)
		},
	}

	return t
}
