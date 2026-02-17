package header

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
)

const (
	// BasicAuthorization is the HTTP Authorization scheme for Basic auth.
	BasicAuthorization = "Basic"

	// BearerAuthorization is the HTTP Authorization scheme for Bearer tokens.
	BearerAuthorization = "Bearer"
)

var (
	// AllAuthorizations lists supported Authorization schemes.
	AllAuthorizations = []string{BasicAuthorization, BearerAuthorization}

	// ErrInvalidAuthorization is returned when an Authorization header is malformed.
	ErrInvalidAuthorization = errors.New("header: authorization is invalid")

	// ErrNotSupportedAuthorization is returned when an Authorization scheme is unsupported.
	ErrNotSupportedAuthorization = errors.New("header: authorization is not supported")
)

// ParseAuthorization parses an Authorization header into scheme and credentials.
//
// Supported schemes are listed in AllAuthorizations.
// It returns ErrInvalidAuthorization when the header does not contain a scheme and value separated by a space.
// It returns ErrNotSupportedAuthorization when the scheme is not supported.
func ParseAuthorization(header string) (string, string, error) {
	key, value, ok := strings.Cut(header, strings.Space)
	if !ok {
		return strings.Empty, strings.Empty, ErrInvalidAuthorization
	}
	if !containsAuthorization(key) {
		return strings.Empty, strings.Empty, ErrNotSupportedAuthorization
	}

	return key, value, nil
}

func containsAuthorization(scheme string) bool {
	return slices.Contains(AllAuthorizations, scheme)
}
