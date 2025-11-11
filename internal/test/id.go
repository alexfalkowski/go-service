package test

import (
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/id/ksuid"
	"github.com/alexfalkowski/go-service/v2/id/nanoid"
	"github.com/alexfalkowski/go-service/v2/id/ulid"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/id/xid"
)

// Generators for test.
var Generators = id.NewMap(id.MapParams{
	KSUID:  ksuid.NewGenerator(),
	NanoID: nanoid.NewGenerator(),
	ULID:   ulid.NewGenerator(rand.NewReader()),
	UUID:   uuid.NewGenerator(),
	XID:    xid.NewGenerator(),
})
