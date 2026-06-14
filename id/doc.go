// Package id provides ID generation abstractions, registries, and wiring used by go-service.
//
// This package defines a small [Generator] interface and provides a registry ([Map]) of generators
// keyed by kind (e.g. "uuid", "ksuid", etc.). A concrete generator can be selected at runtime via
// [NewGenerator] using [Config.Kind].
//
// # Kinds and implementations
//
// Concrete generator implementations live in subpackages under `id/*` (for example
// [github.com/alexfalkowski/go-service/v2/id/uuid], [github.com/alexfalkowski/go-service/v2/id/ksuid],
// [github.com/alexfalkowski/go-service/v2/id/ulid], [github.com/alexfalkowski/go-service/v2/id/nanoid],
// and [github.com/alexfalkowski/go-service/v2/id/xid]). The [Module] wiring constructs these
// generators and registers them into a *[Map].
//
// Generator kinds are selected for different operational properties. The "uuid" generator is the
// default and is optimized for the request metadata hot path. The "xid" generator is intentionally
// available for compact, roughly sortable identifiers, but callers should not treat XIDs as opaque
// or unpredictable values.
//
// These generators are for operational identifiers such as request IDs, webhook IDs, and token jti
// values. They are not a secret-material API and should not be used as passwords, bearer tokens, or
// other credentials.
//
// # Configuration and enablement
//
// ID generation configuration is optional. A nil *[Config] selects the "uuid" kind from the generator
// registry. If configuration is present, [NewGenerator] uses [Config.Kind] exactly; an empty or unknown
// kind returns [ErrNotFound] unless a custom registry contains that key.
//
// Start with [Generator], [Config], [NewGenerator], [Map], and [Module].
package id
