package header

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
)

const (
	// BasicAuthorization is the HTTP Authorization scheme name for Basic authentication.
	//
	// When used in an Authorization header, it is typically formatted as:
	//
	//	Authorization: Basic <credentials>
	//
	// where <credentials> is usually a base64-encoded `username:password` value.
	BasicAuthorization = "Basic"

	// BearerAuthorization is the HTTP Authorization scheme name for Bearer token authentication.
	//
	// When used in an Authorization header, it is typically formatted as:
	//
	//	Authorization: Bearer <token>
	//
	// where <token> is an opaque access token (for example, a JWT).
	BearerAuthorization = "Bearer"
)

var (
	// AllAuthorizations lists the supported Authorization schemes for ParseAuthorization.
	//
	// Values are matched case-insensitively against the parsed scheme token.
	AllAuthorizations = []string{BasicAuthorization, BearerAuthorization}

	// ErrInvalidAuthorization is returned when an Authorization header cannot be parsed.
	//
	// This is returned when the header does not contain a scheme and value separated by a single ASCII space
	// (i.e. it cannot be split as "<scheme> <value>").
	ErrInvalidAuthorization = errors.New("header: authorization is invalid")

	// ErrNotSupportedAuthorization is returned when the Authorization scheme is not supported.
	//
	// This is returned when the parsed scheme is not present in AllAuthorizations.
	ErrNotSupportedAuthorization = errors.New("header: authorization is not supported")
)

// ParseAuthorization parses an HTTP Authorization header into scheme and value.
//
// The expected format is:
//
//	<scheme><space><value>
//
// Supported schemes are listed in AllAuthorizations (for example Basic and Bearer) and are
// matched case-insensitively.
//
// Error behavior:
//   - If the header cannot be split into two parts on the first ASCII space, it returns ErrInvalidAuthorization.
//   - If the parsed scheme is not supported, it returns ErrNotSupportedAuthorization.
//
// On error, the returned scheme and value are empty strings.
func ParseAuthorization(header string) (string, string, error) {
	key, value, ok := strings.Cut(header, strings.Space)
	if !ok {
		return strings.Empty, strings.Empty, ErrInvalidAuthorization
	}

	// Design note: this parser intentionally accepts all supported Authorization schemes.
	// Token middleware consumes only the credential value, so scheme-specific enforcement belongs
	// at the middleware/policy layer rather than in this shared parser.
	key, ok = canonicalAuthorization(key)
	if !ok {
		return strings.Empty, strings.Empty, ErrNotSupportedAuthorization
	}

	return key, value, nil
}

func canonicalAuthorization(scheme string) (string, bool) {
	switch strings.ToLower(scheme) {
	case strings.ToLower(BasicAuthorization):
		return BasicAuthorization, true
	case strings.ToLower(BearerAuthorization):
		return BearerAuthorization, true
	default:
		return strings.Empty, false
	}
}
