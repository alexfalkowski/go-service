package header

import (
	"errors"
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
	ErrInvalidAuthorization = errors.New("authorization is invalid")

	// ErrNotSupportedScheme for http.
	ErrNotSupportedAuthorization = errors.New("authorization is not supported")
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
	for _, v := range AllAuthorizations {
		if v == scheme {
			return true
		}
	}

	return false
}
