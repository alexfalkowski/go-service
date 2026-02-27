package id

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id/ksuid"
	"github.com/alexfalkowski/go-service/v2/id/nanoid"
	"github.com/alexfalkowski/go-service/v2/id/ulid"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/id/xid"
)

// Module wires ID generation into Fx/Dig.
//
// It provides constructors for the built-in generator implementations:
//   - *ksuid.Generator via ksuid.NewGenerator (kind "ksuid")
//   - *nanoid.Generator via nanoid.NewGenerator (kind "nanoid")
//   - *ulid.Generator via ulid.NewGenerator (kind "ulid")
//   - *uuid.Generator via uuid.NewGenerator (kind "uuid")
//   - *xid.Generator via xid.NewGenerator (kind "xid")
//
// It then constructs a generator registry (*id.Map) via NewMap and finally selects a concrete
// Generator via NewGenerator using Config.Kind.
//
// Disabled behavior: when *id.Config is nil/disabled, NewGenerator returns (nil, nil) so downstream
// consumers can treat ID generation as optional.
var Module = di.Module(
	di.Constructor(ksuid.NewGenerator),
	di.Constructor(nanoid.NewGenerator),
	di.Constructor(ulid.NewGenerator),
	di.Constructor(uuid.NewGenerator),
	di.Constructor(xid.NewGenerator),
	di.Constructor(NewMap),
	di.Constructor(NewGenerator),
)
