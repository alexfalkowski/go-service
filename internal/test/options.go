package test

import (
	cache "github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

// WorldOption configures optional features on a World before it is created.
type WorldOption interface {
	apply(opts *options)
}

type options struct {
	verifier      token.Verifier
	access        access.Controller
	rt            http.RoundTripper
	generator     token.Generator
	logger        *logger.Logger
	breaker       *breaker.Config
	clientLimiter *limiter.Config
	serverLimiter *limiter.Config
	pg            *pg.Config
	cache         *worldCacheOptions
	httpHealth    *worldHealthOptions
	grpcHealth    *worldHealthOptions
	transport     *transport.Config
	telemetry     string
	loggerConfig  string
	secure        bool
	compression   bool
	http          bool
	grpc          bool
	debug         bool
	hello         bool
	rest          bool
	registerCache bool
}

type worldOptionFunc func(*options)

func (f worldOptionFunc) apply(o *options) {
	f(o)
}

type worldCacheOptions struct {
	config *cache.Config
	driver driver.Driver
}

func (o *options) cacheOptions() *worldCacheOptions {
	if o.cache == nil {
		o.cache = &worldCacheOptions{}
	}

	return o.cache
}

type worldHealthOptions struct {
	name         string
	url          string
	observations []HealthObservation
}

// HealthObservation describes a health observer kind and the probe names it should track.
type HealthObservation struct {
	Kind  string
	Names []string
}

// HealthObserve builds a HealthObservation for health-related world options.
func HealthObserve(kind string, names ...string) HealthObservation {
	return HealthObservation{
		Kind:  kind,
		Names: append([]string(nil), names...),
	}
}

// WithWorldSecure enables TLS for the transport and debug servers in the test world.
func WithWorldSecure() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.secure = true
	})
}

// WithWorldTelemetry selects the metrics exporter kind used by the world meter.
//
// The current helpers recognize "otlp" for OTLP metrics and fall back to the
// Prometheus test setup for any other value.
func WithWorldTelemetry(kind string) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.telemetry = kind
	})
}

// WithWorldClientLimiter installs the provided client-side rate limiter config.
func WithWorldClientLimiter(config *limiter.Config) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.clientLimiter = config
	})
}

// WithWorldBreaker installs the provided client-side circuit breaker config.
func WithWorldBreaker(config *breaker.Config) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.breaker = config
	})
}

// WithWorldServerLimiter installs the provided server-side rate limiter config.
func WithWorldServerLimiter(config *limiter.Config) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.serverLimiter = config
	})
}

// WithWorldCompression enables transport compression for clients created by the world.
func WithWorldCompression() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.compression = true
	})
}

// WithWorldTransportConfig overrides the world's transport config.
//
// Pair this with [WithWorldHTTP] or [WithWorldGRPC] when the test needs the corresponding transport
// server registered.
func WithWorldTransportConfig(config *transport.Config) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.transport = config
	})
}

// WithWorldRoundTripper overrides the HTTP round tripper used by world clients and event senders.
func WithWorldRoundTripper(rt http.RoundTripper) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.rt = rt
	})
}

// WithWorldHTTP enables registration of the HTTP transport server.
func WithWorldHTTP() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.http = true
	})
}

// WithWorldHTTPHealth registers the HTTP health routes on the world before it starts.
func WithWorldHTTPHealth(name, url string, observations ...HealthObservation) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.http = true
		o.httpHealth = &worldHealthOptions{
			name:         name,
			url:          url,
			observations: append([]HealthObservation(nil), observations...),
		}
	})
}

// WithWorldGRPC enables registration of the gRPC transport server.
func WithWorldGRPC() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.grpc = true
	})
}

// WithWorldGRPCHealth registers the gRPC health service on the world before it starts.
func WithWorldGRPCHealth(name, url string, observations ...HealthObservation) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.grpc = true
		o.grpcHealth = &worldHealthOptions{
			name:         name,
			url:          url,
			observations: append([]HealthObservation(nil), observations...),
		}
	})
}

// WithWorldDebug enables registration of the debug server.
func WithWorldDebug() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.debug = true
	})
}

// WithWorldHello registers the default GET /hello test handler on the world's mux.
func WithWorldHello() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.hello = true
	})
}

// WithWorldCacheConfig overrides the cache built for the world.
//
// When config is nil, the world cache is disabled instead of using the default
// Redis-backed test cache.
func WithWorldCacheConfig(config *cache.Config) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.cacheOptions().config = config
	})
}

// WithWorldCacheDriver overrides the cache driver used for the world's cache.
//
// If no custom cache driver is provided, the driver is built from the selected
// cache config.
func WithWorldCacheDriver(driver driver.Driver) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.cacheOptions().driver = driver
	})
}

// WithWorldRegisterCache registers the world's cache with the generic cache helpers.
//
// This mutates the package-level cache used by cache.Get and cache.Persist.
// NewWorld resets that package-level registration with tb.Cleanup, so this
// option is intended for tests that specifically exercise the generic cache
// helpers.
func WithWorldRegisterCache() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.registerCache = true
	})
}

// WithWorldPGConfig enables Postgres for the world using config or the default
// test config when config is nil.
func WithWorldPGConfig(config *pg.Config) WorldOption {
	return worldOptionFunc(func(o *options) {
		if config != nil {
			o.pg = config
		} else {
			o.pg = NewPGConfig()
		}
	})
}

// WithWorldLogger injects a prebuilt logger into the world instead of constructing one from config.
func WithWorldLogger(logger *logger.Logger) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.logger = logger
	})
}

// WithWorldLoggerConfig selects the named logger config variant used when the world builds its own logger.
//
// Recognized values are "json", "text", "tint", and "otlp". Empty or
// unrecognized values use the OTLP logger config.
func WithWorldLoggerConfig(config string) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.loggerConfig = config
	})
}

// WithWorldRest configures NewWorld to create a REST client backed by the world's HTTP transport.
func WithWorldRest() WorldOption {
	return worldOptionFunc(func(o *options) {
		o.rest = true
	})
}

// WithWorldToken overrides the token generator and verifier used by world clients and servers.
func WithWorldToken(generator token.Generator, verifier token.Verifier) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.generator = generator
		o.verifier = verifier
	})
}

// WithWorldAccessController overrides the access controller used by world servers.
func WithWorldAccessController(controller access.Controller) WorldOption {
	return worldOptionFunc(func(o *options) {
		o.access = controller
	})
}

func worldOptions(opts ...WorldOption) *options {
	os := &options{}
	for _, o := range opts {
		o.apply(os)
	}

	return os
}

func (w *World) registerOptions(opts *options) {
	if opts.hello {
		w.HandleHello()
	}
	if opts.httpHealth != nil {
		w.RegisterHTTPHealth(opts.httpHealth.name, opts.httpHealth.url, opts.httpHealth.observations...)
	}
	if opts.grpcHealth != nil {
		w.RegisterGRPCHealth(opts.grpcHealth.name, opts.grpcHealth.url, opts.grpcHealth.observations...)
	}
}

func transportConfig(opts *options) *transport.Config {
	if opts.transport != nil {
		return opts.transport
	}

	if opts.secure {
		return NewSecureTransportConfig()
	}

	return NewInsecureTransportConfig()
}

func debugConfig(opts *options) *debug.Config {
	if opts.secure {
		return NewSecureDebugConfig()
	}

	return NewInsecureDebugConfig()
}

func tlsConfig(opts *options) *tls.Config {
	if opts.secure {
		return NewTLSClientConfig()
	}

	return nil
}
