package otlp

import (
	"net/netip"

	"github.com/alexfalkowski/go-service/v2/context"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
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

// Endpoint describes an explicitly configured OTLP endpoint.
type Endpoint struct {
	// TLS configures TLS for gRPC exporters.
	TLS *tls.Config

	// Headers contains metadata sent to the OTLP collector.
	Headers map[string]string

	// Protocol selects the OTLP exporter protocol.
	Protocol string

	// Address is the configured OTLP destination.
	//
	// HTTP exporters expect a URL. gRPC exporters expect host:port.
	Address string
}

// ValidateEndpoint validates an explicitly configured OTLP endpoint.
func ValidateEndpoint(endpoint Endpoint) error {
	if strings.IsEmpty(endpoint.Address) {
		return ErrMissingEndpoint
	}

	switch endpoint.Protocol {
	case ProtocolHTTP:
		return validateHTTP(endpoint)
	case ProtocolGRPC:
		return validateGRPC(endpoint)
	default:
		return ErrInvalidProtocol
	}
}

func validateHTTP(endpoint Endpoint) error {
	u, err := url.Parse(endpoint.Address)
	if err != nil {
		return err
	}

	if (u.Scheme != "http" && u.Scheme != "https") || strings.IsEmpty(u.Host) {
		return ErrInvalidEndpoint
	}

	if u.Scheme != "http" || len(endpoint.Headers) == 0 || isLoopback(u.Hostname()) {
		return nil
	}

	return ErrInsecureEndpoint
}

func validateGRPC(endpoint Endpoint) error {
	host, port, err := net.SplitHostPort(endpoint.Address)
	if err != nil || strings.IsEmpty(host) || strings.IsEmpty(port) {
		return ErrInvalidEndpoint
	}
	portNumber, err := net.LookupPort(context.Background(), "tcp", port)
	if err != nil || portNumber == 0 {
		return ErrInvalidEndpoint
	}

	if len(endpoint.Headers) == 0 || endpoint.TLS != nil || isLoopback(host) {
		return nil
	}

	return ErrInsecureEndpoint
}

func isLoopback(host string) bool {
	if strings.ToLower(host) == "localhost" {
		return true
	}

	addr, err := netip.ParseAddr(host)
	return err == nil && addr.IsLoopback()
}
