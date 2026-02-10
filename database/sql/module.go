package sql

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module provides the Fx module that wires SQL database support.
var Module = di.Module(
	pg.Module,
)
