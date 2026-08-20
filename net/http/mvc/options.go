package mvc

import "github.com/alexfalkowski/go-service/v2/net/http"

// RouteOption configures MVC route registration.
type RouteOption interface {
	apply(opts *routeOptions)
}

type routeOptions struct {
	unauthenticated bool
}

// StaticOption configures static file response behavior.
type StaticOption interface {
	apply(opts *staticOptions)
}

type staticOptions struct {
	cacheControl string
	routeOptions
}

type routeOptionFunc func(*routeOptions)

func (f routeOptionFunc) apply(o *routeOptions) {
	f(o)
}

type staticOptionFunc func(*staticOptions)

func (f staticOptionFunc) apply(o *staticOptions) {
	f(o)
}

// WithRouteUnauthenticated marks an MVC route as not requiring transport token authentication.
func WithRouteUnauthenticated() RouteOption {
	return routeOptionFunc(func(options *routeOptions) {
		options.unauthenticated = true
	})
}

// WithCacheControl sets the Cache-Control response header for a static route.
func WithCacheControl(value string) StaticOption {
	return staticOptionFunc(func(options *staticOptions) {
		options.cacheControl = value
	})
}

// WithStaticUnauthenticated marks a static route as not requiring transport token authentication.
func WithStaticUnauthenticated() StaticOption {
	return staticOptionFunc(func(options *staticOptions) {
		options.unauthenticated = true
	})
}

func (options *routeOptions) httpOptions() []http.RouteOption {
	if options.unauthenticated {
		return []http.RouteOption{http.WithRouteUnauthenticated()}
	}

	return nil
}

func options(opts ...StaticOption) *staticOptions {
	options := &staticOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}

	return options
}
