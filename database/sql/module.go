package sql

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module for fx.
var Module = di.Module(
	pg.Module,
)
