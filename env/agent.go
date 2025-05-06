package env

import "path/filepath"

// NewUserAgent for this service.
func NewUserAgent(name Name, ver Version) UserAgent {
	return UserAgent(filepath.Join(name.String(), ver.String()))
}

// UserAgent for this service.
type UserAgent string

// String representation of the user agent.
func (ua UserAgent) String() string {
	return string(ua)
}
