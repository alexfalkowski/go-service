package url

import (
	"fmt"
	"net/url"

	"github.com/alexfalkowski/go-service/v3/bytes"
	"github.com/alexfalkowski/go-service/v3/encoding/json"
	"github.com/alexfalkowski/go-service/v3/encoding/yaml"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/net"
	"github.com/alexfalkowski/go-service/v3/runtime"
	"github.com/alexfalkowski/go-service/v3/strings"
)

// URL is an alias for [net/url.URL].
type URL = url.URL

// Parse parses a raw URL into a URL structure.
func Parse(rawURL string) (*URL, error) {
	return url.Parse(rawURL)
}

// JoinPath returns a URL string with the provided path elements joined to base.
func JoinPath(base string, elem ...string) (string, error) {
	return url.JoinPath(base, elem...)
}

// SplitPath splits a slash-prefixed path into its first segment and the remaining path.
//
// It accepts paths shaped like "/first/rest" and returns ("first", "rest", true).
// If the path is not slash-prefixed, lacks a second segment, or has an empty
// first/rest segment, it returns ("", "", false).
func SplitPath(path string) (string, string, bool) {
	if !strings.HasPrefix(path, "/") {
		return strings.Empty, strings.Empty, false
	}

	first, rest, ok := strings.Cut(path[1:], "/")
	if !ok || strings.IsAnyEmpty(first, rest) {
		return strings.Empty, strings.Empty, false
	}

	return first, rest, true
}

// ErrInvalidSecureURL is returned when a value cannot be represented as a [SecureURL].
var ErrInvalidSecureURL = errors.New("url: expected an absolute HTTPS URL or a loopback HTTP URL")

// ErrInsecureURL identifies a cleartext HTTP URL with a non-loopback hostname.
// It wraps [ErrInvalidSecureURL].
var ErrInsecureURL = fmt.Errorf("%w: non-loopback HTTP", ErrInvalidSecureURL)

// SecureURL is a parsed, absolute HTTPS URL with a hostname. It also accepts HTTP
// URLs whose hostname is localhost or a loopback IP address, such as 127.0.0.1 or ::1.
// It is represented as a string in JSON, HJSON, TOML, and YAML configuration.
// Use a *SecureURL config field and add the `validate:"required"` tag when
// the URL is required. A nil pointer or a zero SecureURL represents an unset URL.
type SecureURL struct {
	url *URL
}

// ParseSecureURL parses raw as an absolute HTTPS URL or a loopback HTTP URL.
// Invalid input returns nil and an error matching [ErrInvalidSecureURL].
// Non-loopback HTTP URLs return [ErrInsecureURL]. Errors omit the raw input.
func ParseSecureURL(raw string) (*SecureURL, error) {
	url, err := Parse(raw)
	if err != nil {
		return nil, ErrInvalidSecureURL
	}

	host := url.Hostname()
	if strings.IsEmpty(host) {
		return nil, ErrInvalidSecureURL
	}

	https := url.Scheme == "https"
	http := url.Scheme == "http"
	loopbackHTTP := http && net.IsLoopbackHost(host)
	if https || loopbackHTTP {
		return &SecureURL{url: url}, nil
	}
	if http {
		return nil, ErrInsecureURL
	}

	return nil, ErrInvalidSecureURL
}

// MustParseSecureURL parses raw as a [SecureURL] and panics if parsing fails.
func MustParseSecureURL(raw string) *SecureURL {
	parsed, err := ParseSecureURL(raw)
	runtime.Must(err)
	return parsed
}

// IsZero reports whether u is unset.
func (u SecureURL) IsZero() bool {
	return u.url == nil
}

// String returns the serialized URL, or an empty string if u is unset.
func (u SecureURL) String() string {
	if u.IsZero() {
		return strings.Empty
	}

	return u.url.String()
}

// URL returns a deep copy of the parsed URL, or nil if u is unset.
func (u SecureURL) URL() *URL {
	return u.url.Clone()
}

// MarshalText encodes u as a URL string. An unset value returns [ErrInvalidSecureURL].
func (u SecureURL) MarshalText() ([]byte, error) {
	if u.IsZero() {
		return nil, ErrInvalidSecureURL
	}

	return bytes.Clone(strings.Bytes(u.String())), nil
}

// UnmarshalText parses text as a [SecureURL], updating u only on success.
func (u *SecureURL) UnmarshalText(text []byte) error {
	// The url URL can retain substrings of the input.
	url, err := ParseSecureURL(bytes.String(bytes.Clone(text)))
	if err != nil {
		return err
	}

	*u = *url
	return nil
}

// UnmarshalJSON decodes a JSON string as a [SecureURL].
func (u *SecureURL) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return ErrInvalidSecureURL
	}

	return u.UnmarshalText(strings.Bytes(text))
}

// UnmarshalYAML decodes a YAML string as a [SecureURL].
func (u *SecureURL) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode || node.ShortTag() != "!!str" {
		return ErrInvalidSecureURL
	}

	return u.UnmarshalText(strings.Bytes(node.Value))
}
