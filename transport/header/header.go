package header

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
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

	// ErrNotSupportedAuthorization for http.
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
	return slices.Contains(AllAuthorizations, scheme)
}
