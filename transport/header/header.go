package header

import (
	"errors"
	"slices"
	"strings"
)

const (
	// BasicAuthorization scheme.
	BasicAuthorization = "Basic"

	// BearerAuthorization scheme.
	BearerAuthorization = "Bearer"
)

var (
	// AllAuthorizations supported by the header.
	AllAuthorizations = []string{BasicAuthorization, BearerAuthorization}

	// ErrInvalidAuthorization header.
	ErrInvalidAuthorization = errors.New("header: authorization is invalid")

	// ErrNotSupportedScheme for http.
	ErrNotSupportedAuthorization = errors.New("header: authorization is not supported")
)

// ParseAuthorization header into type and credentials or error.
func ParseAuthorization(header string) (string, string, error) {
	s := strings.Split(header, " ")

	if len(s) != 2 {
		return "", "", ErrInvalidAuthorization
	}

	if !containsAuthorization(s[0]) {
		return "", "", ErrNotSupportedAuthorization
	}

	return s[0], s[1], nil
}

func containsAuthorization(scheme string) bool {
	return slices.ContainsFunc(AllAuthorizations, func(s string) bool { return s == scheme })
}
