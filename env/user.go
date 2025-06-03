package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewUserID for this service.
func NewUserID(name Name) UserID {
	return UserID(cmp.Or(os.Getenv("SERVICE_USER_ID"), name.String()))
}

// UserID of the service.
type UserID string

// String representation of the user id.
func (i UserID) String() string {
	return string(i)
}
