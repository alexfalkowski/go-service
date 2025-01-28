package header

import (
	"errors"
	"slices"
	"strings"
	"unique"
)

var (
	// BasicAuthorization scheme.
	BasicAuthorization = unique.Make("Basic")

	// BearerAuthorization scheme.
	BearerAuthorization = unique.Make("Bearer")
)

var (
	// AllAuthorizations supported by the header.
	AllAuthorizations = []string{
		BasicAuthorization.Value(),
		BearerAuthorization.Value(),
	}

	// ErrInvalidAuthorization header.
	ErrInvalidAuthorization = errors.New("header: authorization is invalid")

	// ErrNotSupportedScheme for http.
	ErrNotSupportedAuthorization = errors.New("header: authorization is not supported")
)

// ParseAuthorization header into type and credentials or error.
func ParseAuthorization(header string) (string, string, error) {
	key, value, ok := strings.Cut(header, " ")
	if !ok {
		return "", "", ErrInvalidAuthorization
	}

	if !containsAuthorization(key) {
		return "", "", ErrNotSupportedAuthorization
	}

	return key, value, nil
}

func containsAuthorization(scheme string) bool {
	return slices.ContainsFunc(AllAuthorizations, func(s string) bool { return s == scheme })
}
