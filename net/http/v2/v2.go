package v2

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// Transport for v1.
func Transport(cfg *tls.Config) http.RoundTripper {
	t := &http2.Transport{
		AllowHTTP:        true,
		IdleConnTimeout:  90 * time.Second,
		ReadIdleTimeout:  90 * time.Second,
		PingTimeout:      90 * time.Second,
		WriteByteTimeout: 90 * time.Second,
		TLSClientConfig:  cfg,
	}

	if cfg == nil {
		d := &net.Dialer{Timeout: time.Minute, KeepAlive: 30 * time.Second}

		t.AllowHTTP = true
		t.DialTLSContext = func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return d.DialContext(ctx, network, addr)
		}
	}

	return t
}
