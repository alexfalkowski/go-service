package http

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// IgnoreRedirect returns redirect responses without following them.
func IgnoreRedirect(_ *Request, _ []*Request) error {
	return ErrUseLastResponse
}

// SameOrigin reports whether prev and next use the same URL origin.
func SameOrigin(prev, next *url.URL) bool {
	if prev == nil || next == nil {
		return false
	}

	return strings.ToLower(prev.Scheme) == strings.ToLower(next.Scheme) &&
		strings.ToLower(prev.Hostname()) == strings.ToLower(next.Hostname()) &&
		originPort(prev) == originPort(next)
}

func originPort(u *url.URL) string {
	if port := u.Port(); !strings.IsEmpty(port) {
		return port
	}

	switch strings.ToLower(u.Scheme) {
	case "http":
		return "80"
	case "https":
		return "443"
	default:
		return strings.Empty
	}
}

// IsCrossOriginRedirect reports whether req is a redirected request whose previous request used a different origin.
func IsCrossOriginRedirect(req *Request) bool {
	if req == nil || req.Response == nil || req.Response.Request == nil {
		return false
	}

	return !SameOrigin(req.Response.Request.URL, req.URL)
}

// SameOriginRedirect limits same-origin redirect chains to the standard Go limit of ten requests.
//
// It is intended for clients that add credentials or signatures in RoundTripper middleware, where Go's
// default cross-origin sensitive-header stripping is not enough because middleware may mint fresh credentials
// for each redirected request.
func SameOriginRedirect(req *Request, via []*Request) error {
	if len(via) == 0 {
		return nil
	}

	prev := via[len(via)-1].URL
	next := req.URL
	if !SameOrigin(prev, next) {
		return ErrUseLastResponse
	}

	if len(via) >= 10 {
		return errors.New("stopped after 10 redirects")
	}

	return nil
}
