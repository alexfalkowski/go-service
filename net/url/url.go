package url

import (
	"net/url"

	"github.com/alexfalkowski/go-service/v2/strings"
)

// URL is an alias for net/url.URL.
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
