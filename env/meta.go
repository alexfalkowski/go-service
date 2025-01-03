package env

import (
	"github.com/alexfalkowski/go-service/os"
)

// NewName for this service.
func NewName() Name {
	return Name(os.ExecutableName())
}

// Name of the service.
type Name string

// NewUserAgent for this service.
func NewUserAgent(name Name, ver Version) UserAgent {
	return UserAgent(string(name) + "/" + ver.String())
}

// UserAgent for this service.
type UserAgent string
