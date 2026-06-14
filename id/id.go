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
// It is returned by [NewGenerator] when [Config.Kind] does not match any generator registered in the [Map].
var ErrNotFound = errors.New("id: generator not found")

// MapParams defines dependencies used to construct a generator Map.
//
// It is intended for dependency injection ([go.uber.org/fx]/[go.uber.org/dig]). The default wiring is provided by [Module].
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

// NewMap constructs a Map from the supplied default generator dependencies.
//
// The returned registry includes these kinds:
//   - "uuid"
//   - "ksuid"
//   - "ulid"
//   - "nanoid"
//   - "xid"
//
// The map is fixed by standard wiring, which supplies each generator dependency before NewMap runs.
// Manual callers that pass nil generator fields register nil values for those kinds. Services that
// need custom generator kinds should provide custom DI wiring for [Map] or [Generator] selection.
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
func (m *Map) Get(kind string) Generator {
	return m.generators[kind]
}

// NewGenerator selects a [Generator] based on [Config.Kind] from the provided registry.
//
// Default behavior: if config is nil/disabled, [NewGenerator] returns the generator registered under
// kind "uuid".
//
// Enabled behavior: if config is enabled, [NewGenerator] looks up the generator for [Config.Kind] in m.
// If the kind is registered with a non-nil generator, it returns that generator. If the kind is not
// registered or is registered as nil, it returns [ErrNotFound].
func NewGenerator(config *Config, m *Map) (Generator, error) {
	if !config.IsEnabled() {
		return m.Get("uuid"), nil
	}

	g := m.Get(config.Kind)
	if g != nil {
		return g, nil
	}

	return nil, ErrNotFound
}
