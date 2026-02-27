package env

// NewUserAgent returns the service User-Agent value.
//
// The returned value is formatted as:
//
//	"<name>/<version>"
//
// Where name and version are derived from the provided Name and Version values (including any
// normalization performed by Version.String).
//
// This value is commonly used for outbound HTTP clients so requests can be attributed to a specific
// service and version by upstreams.
func NewUserAgent(name Name, ver Version) UserAgent {
	return UserAgent(name.String() + "/" + ver.String())
}

// UserAgent is the HTTP User-Agent value for this service.
type UserAgent string

// String returns the User-Agent value as a string.
func (ua UserAgent) String() string {
	return string(ua)
}
