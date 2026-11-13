package id

import (
	"github.com/alexfalkowski/go-service/v3/id/ksuid"
	"github.com/alexfalkowski/go-service/v3/id/nanoid"
	"github.com/alexfalkowski/go-service/v3/id/ulid"
	"github.com/alexfalkowski/go-service/v3/id/uuid"
	"github.com/alexfalkowski/go-service/v3/id/xid"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires ID generation into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for the built-in generator implementations:
//   - *[ksuid.Generator] via [ksuid.NewGenerator] (kind "ksuid")
//   - *[nanoid.Generator] via [nanoid.NewGenerator] (kind "nanoid")
//   - *[ulid.Generator] via [ulid.NewGenerator] (kind "ulid")
//   - *[uuid.Generator] via [uuid.NewGenerator] (kind "uuid")
//   - *[xid.Generator] via [xid.NewGenerator] (kind "xid")
//
// It then constructs a generator registry (*[Map]) via [NewMap] and finally selects a concrete
// [Generator] via [NewGenerator] using [Config.Kind].
//
// Default behavior: when *[Config] is nil/disabled, [NewGenerator] returns the "uuid" generator so
// downstream consumers always have a usable ID source.
//
// The ULID generator requires [github.com/alexfalkowski/go-service/v3/crypto/rand.Reader]. Standard
// top-level wiring supplies that dependency through [github.com/alexfalkowski/go-service/v3/module.Library]
// via [github.com/alexfalkowski/go-service/v3/crypto.Module]. Direct lower-level composition of
// Module must provide the same reader dependency.
var Module = signal.Module(
	signal.Constructor(ksuid.NewGenerator),
	signal.Constructor(nanoid.NewGenerator),
	signal.Constructor(ulid.NewGenerator),
	signal.Constructor(uuid.NewGenerator),
	signal.Constructor(xid.NewGenerator),
	signal.Constructor(NewMap),
	signal.Constructor(NewGenerator),
)
