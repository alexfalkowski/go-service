package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewUserID returns the service user id.
//
// It prefers the SERVICE_USER_ID environment variable when set; otherwise it falls back to the service name.
func NewUserID(name Name) UserID {
	return UserID(cmp.Or(os.Getenv("SERVICE_USER_ID"), name.String()))
}

// UserID of the service.
type UserID string

// String returns the user id as a string.
func (i UserID) String() string {
	return string(i)
}
