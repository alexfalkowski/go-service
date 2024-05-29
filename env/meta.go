package env

import (
	"fmt"

	"github.com/alexfalkowski/go-service/os"
)

// Name of the service.
type Name string

// NewName for this service.
func NewName() Name {
	return Name(os.ExecutableName())
}

// UserAgent for this service.
type UserAgent string

// NewUserAgent for this service.
func NewUserAgent(name Name, ver Version) UserAgent {
	return UserAgent(fmt.Sprintf("%s/%s", name, ver))
}
