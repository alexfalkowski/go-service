package header

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
	"golang.org/x/net/http/httpguts"
)

// BearerAuthorization is the HTTP Authorization scheme name for Bearer token authentication.
//
// When used in an Authorization header, it is typically formatted as:
//
//	Authorization: Bearer <token>
//
// where <token> is an opaque access token (for example, a JWT).
const BearerAuthorization = "Bearer"

// ForwardedIPs lists the forwarding headers used to derive a client IP address.
//
// Metadata extraction accepts these as trusted inputs, so callers should only rely on derived IPs behind
// trusted edge infrastructure that strips or overwrites client-supplied forwarding headers.
//
// The order reflects the preferred source when multiple headers are present.
var ForwardedIPs = [...]ForwardedIP{
	{HTTP: "X-Real-Ip", GRPC: "x-real-ip"},
	{HTTP: "CF-Connecting-Ip", GRPC: "cf-connecting-ip"},
	{HTTP: "True-Client-Ip", GRPC: "true-client-ip"},
	{HTTP: "X-Forwarded-For", GRPC: "x-forwarded-for"},
}

var (
	// ErrInvalidAuthorization is returned when an Authorization header cannot be parsed.
	//
	// This is returned when the header does not contain a scheme and value separated by a single ASCII space
	// (i.e. it cannot be split as "<scheme> <value>").
	ErrInvalidAuthorization = errors.New("header: authorization is invalid")

	// ErrNotSupportedAuthorization is returned when the Authorization scheme is not supported.
	//
	// This is returned when the parsed scheme is not Bearer.
	ErrNotSupportedAuthorization = errors.New("header: authorization is not supported")
)

// ForwardedIP describes one forwarding header used to derive a client IP address.
type ForwardedIP struct {
	// HTTP is the canonical HTTP header key.
	//
	// net/http.Header.Get avoids allocating when the lookup key is already canonical.
	HTTP string

	// GRPC is the lowercase gRPC metadata key and stored IP address kind.
	//
	// gRPC metadata is normalized to lowercase, so lowercase keys hit the direct
	// metadata lookup path instead of falling back to case-insensitive scanning.
	GRPC string
}

// ValidFieldName reports whether name is a valid HTTP header field name.
func ValidFieldName(name string) bool {
	return httpguts.ValidHeaderFieldName(name)
}

// ValidFieldValue reports whether value is a valid HTTP header field value.
func ValidFieldValue(value string) bool {
	return httpguts.ValidHeaderFieldValue(value)
}

// ParseBearer parses an HTTP Authorization header and returns its bearer token value.
//
// The expected format is:
//
//	Bearer <token>
//
// The Bearer scheme is matched case-insensitively.
//
// Error behavior:
//   - If the header cannot be split into two parts on the first ASCII space, it returns ErrInvalidAuthorization.
//   - If the parsed value is empty, it returns ErrInvalidAuthorization.
//   - If the parsed scheme is not Bearer, it returns ErrNotSupportedAuthorization.
//
// On error, the returned value is an empty string.
func ParseBearer(header string) (string, error) {
	key, value, ok := strings.Cut(header, strings.Space)
	if !ok {
		return strings.Empty, ErrInvalidAuthorization
	}

	if strings.ToLower(key) != strings.ToLower(BearerAuthorization) {
		return strings.Empty, ErrNotSupportedAuthorization
	}

	value = strings.TrimSpace(value)
	if strings.IsEmpty(value) {
		return strings.Empty, ErrInvalidAuthorization
	}

	return value, nil
}
