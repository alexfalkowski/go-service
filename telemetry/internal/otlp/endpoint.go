package otlp

import (
	"net/netip"
	"net/url"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ErrInsecureEndpoint is returned when secret headers would be sent over a non-local HTTP endpoint.
var ErrInsecureEndpoint = errors.New("otlp: insecure endpoint")

// ValidateEndpoint rejects accidental cleartext OTLP endpoints when exporter headers are configured.
func ValidateEndpoint(rawURL string, headers map[string]string) error {
	if strings.IsEmpty(rawURL) || len(headers) == 0 {
		return nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	if u.Scheme != "http" || isLoopback(u.Hostname()) {
		return nil
	}

	return ErrInsecureEndpoint
}

func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}

	addr, err := netip.ParseAddr(host)
	return err == nil && addr.IsLoopback()
}
