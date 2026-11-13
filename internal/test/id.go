package test

import (
	"github.com/alexfalkowski/go-service/v3/crypto/rand"
	"github.com/alexfalkowski/go-service/v3/id"
	"github.com/alexfalkowski/go-service/v3/id/ksuid"
	"github.com/alexfalkowski/go-service/v3/id/nanoid"
	"github.com/alexfalkowski/go-service/v3/id/ulid"
	"github.com/alexfalkowski/go-service/v3/id/uuid"
	"github.com/alexfalkowski/go-service/v3/id/xid"
)

// Generators contains the standard ID generators exercised by tests.
var Generators = id.NewMap(id.MapParams{
	KSUID:  ksuid.NewGenerator(),
	NanoID: nanoid.NewGenerator(),
	ULID:   ulid.NewGenerator(rand.NewReader()),
	UUID:   uuid.NewGenerator(),
	XID:    xid.NewGenerator(),
})

// IDSequenceGenerator is an [id.Generator] test double that returns configured IDs in order.
type IDSequenceGenerator struct {
	IDs []string
}

// Generate returns the next configured ID.
func (g *IDSequenceGenerator) Generate() string {
	id := g.IDs[0]
	g.IDs = g.IDs[1:]

	return id
}

// StaticIDGenerator is an [id.Generator] test double that always returns the same ID.
type StaticIDGenerator string

// Generate returns the configured ID.
func (g StaticIDGenerator) Generate() string {
	return string(g)
}
