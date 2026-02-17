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

// ErrNotFound is returned when the configured generator kind is unknown.
var ErrNotFound = errors.New("id: generator not found")

// MapParams defines dependencies used to construct a generator Map.
type MapParams struct {
	di.In
	KSUID  *ksuid.Generator
	NanoID *nanoid.Generator
	ULID   *ulid.Generator
	UUID   *uuid.Generator
	XID    *xid.Generator
}

// NewMap constructs a Map pre-populated with default generators.
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

// Map holds ID generators keyed by kind.
type Map struct {
	generators map[string]Generator
}

// Get returns the generator registered for kind.
func (f *Map) Get(kind string) Generator {
	return f.generators[kind]
}

// NewGenerator selects a generator based on config.Kind.
//
// It returns (nil, nil) when config is disabled.
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
