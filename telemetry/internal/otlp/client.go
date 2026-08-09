package otlp

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
)

// NewHTTPClient constructs an OTLP HTTP client that returns redirect responses
// without following them. A zero timeout uses the OpenTelemetry default of ten
// seconds.
func NewHTTPClient(fs *os.FS, cfg *config.Config, timeout time.Duration) (*http.Client, error) {
	dialer := &net.Dialer{
		Timeout:   (30 * time.Second).Duration(),
		KeepAlive: (30 * time.Second).Duration(),
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       (90 * time.Second).Duration(),
		TLSHandshakeTimeout:   (10 * time.Second).Duration(),
		ExpectContinueTimeout: time.Second.Duration(),
	}
	if cfg.IsEnabled() {
		tlsConfig, err := client.NewConfig(fs, cfg)
		if err != nil {
			return nil, err
		}
		transport.TLSClientConfig = tlsConfig
	}
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout.Duration(),
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return client, nil
}
