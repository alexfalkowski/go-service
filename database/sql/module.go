package sql

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires SQL database support into Fx/Dig.
//
// This module is a thin composition layer that includes driver-specific SQL modules. At present it
// includes PostgreSQL support via `database/sql/pg.Module`.
//
// Consumers typically include this module alongside configuration wiring so that driver-specific Open
// functions can read DSNs/pool settings and return a master/slave connection pool.
var Module = di.Module(
	pg.Module,
)
