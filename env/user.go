package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewUserID returns a service user identifier.
//
// It prefers the SERVICE_USER_ID environment variable when set (non-empty); otherwise it falls back
// to the service name.
//
// This value is commonly used when a stable "user" identity is required for integrations that need a
// username/user-id concept, but where the service itself is the actor.
func NewUserID(name Name) UserID {
	return UserID(cmp.Or(os.Getenv("SERVICE_USER_ID"), name.String()))
}

// UserID is the service user identifier.
type UserID string

// String returns the user id value as a string.
func (i UserID) String() string {
	return string(i)
}
