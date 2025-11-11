package id

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id/ksuid"
	"github.com/alexfalkowski/go-service/v2/id/nanoid"
	"github.com/alexfalkowski/go-service/v2/id/ulid"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/id/xid"
)

// Module for fx.
var Module = di.Module(
	di.Constructor(ksuid.NewGenerator),
	di.Constructor(nanoid.NewGenerator),
	di.Constructor(ulid.NewGenerator),
	di.Constructor(uuid.NewGenerator),
	di.Constructor(xid.NewGenerator),
	di.Constructor(NewMap),
	di.Constructor(NewGenerator),
)
