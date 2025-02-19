package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/os"
)

// NewName for this service.
func NewName() Name {
	return Name(cmp.Or(os.GetVariable("SERVICE_NAME"), os.ExecutableName()))
}

// Name of the service.
type Name string

// String representation of the name.
func (n Name) String() string {
	return string(n)
}

// NewUserAgent for this service.
func NewUserAgent(name Name, ver Version) UserAgent {
	return UserAgent(name.String() + "/" + ver.String())
}

// UserAgent for this service.
type UserAgent string

// String representation of the user agent.
func (ua UserAgent) String() string {
	return string(ua)
}
