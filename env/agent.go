package env

// NewUserAgent returns the service User-Agent value in the form "<name>/<version>".
func NewUserAgent(name Name, ver Version) UserAgent {
	return UserAgent(name.String() + "/" + ver.String())
}

// UserAgent for this service.
type UserAgent string

// String returns the user agent as a string.
func (ua UserAgent) String() string {
	return string(ua)
}
