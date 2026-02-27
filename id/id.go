package id

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id/ksuid"
	"github.com/alexfalkowski/go-service/v2/id/nanoid"
	"github.com/alexfalkowski/go-service/v2/id/ulid"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/id/xid"
)

// ErrNotFound is returned when generator selection fails because the configured kind is unknown.
//
// It is returned by NewGenerator when Config.Kind does not match any generator registered in the Map.
var ErrNotFound = errors.New("id: generator not found")

// MapParams defines dependencies used to construct a generator Map.
//
// It is intended for dependency injection (Fx/Dig). The default wiring is provided by id.Module.
type MapParams struct {
	di.In

	// KSUID is the KSUID generator registered under kind "ksuid".
	KSUID *ksuid.Generator

	// NanoID is the NanoID generator registered under kind "nanoid".
	NanoID *nanoid.Generator

	// ULID is the ULID generator registered under kind "ulid".
	ULID *ulid.Generator

	// UUID is the UUID generator registered under kind "uuid".
	UUID *uuid.Generator

	// XID is the XID generator registered under kind "xid".
	XID *xid.Generator
}

// NewMap constructs a Map pre-populated with default generators.
//
// The returned registry includes these kinds:
//   - "uuid"
//   - "ksuid"
//   - "ulid"
//   - "nanoid"
//   - "xid"
//
// Callers can add additional kinds or override existing kinds by mutating the returned Map
// (this package does not currently expose a Register method; the map is typically fixed by wiring).
func NewMap(params MapParams) *Map {
	return &Map{
		generators: map[string]Generator{
			"ksuid":  params.KSUID,
			"nanoid": params.NanoID,
			"ulid":   params.ULID,
			"uuid":   params.UUID,
			"xid":    params.XID,
		},
	}
}

// Map is a registry of ID generators keyed by kind (for example "uuid" or "ksuid").
//
// Map is not concurrency-safe. It is typically constructed during initialization and treated
// as read-only thereafter.
type Map struct {
	generators map[string]Generator
}

// Get returns the generator registered for kind.
//
// If no generator is registered for kind, Get returns nil.
func (f *Map) Get(kind string) Generator {
	return f.generators[kind]
}

// NewGenerator selects a Generator based on config.Kind from the provided registry.
//
// Disabled behavior: if config is nil/disabled, NewGenerator returns (nil, nil).
//
// Enabled behavior: if config is enabled, NewGenerator looks up the generator for config.Kind in m.
// If the kind is registered, it returns that generator. If the kind is not registered, it returns ErrNotFound.
func NewGenerator(config *Config, m *Map) (Generator, error) {
	if !config.IsEnabled() {
		return nil, nil
	}

	g := m.Get(config.Kind)
	if g != nil {
		return g, nil
	}

	return nil, ErrNotFound
}
