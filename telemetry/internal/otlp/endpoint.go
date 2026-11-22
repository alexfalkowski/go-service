package otlp

import (
	"github.com/alexfalkowski/go-service/v3/context"
	tls "github.com/alexfalkowski/go-service/v3/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/net"
	"github.com/alexfalkowski/go-service/v3/net/url"
	"github.com/alexfalkowski/go-service/v3/strings"
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

// ResolveEndpoint validates an explicitly configured OTLP endpoint and returns the address for the exporter.
// It returns the original address so the exporter parses the same input that was validated.
func ResolveEndpoint(endpoint Endpoint) (string, error) {
	if strings.IsEmpty(endpoint.Address) {
		return strings.Empty, ErrMissingEndpoint
	}

	switch endpoint.Protocol {
	case ProtocolHTTP:
		return resolveHTTP(endpoint)
	case ProtocolGRPC:
		return resolveGRPC(endpoint)
	default:
		return strings.Empty, ErrInvalidProtocol
	}
}

func resolveHTTP(endpoint Endpoint) (string, error) {
	if len(endpoint.Headers) == 0 {
		return resolveHeaderlessHTTP(endpoint.Address)
	}

	_, err := url.ParseSecureURL(endpoint.Address)
	if errors.Is(err, url.ErrInsecureURL) {
		return strings.Empty, ErrInsecureEndpoint
	}
	if err != nil {
		return strings.Empty, ErrInvalidEndpoint
	}

	return endpoint.Address, nil
}

func resolveHeaderlessHTTP(address string) (string, error) {
	u, err := url.Parse(address)
	if err != nil {
		return strings.Empty, err
	}

	httpScheme := u.Scheme == "http" || u.Scheme == "https"
	if !httpScheme || strings.IsEmpty(u.Hostname()) {
		return strings.Empty, ErrInvalidEndpoint
	}

	return address, nil
}

func resolveGRPC(endpoint Endpoint) (string, error) {
	host, port, err := net.SplitHostPort(endpoint.Address)
	if err != nil || strings.IsEmpty(host) || strings.IsEmpty(port) {
		return strings.Empty, ErrInvalidEndpoint
	}

	portNumber, err := net.LookupPort(context.Background(), "tcp", port)
	if err != nil || portNumber == 0 {
		return strings.Empty, ErrInvalidEndpoint
	}

	if len(endpoint.Headers) == 0 || endpoint.TLS != nil || net.IsLoopbackHost(host) {
		return endpoint.Address, nil
	}

	return strings.Empty, ErrInsecureEndpoint
}
