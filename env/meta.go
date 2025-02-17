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

// String representation of the name.
func (n Name) String() string {
	return string(n)
}

// NewUserAgent for this service.
func NewUserAgent(name Name, ver Version) UserAgent {
	return UserAgent(string(name) + "/" + ver.String())
}

// UserAgent for this service.
type UserAgent string

// String representation of the user agent.
func (ua UserAgent) String() string {
	return string(ua)
}
