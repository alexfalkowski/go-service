package otlp

import (
	"net/netip"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ProtocolGRPC selects OTLP/gRPC exporters.
const ProtocolGRPC = "grpc"

// ProtocolHTTP selects OTLP/HTTP exporters.
const ProtocolHTTP = "http"

// ErrMissingEndpoint is returned when an OTLP exporter is enabled without an explicit endpoint.
var ErrMissingEndpoint = errors.New("otlp: missing endpoint")

// ErrInvalidEndpoint is returned when an OTLP exporter endpoint is not valid for
// the selected protocol.
var ErrInvalidEndpoint = errors.New("otlp: invalid endpoint")

// ErrInvalidProtocol is returned when an OTLP exporter protocol is unsupported.
var ErrInvalidProtocol = errors.New("otlp: invalid protocol")

// ErrInsecureEndpoint is returned when secret headers would be sent over a
// non-local cleartext endpoint.
var ErrInsecureEndpoint = errors.New("otlp: insecure endpoint")

// ValidateEndpoint validates an explicitly configured OTLP endpoint.
func ValidateEndpoint(protocol, endpoint string, headers map[string]string) error {
	if strings.IsEmpty(endpoint) {
		return ErrMissingEndpoint
	}

	switch protocol {
	case ProtocolHTTP:
		return validateHTTP(endpoint, headers)
	case ProtocolGRPC:
		return validateGRPC(endpoint, headers)
	default:
		return ErrInvalidProtocol
	}
}

func validateHTTP(rawURL string, headers map[string]string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	if (u.Scheme != "http" && u.Scheme != "https") || strings.IsEmpty(u.Host) {
		return ErrInvalidEndpoint
	}

	if u.Scheme != "http" || len(headers) == 0 || isLoopback(u.Hostname()) {
		return nil
	}

	return ErrInsecureEndpoint
}

func validateGRPC(endpoint string, headers map[string]string) error {
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil || strings.IsEmpty(host) {
		return ErrInvalidEndpoint
	}
	if _, err := net.LookupPort(context.Background(), "tcp", port); err != nil {
		return ErrInvalidEndpoint
	}

	if len(headers) == 0 || isLoopback(host) {
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
